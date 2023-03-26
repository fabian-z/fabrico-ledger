package main

import (
	context "context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/libp2p/zeroconf/v2"
)

// Randomly generated UUID identifying this distributed ledger
const serviceUUID = "91a2d4ae-50d2-42d3-8744-4108cb80e134"
const dnsService = "_fabrico-ledger._tcp"
const dnsDomain = "local."

// TODO change to UUID as well?
type NodeID uint64
type FQDN string

type Peer struct {
	Self     bool
	PeerID   NodeID
	Hostname FQDN
	Port     uint16
	// TODO Certificate infos?
}

// Discoverer describes a stateful interface which uses LAN / WAN protocols to discover nodes
type Discoverer interface {
	Start(self Peer) error
	Stop() error
	AddPeer(node Peer) error // Used to add peers manually to discovery
	GetPeers() <-chan Peer
}

// List
type ListDiscoverer struct {
	sync.Mutex
	Peers     []Peer
	Listeners []chan Peer
}

func (d *ListDiscoverer) AddPeer(node Peer) error {
	d.Lock()
	defer d.Unlock()
	d.Peers = append(d.Peers, node)

	// inform listeners asynchronously
	go func() {
		for _, listener := range d.Listeners {
			listener <- node
		}
	}()

	return nil
}

func (d *ListDiscoverer) Start(self Peer) error {
	//no op
	return nil
}

func (d *ListDiscoverer) Stop() error {
	for _, listener := range d.Listeners {
		close(listener)
	}
	return nil
}

func (d *ListDiscoverer) GetPeers() <-chan Peer {
	d.Lock()
	var currentPeers []Peer
	copy(currentPeers, d.Peers)
	listener := make(chan Peer, 1)
	d.Listeners = append(d.Listeners, listener)
	d.Unlock()

	go func() {
		// inform of current peers asynchronously
		for _, peer := range currentPeers {
			if peer.Self {
				continue
			}
			listener <- peer
		}
	}()

	return listener
}

type MdnsDiscoverer struct {
	server         *zeroconf.Server
	cancelDiscover context.CancelFunc
	entries        chan *zeroconf.ServiceEntry
	selfInstance   string

	sync.Mutex
	Peers     []Peer
	Listeners []chan Peer
}

func (d *MdnsDiscoverer) AddPeer(node Peer) error {
	d.Lock()
	defer d.Unlock()
	d.Peers = append(d.Peers, node)

	// inform listeners asynchronously
	go func() {
		for _, listener := range d.Listeners {
			listener <- node
		}
	}()

	return nil
}

func (d *MdnsDiscoverer) Start(self Peer) error {

	// Start server
	// Note: Instance name must be unique to avoid conflicts

	d.selfInstance = fmt.Sprintf("node-%v", self.PeerID)

	var err error
	d.server, err = zeroconf.Register(d.selfInstance, dnsService, "local.", int(self.Port), nil, nil)
	if err != nil {
		return err
	}

	// Start discovery
	d.entries = make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == d.selfInstance {
				continue
			}
			instanceSplit := strings.Split(entry.Instance, "-")
			if len(instanceSplit) != 2 {
				continue
			}
			peerID, err := strconv.ParseUint(instanceSplit[1], 10, 64)
			if err != nil {
				continue
			}

			log.Println(entry)
			d.AddPeer(Peer{
				PeerID:   NodeID(peerID),
				Hostname: "localhost",
				Port:     uint16(entry.Port),
			})
		}
	}(d.entries)

	var ctx context.Context
	ctx, d.cancelDiscover = context.WithCancel(context.Background())
	// Discover all services on the network (e.g. _workstation._tcp)
	go zeroconf.Browse(ctx, dnsService, "local.", d.entries)

	return nil
}

func (d *MdnsDiscoverer) Stop() error {
	if d.server != nil {
		d.server.Shutdown()
	}
	if d.cancelDiscover != nil {
		d.cancelDiscover()
	}
	if d.entries != nil {
		close(d.entries)
	}
	return nil
}

func (d *MdnsDiscoverer) GetPeers() <-chan Peer {
	d.Lock()
	var currentPeers []Peer
	copy(currentPeers, d.Peers)
	listener := make(chan Peer, 1)
	d.Listeners = append(d.Listeners, listener)
	d.Unlock()

	go func() {
		// inform of current peers asynchronously
		for _, peer := range currentPeers {
			if peer.Self {
				continue
			}
			listener <- peer
		}
	}()

	return listener
}

// TODO Implement DHT for WAN
type DHTDiscovery struct {
}

func (d *DHTDiscovery) AddPeer(node Peer) error {
	return nil
}

func (d *DHTDiscovery) Start(self Peer) error {
	return nil
}

func (d *DHTDiscovery) Stop() error {
	return nil

}

func (d *DHTDiscovery) GetPeers() <-chan Peer {
	return nil
}
