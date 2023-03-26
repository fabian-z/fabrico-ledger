// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/SmartBFT-Go/consensus/v2/pkg/consensus"
	"github.com/SmartBFT-Go/consensus/v2/pkg/types"
	"github.com/SmartBFT-Go/consensus/v2/pkg/wal"
	"github.com/SmartBFT-Go/consensus/v2/smartbftprotos"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var fastConfig = types.Configuration{
	RequestBatchMaxCount:          10,
	RequestBatchMaxBytes:          10 * 1024 * 1024,
	RequestBatchMaxInterval:       10 * time.Millisecond,
	IncomingMessageBufferSize:     200,
	RequestPoolSize:               40,
	RequestForwardTimeout:         500 * time.Millisecond,
	RequestComplainTimeout:        2 * time.Second,
	RequestAutoRemoveTimeout:      3 * time.Minute,
	ViewChangeResendInterval:      5 * time.Second,
	ViewChangeTimeout:             1 * time.Minute,
	LeaderHeartbeatTimeout:        1 * time.Minute,
	LeaderHeartbeatCount:          10,
	NumOfTicksBehindBeforeSyncing: 10,
	CollectTimeout:                200 * time.Millisecond,
	LeaderRotation:                false,
	RequestMaxBytes:               10 * 1024,
	RequestPoolSubmitTimeout:      5 * time.Second,
}

// App implements all interfaces required by an application using this library
type App struct {
	ID              NodeID
	Delivered       chan *AppRecord
	Consensus       *consensus.Consensus
	Node            *Node
	logLevel        zap.AtomicLevel
	latestMD        *smartbftprotos.ViewMetadata
	lastDecision    *types.Decision
	clock           *time.Ticker
	heartbeatTime   chan time.Time
	viewChangeTime  chan time.Time
	secondClock     *time.Ticker
	logger          *zap.SugaredLogger
	lastRecord      lastRecord
	verificationSeq uint64

	// Signature Data
	nodeCert *x509.Certificate
	nodeKey  ed25519.PrivateKey
	caCert   *x509.CertPool

	// Fabrication Data Storage
	Store FabricationDataStore
}

type lastRecord struct {
	proposal   types.Proposal
	signatures []types.Signature
}

// Submit submits the client request
func (a *App) Submit(req Request) {
	a.Consensus.SubmitRequest(req.ToBytes())
}

// Sync synchronizes and returns the latest decision
func (a *App) Sync() types.SyncResponse {

	reconfigSync := types.ReconfigSync{InReplicatedDecisions: false}

	// TODO enhance experimental implementation, currently chooses first live node which is not itself
	// also has to implement:
	// further verifications (client signatures, etc)
	// load balancing
	// retries

	a.Node.Lock()
	var srcConn *grpc.ClientConn
	for id, nodeConn := range a.Node.nodeChannels {
		if id == a.ID {
			continue
		}
		srcConn = nodeConn
		a.logger.Debugf("Sync connection with node %v", id)
		break
	}
	a.Node.Unlock()

	if srcConn == nil {
		//Single node started
		return types.SyncResponse{Latest: *a.lastDecision, Reconfig: reconfigSync}
	}

	client := NewNodeExchangeClient(srcConn)

	fetchFrom := &BlockPosition{
		ViewId:         a.Consensus.Metadata.GetViewId(),
		LatestSequence: a.Consensus.Metadata.LatestSequence,
	}

	stream, err := client.FetchBlocks(context.TODO(), fetchFrom)
	if err != nil {
		// TODO implement retry logic instead of panicking
		a.logger.Panicln(err)
	}

	for {
		record, err := stream.Recv()
		if err == io.EOF {
			err = stream.CloseSend()
			if err != nil {
				a.logger.Panic(err)
			}
			break
		}

		if err != nil {
			a.logger.Panic(err)
		}

		batchPayload, err := asn1.Marshal(struct{ Requests [][]byte }{Requests: record.Batch})
		if err != nil {
			a.logger.Panic(err)
		}

		proposal := types.Proposal{
			Payload:  batchPayload,
			Metadata: record.Metadata,
		}
		a.Deliver(proposal, nil)
		for _, req := range record.Batch {
			request := requestFromBytes(req)
			if request.Reconfig.InLatestDecision {
				reconfig := request.Reconfig.recconfigToUint(a.ID)
				reconfigSync = types.ReconfigSync{
					InReplicatedDecisions: true,
					CurrentNodes:          reconfig.CurrentNodes,
					CurrentConfig:         reconfig.CurrentConfig,
				}
			}
		}

	}

	return types.SyncResponse{Latest: *a.lastDecision, Reconfig: reconfigSync}
}

// RequestID returns info about the given request
func (a *App) RequestID(req []byte) types.RequestInfo {
	txn := requestFromBytes(req)
	return types.RequestInfo{
		ClientID: txn.ClientID,
		ID:       txn.ID,
	}
}

// VerifyProposal verifies the given proposal and returns the included requests
func (a *App) VerifyProposal(proposal types.Proposal) ([]types.RequestInfo, error) {
	blockData := batchFromBytes(proposal.Payload)
	requests := make([]types.RequestInfo, 0)
	for _, t := range blockData.Requests {
		req := requestFromBytes(t)
		reqInfo := types.RequestInfo{ID: req.ID, ClientID: req.ClientID}
		requests = append(requests, reqInfo)
	}
	return requests, nil
}

// RequestsFromProposal returns from the given proposal the included requests' info
func (a *App) RequestsFromProposal(proposal types.Proposal) []types.RequestInfo {
	blockData := batchFromBytes(proposal.Payload)
	requests := make([]types.RequestInfo, 0)
	for _, t := range blockData.Requests {
		req := requestFromBytes(t)
		reqInfo := types.RequestInfo{ID: req.ID, ClientID: req.ClientID}
		requests = append(requests, reqInfo)
	}
	return requests
}

// VerifyRequest verifies the given request and returns its info
func (a *App) VerifyRequest(val []byte) (types.RequestInfo, error) {
	req := requestFromBytes(val)
	return types.RequestInfo{ID: req.ID, ClientID: req.ClientID}, nil
}

// VerifyConsenterSig verifies a nodes signature on the given proposal
// Returns auxiliary data and error
func (a *App) VerifyConsenterSig(signature types.Signature, proposal types.Proposal) ([]byte, error) {

	appSig := &AppSignature{}
	rest, err := asn1.Unmarshal(signature.Value, appSig)
	if len(rest) > 0 {
		return nil, errors.New("unexpected trailing data")
	}
	if err != nil {
		return nil, err
	}
	if appSig.Version != 1 {
		return nil, errors.New("unexpected signature version")
	}

	cert, err := x509.ParseCertificate(appSig.Certificate)
	if err != nil {
		return nil, err
	}

	opts := x509.VerifyOptions{
		Roots: a.caCert,
	}

	if _, err := cert.Verify(opts); err != nil {
		return nil, errors.New("failed to verify certificate: " + err.Error())
	}

	if cert.PublicKeyAlgorithm != x509.Ed25519 {
		return nil, errors.New("unexpected public key algorithm")
	}

	pub, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("unexpected public key type")
	}

	// TODO standardize CN formatting
	if cert.Subject.CommonName != fmt.Sprintf("node%v", signature.ID) {
		return nil, errors.New("unexpected signer common name")
	}

	if !ed25519.Verify(pub, signature.Msg, appSig.Signature) {
		return nil, errors.New("invalid message signature")
	}

	sigMsg := &AppSignedMessage{}
	rest, err = asn1.Unmarshal(signature.Msg, sigMsg)
	if len(rest) > 0 {
		return nil, errors.New("unexpected trailing data")
	}
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(sigMsg.Payload, []byte(proposal.Digest())) {
		return nil, errors.New("invalid proposal for signature")
	}

	return sigMsg.AuxiliaryData, nil
}

func (a *App) AuxiliaryData(msg []byte) []byte {
	sigMsg := &AppSignedMessage{}
	rest, err := asn1.Unmarshal(msg, sigMsg)
	if len(rest) > 0 {
		panic("Unexpected trailing data")
	}
	if err != nil {
		panic(err)
	}

	return sigMsg.AuxiliaryData
}

// VerifySignature verifies a signature
func (a *App) VerifySignature(sig types.Signature) error {

	sigMsg := &AppSignature{}
	rest, err := asn1.Unmarshal(sig.Value, sigMsg)
	if len(rest) > 0 {
		return errors.New("unexpected trailing data")
	}
	if err != nil {
		return err
	}
	if sigMsg.Version != 1 {
		return errors.New("unexpected signature version")
	}

	cert, err := x509.ParseCertificate(sigMsg.Certificate)
	if err != nil {
		return err
	}

	opts := x509.VerifyOptions{
		Roots: a.caCert,
	}

	if _, err := cert.Verify(opts); err != nil {
		return errors.New("failed to verify certificate: " + err.Error())
	}

	if cert.PublicKeyAlgorithm != x509.Ed25519 {
		return errors.New("unexpected public key algorithm")
	}

	pub, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return errors.New("unexpected public key type")
	}

	// TODO standardize CN formatting
	if cert.Subject.CommonName != fmt.Sprintf("node%v", sig.ID) {
		return errors.New("unexpected signer common name")
	}

	if ed25519.Verify(pub, sig.Msg, sigMsg.Signature) {
		return nil
	} else {
		return errors.New("invalid message signature")
	}

}

// VerificationSequence returns the current verification sequence
func (a *App) VerificationSequence() uint64 {
	return atomic.LoadUint64(&a.verificationSeq)
}

type AppSignature struct {
	Version     int
	Signature   []byte
	Certificate []byte
}
type AppSignedMessage struct {
	Payload       []byte
	AuxiliaryData []byte `asn1:"omitempty"`
}

// Sign signs on the given value
func (a *App) Sign(in []byte) []byte {
	msg := AppSignature{
		Version:     1,
		Signature:   ed25519.Sign(a.nodeKey, in),
		Certificate: a.nodeCert.Raw,
	}

	out, err := asn1.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return out
}

// SignProposal signs on the given proposal
func (a *App) SignProposal(proposal types.Proposal, aux []byte) *types.Signature {
	// Aux Contains the Proposals received
	if len(aux) == 0 && len(a.Node.peers) > 1 {
		panic(fmt.Sprintf("didn't receive prepares from anyone, n=%d", len(a.Node.peers)))
	}

	msg, err := asn1.Marshal(AppSignedMessage{
		Payload:       []byte(proposal.Digest()),
		AuxiliaryData: aux,
	})
	if err != nil {
		panic(err)
	}

	sig, err := asn1.Marshal(AppSignature{
		Version:     1,
		Signature:   ed25519.Sign(a.nodeKey, msg),
		Certificate: a.nodeCert.Raw,
	})
	if err != nil {
		panic(err)
	}

	return &types.Signature{ID: uint64(a.ID), Value: sig, Msg: msg}
}

// AssembleProposal assembles a new proposal from the given requests
func (a *App) AssembleProposal(metadata []byte, requests [][]byte) types.Proposal {
	return types.Proposal{
		VerificationSequence: int64(atomic.LoadUint64(&a.verificationSeq)),
		Payload:              batch{Requests: requests}.toBytes(),
		Metadata:             metadata,
	}
}

func (a *App) MembershipChange() bool {
	return false
}

// Deliver delivers the given proposal
func (a *App) Deliver(proposal types.Proposal, signatures []types.Signature) types.Reconfig {
	defer func() {
		a.lastRecord = lastRecord{
			proposal:   proposal,
			signatures: signatures,
		}
	}()
	record := &AppRecord{
		Metadata: proposal.Metadata,
		Batch:    batchFromBytes(proposal.Payload),
	}
	a.Node.cb.add(record)
	a.lastDecision = &types.Decision{
		Proposal:   proposal,
		Signatures: signatures,
	}

	prevSeq := a.latestMD.LatestSequence

	a.latestMD = &smartbftprotos.ViewMetadata{}
	if err := proto.Unmarshal(proposal.Metadata, a.latestMD); err != nil {
		a.logger.Panic(err)
	}

	if prevSeq == a.latestMD.LatestSequence {
		a.logger.Panicf("Committed sequence %d twice", prevSeq)
	}

	a.Delivered <- record

	for _, req := range record.Batch.Requests {
		request := requestFromBytes(req)
		if request.Reconfig.InLatestDecision {
			reconfig := request.Reconfig.recconfigToUint(a.ID)
			return types.Reconfig{InLatestDecision: true, CurrentNodes: reconfig.CurrentNodes, CurrentConfig: reconfig.CurrentConfig}
		}
	}

	return types.Reconfig{InLatestDecision: false}
}

func newNode(id NodeID, walDir string, tlsPaths TLSPaths, rotateLeader bool, decisionsPerLeader uint64) *App {
	logConfig := zap.NewDevelopmentConfig()
	//logConfig := zap.NewProductionConfig()
	logger, _ := logConfig.Build()
	logger = logger.With(zap.Int64("id", int64(id)))
	sugaredLogger := logger.Sugar()

	cert, caPool, key, err := loadCertificate(tlsPaths)
	if err != nil {
		sugaredLogger.Panicf("Failed to initialize WAL: %s", err)
	}

	app := &App{
		clock:        time.NewTicker(time.Second),
		secondClock:  time.NewTicker(time.Second),
		ID:           id,
		Delivered:    make(chan *AppRecord, 100),
		logLevel:     logConfig.Level,
		latestMD:     &smartbftprotos.ViewMetadata{},
		lastDecision: &types.Decision{},
		logger:       sugaredLogger,

		nodeCert: cert,
		caCert:   caPool,
		nodeKey:  key,

		Store: NewMemoryStore(),
	}

	config := fastConfig
	config.SelfID = uint64(id)
	config.SyncOnStart = true
	//config.LeaderRotation = rotateLeader
	//config.DecisionsPerLeader = decisionsPerLeader

	writeAheadLog, walInitialEntries, err := wal.InitializeAndReadAll(app.logger, filepath.Join(walDir, nodeName), nil)
	if err != nil {
		sugaredLogger.Panicf("Failed to initialize WAL: %s", err)
	}

	if app.Consensus != nil && app.Consensus.Config.DecisionsPerLeader > 0 {
		config.DecisionsPerLeader = app.Consensus.Config.DecisionsPerLeader
	}
	if app.Consensus != nil && app.Consensus.Config.LeaderRotation {
		config.LeaderRotation = true
	}

	c := &consensus.Consensus{
		Config:             config,
		ViewChangerTicker:  app.secondClock.C,
		Scheduler:          app.clock.C,
		Logger:             app.logger,
		WAL:                writeAheadLog,
		Metadata:           *app.latestMD,
		Verifier:           app,
		Signer:             app,
		MembershipNotifier: app,
		RequestInspector:   app,
		Assembler:          app,
		Synchronizer:       app,
		Application:        app,
		WALInitialContent:  walInitialEntries,
		LastProposal:       app.lastRecord.proposal,
		LastSignatures:     app.lastRecord.signatures,
	}
	if app.heartbeatTime != nil {
		app.clock.Stop()
		c.Scheduler = app.heartbeatTime
	}
	if app.viewChangeTime != nil {
		app.secondClock.Stop()
		c.ViewChangerTicker = app.viewChangeTime
	}

	node, err := StartNode(id, c, app, tlsPaths)

	if err != nil {
		sugaredLogger.Panicf("Failed to start Node: %s", err)
	}

	c.Comm = node

	app.Consensus = c

	app.Node = node
	return app
}

// TODO check refactoring to unify network, ecies, and app signatures code paths?
func loadCertificate(paths TLSPaths) (*x509.Certificate, *x509.CertPool, ed25519.PrivateKey, error) {
	// Load server's certificate and private key

	certPEM, err := ioutil.ReadFile(paths.NodeCertificate)
	if err != nil {
		return nil, nil, nil, err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, nil, nil, errors.New("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	keyPEM, err := ioutil.ReadFile(paths.NodeKey)
	if err != nil {
		return nil, nil, nil, err
	}

	block, _ = pem.Decode(keyPEM)
	if block == nil {
		return nil, nil, nil, errors.New("failed to parse private key from PEM")
	}

	privParsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	priv, ok := privParsed.(ed25519.PrivateKey)
	if !ok {
		return nil, nil, nil, errors.New("invalid private key type")
	}

	caPEM, err := ioutil.ReadFile(paths.CaCertificate)
	if err != nil {
		return nil, nil, nil, err
	}

	pool := x509.NewCertPool()
	ok = pool.AppendCertsFromPEM(caPEM)
	if !ok {
		return nil, nil, nil, errors.New("no CA certificate found in path")
	}

	return cert, pool, priv, nil
}
