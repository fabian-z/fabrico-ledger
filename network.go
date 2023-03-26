// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SmartBFT-Go/consensus/v2/pkg/types"
	"github.com/SmartBFT-Go/consensus/v2/smartbftprotos"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	incBuffSize = 1000
)

// Interface used to communicate network message to consensus backend
type handler interface {
	HandleMessage(sender uint64, m *smartbftprotos.Message)
	HandleRequest(sender uint64, req []byte)
	Stop()
}

type msgFrom struct {
	message proto.Message
	from    NodeID
}

// Node represents a node in a network
type Node struct {
	sync.RWMutex
	running      sync.WaitGroup
	id           NodeID
	shutdownChan chan struct{}

	transportCred credentials.TransportCredentials

	in chan Consensus

	peers map[NodeID]*Peer
	port  string

	listener net.Listener

	nodeChannels  map[NodeID]*grpc.ClientConn
	nodeExchanges map[NodeID]NodeExchangeClient

	h          handler
	cb         *committedBatches
	discoverer Discoverer // TODO implement multiple discovery methods
	app        *App
}

// AddOrUpdateNode adds or updates a node in the network
func StartNode(id NodeID, h handler, app *App, tlsPaths TLSPaths) (*Node, error) {
	var err error
	port := 3000 + int(id)

	node := &Node{
		in:           make(chan Consensus, incBuffSize),
		h:            h,
		shutdownChan: make(chan struct{}),
		peers:        make(map[NodeID]*Peer),

		nodeChannels:  make(map[NodeID]*grpc.ClientConn),
		nodeExchanges: make(map[NodeID]NodeExchangeClient),
		id:            id,
		app:           app,
	}

	node.transportCred, err = loadTLSCredentials(tlsPaths)
	if err != nil {
		return nil, err
	}

	selfPeer := &Peer{
		PeerID: id,
		Port:   uint16(port),
		Self:   true,
	}

	node.peers[id] = selfPeer
	node.cb = newCommittedBatches()
	/*
		node.discoverer = new(ListDiscoverer)
		node.discoverer.Start(*selfPeer)
		peerChan := node.discoverer.GetPeers()*/

	node.discoverer = new(MdnsDiscoverer)
	node.discoverer.Start(*selfPeer)
	peerChan := node.discoverer.GetPeers()

	go func() {
		for peer := range peerChan {
			if peer.PeerID == id {
				continue
			}
			node.Lock()
			node.peers[peer.PeerID] = &peer
			node.Unlock()

			err := node.Connect(peer.PeerID)
			if err != nil {
				node.app.logger.Error("Error connecting to node:", err)
				continue
			}

			// TODO improve reconfig logic!
			// Executing reconfig from all clients does not work well, limiting to one node imposes SPOF

			leaderID := node.app.Consensus.GetLeaderID()
			node.app.logger.Info("Leader during node add: ", leaderID)
			if id == NodeID(leaderID) {
				node.app.logger.Debug("Starting reconfig with Nodes: ", node.Nodes())

				node.app.Submit(Request{
					ClientID: "reconfig",
					ID:       fmt.Sprintf("add_node-%v_%v", peer.PeerID, time.Now().Unix()),
					Reconfig: Reconfig{
						InLatestDecision: true,
						CurrentNodes:     nodesToInt(node.Nodes()),
						CurrentConfig:    recconfigToInt(types.Reconfig{CurrentConfig: app.Consensus.Config}).CurrentConfig,
					},
				})
			}

		}
	}()

	// Useful for manually specifying peers on commandline
	for _, flagPeer := range flagPeers {
		splitPeer := strings.Split(flagPeer, ":")
		if len(splitPeer) != 3 {
			node.app.logger.Error("Error parsing peer id:host:port from command-line: ", flagPeer)
			continue
		}
		peerId, err := strconv.ParseUint(splitPeer[0], 10, 64)
		if err != nil {
			node.app.logger.Error("Error parsing peer id from command-line: ", err)
			continue
		}
		host := splitPeer[1]
		port, err := strconv.ParseUint(splitPeer[2], 10, 16)
		if err != nil {
			node.app.logger.Error("Error parsing peer port from command-line: ", err)
			continue
		}
		err = node.discoverer.AddPeer(Peer{
			Self:     false,
			PeerID:   NodeID(peerId),
			Hostname: FQDN(host),
			Port:     uint16(port),
		})
		if err != nil {
			node.app.logger.Error("Error adding peer from command-line: ", err)
			continue
		}
	}

	node.listener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(node.transportCred),
	)
	RegisterNodeExchangeServer(grpcServer, &nodeExchange{NodeId: uint64(node.id), messageChannel: node.in, committedBatches: node.cb, store: node.app.Store})
	go grpcServer.Serve(node.listener)

	go node.serve()
	return node, nil
}

func (n *Node) Connect(id NodeID) error {

	retry := func() {
		n.Lock()
		n.nodeChannels[id].Close()
		delete(n.nodeChannels, id)
		n.Unlock()
		go n.Connect(id)
	}

	n.Lock()
	peer := n.peers[id]
	n.Unlock()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", peer.Hostname, peer.Port), grpc.WithTimeout(5*time.Second), grpc.WithBlock(), grpc.WithTransportCredentials(n.transportCred), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                10 * time.Second,
		PermitWithoutStream: true,
	}))

	if err != nil {
		go retry()
		return err
	}

	n.Lock()
	n.nodeChannels[id] = conn
	n.Unlock()

	client := NewNodeExchangeClient(conn)

	n.Lock()
	n.nodeExchanges[id] = client
	n.Unlock()

	return nil
}

func (n *Node) send(source, target NodeID, msg proto.Message) {
	n.Lock()
	dstNode, found := n.nodeExchanges[target]
	n.Unlock()

	if !found {
		panic("node doesn't exist")
	}

	any, err := anypb.New(msg)
	if err != nil {
		panic(err)
	}
	//TODO proper error handling, shutdown
	_, err = dstNode.ConsensusMessage(context.TODO(), &Consensus{
		Node:    uint64(n.id),
		Message: any,
	})

	if err != nil {
		n.app.logger.Error("Dropped msg from", source, "to", target, "due to overflow")
	}
}

// SendConsensus sends a consensus related message to a target node
func (node *Node) SendConsensus(targetID uint64, m *smartbftprotos.Message) {
	node.send(node.id, NodeID(targetID), m)
}

// SendTransaction sends a client's request to a target node
func (node *Node) SendTransaction(targetID uint64, request []byte) {
	node.send(node.id, NodeID(targetID), &FwdMessage{Sender: uint64(node.id), Payload: request})
}

// Nodes returns the ids of all nodes in the network
func (node *Node) Nodes() []uint64 {
	//node.app.logger.Debug("Nodes called!")

	var res []uint64

	node.Lock()
	for k, _ := range node.peers {
		res = append(res, uint64(k))
	}
	node.Unlock()

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res
}

func (node *Node) serve() {
	for {
		node.app.logger.Debug("Trying to receive message")

		select {
		case <-node.shutdownChan:
			return
		case inMsg := <-node.in:
			node.RLock()
			handler := node.h
			node.RUnlock()

			id := inMsg.Node
			any := inMsg.Message

			node.app.logger.Debug("Received message from:", id)

			switch any.TypeUrl {
			case "type.googleapis.com/smartbftprotos.Message":
				msg := &smartbftprotos.Message{}
				err := any.UnmarshalTo(msg)
				if err != nil {
					node.app.logger.Panic(err)
				}

				//switch msg := msg.(type) {
				//case *smartbftprotos.Message:
				handler.HandleMessage(uint64(id), msg)

			case "type.googleapis.com/fabrico.FwdMessage":

				msg := &FwdMessage{}
				err := any.UnmarshalTo(msg)
				if err != nil {
					node.app.logger.Panic(err)
				}

				//switch msg := msg.(type) {
				//case *smartbftprotos.Message:
				handler.HandleRequest(uint64(id), msg.Payload)

			}

			//default:
			// TODO handle client request!
			//
			//}
		}
	}
}

type nodeExchange struct {
	NodeId           uint64
	messageChannel   chan<- Consensus
	committedBatches *committedBatches
	store            FabricationDataStore
	UnimplementedNodeExchangeServer
}

func (n *nodeExchange) ConsensusMessage(ctx context.Context, msg *Consensus) (*emptypb.Empty, error) {

	if msg.Node == n.NodeId {
		// Do not receive messages sent by ourselves
		return &emptypb.Empty{}, nil
	}

	n.messageChannel <- *msg
	return &emptypb.Empty{}, nil
}

func (n *nodeExchange) FetchBlocks(pos *BlockPosition, stream NodeExchange_FetchBlocksServer) error {

	records := n.committedBatches.readAll(smartbftprotos.ViewMetadata{
		ViewId:         pos.ViewId,
		LatestSequence: pos.LatestSequence,
	})

	for _, record := range records {
		err := stream.Send(&BlockRecord{
			Metadata: record.Metadata,
			Batch:    record.Batch.Requests,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *nodeExchange) DownloadContent(id *ContentID, stream NodeExchange_DownloadContentServer) error {

	var address FabricationDataHash

	if len(id.Id) != len(address) {
		return errors.New("invalid content address hash length")
	}

	copy(address[:], id.Id)

	content, err := n.store.GetData(address)
	if err != nil {
		return err
	}

	// Send file in 1MiB chunks, default GRPC limit is ~4M(i)B?
	for _, chunk := range splitByteSlice(content, 1024*1024) {
		err := stream.Send(&ContentChunk{
			Chunk: chunk,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Utility functions

func loadTLSCredentials(paths TLSPaths) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(paths.NodeCertificate, paths.NodeKey)
	if err != nil {
		return nil, err
	}

	caPEM, err := ioutil.ReadFile(paths.CaCertificate)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(caPEM)
	if !ok {
		return nil, errors.New("no CA certificate found in path")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:      pool,
		ClientCAs:    pool,
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // only verified P2P Connections
	}

	return credentials.NewTLS(config), nil
}

// https://gist.github.com/xlab/6e204ef96b4433a697b3
func splitByteSlice(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
