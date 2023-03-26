package main

import (
	"encoding/asn1"
	"encoding/binary"
	"sync"

	"github.com/SmartBFT-Go/consensus/v2/smartbftprotos"
	"google.golang.org/protobuf/proto"
)

type AllowCount struct {
	NodeID         NodeID
	RemainingCount int
}

func newCommittedBatches() *committedBatches {
	return &committedBatches{
		knownFiles:   make(map[FabricationDataHash]NodeID),
		allowedNodes: make(map[FabricationDataHash][]*AllowCount),
	}
}

type committedBatches struct {
	lock     sync.RWMutex
	latestMD smartbftprotos.ViewMetadata
	records  []*AppRecord

	// aggregated fields
	aggregationLock sync.Mutex
	knownFiles      map[FabricationDataHash]NodeID // mapping from file hash to originating node id
	allowedNodes    map[FabricationDataHash][]*AllowCount
}

func (cb *committedBatches) add(record *AppRecord) {
	cb.lock.Lock()

	md := &smartbftprotos.ViewMetadata{}
	if err := proto.Unmarshal(record.Metadata, md); err != nil {
		cb.lock.Unlock() // in case we recover later
		panic(err)
	}

	if cb.latestMD.ViewId > md.ViewId {
		cb.lock.Unlock()
		return
	}
	if cb.latestMD.LatestSequence >= md.LatestSequence {
		cb.lock.Unlock()
		return
	}
	cb.latestMD = *md
	cb.records = append(cb.records, record)

	cb.lock.Unlock()
	// Process aggregations
	// *Has to be done asynchronously to prevent consensus failures!*
	// TODO Client signatures (prevent spoofing)

	go func() {
		for _, reqBytes := range record.Batch.Requests {
			request := requestFromBytes(reqBytes)

			switch RequestType(request.Type) {
			case SystemReserved:
				// No aggregations from system messages
				continue

			case AddFile:
				//Payload is 64 bits hash, 64 bit uint originating node
				var address FabricationDataHash
				copy(address[:], request.Payload[:64])

				originatingNode := binary.LittleEndian.Uint64(request.Payload[64:])

				cb.aggregationLock.Lock()
				cb.knownFiles[address] = NodeID(originatingNode)
				cb.aggregationLock.Unlock()

				// TODO add more verifications for file ownership?
				// TODO difference between originating and distributing nodes

			case AllowFabrication:
				//64 bits Hash, 64 Bits Uint64 Allowed Node, Count Sets Maximum Parts

				var address FabricationDataHash
				copy(address[:], request.Payload[:64])

				allowedNode := binary.LittleEndian.Uint64(request.Payload[64:])

				allowCounts, ok := cb.allowedNodes[address]
				if !ok {
					cb.aggregationLock.Lock()
					cb.allowedNodes[address] = []*AllowCount{{
						NodeID:         NodeID(allowedNode),
						RemainingCount: request.Count,
					}}
					cb.aggregationLock.Unlock()
					continue
				}

				setCount := false
				// Check if node has allow count
				for _, count := range allowCounts {
					if count.NodeID == NodeID(allowedNode) {
						count.RemainingCount = count.RemainingCount + request.Count
						setCount = true
						break
					}
				}

				// Node needs new allow count
				if !setCount {
					cb.aggregationLock.Lock()
					cb.allowedNodes[address] = append(cb.allowedNodes[address], &AllowCount{
						NodeID:         NodeID(allowedNode),
						RemainingCount: request.Count,
					})
					cb.aggregationLock.Unlock()
				}

			case AnnounceFabrication:
				//64bits Hash, Count Number of Parts intended to produce
				panic("not yet implemented")

			case CancelFabrication:
				//64bits Hash, Count Number of Parts not produced
				panic("not yet implemented")

			}
		}
	}()

}

func (cb *committedBatches) readAll(from smartbftprotos.ViewMetadata) []*AppRecord {
	cb.lock.RLock()
	defer cb.lock.RUnlock()

	var res []*AppRecord

	for _, entry := range cb.records {
		md := &smartbftprotos.ViewMetadata{}
		if err := proto.Unmarshal(entry.Metadata, md); err != nil {
			panic(err)
		}
		if md.ViewId < from.ViewId || md.LatestSequence <= from.LatestSequence {
			continue
		}
		res = append(res, &AppRecord{
			Metadata: entry.Metadata,
			Batch:    entry.Batch,
		})
	}
	return res
}

// Request represents a client's request
type Request struct {
	ClientID string // Currently always originating NodeID
	ID       string // Request UUID
	Type     RequestType
	Payload  []byte // Message type specific payload data
	Count    int    // Message type specific count value
	Reconfig Reconfig
}

// AddFile Payload: 64 bits Hash, 64 Bits Uint64 Originating Node
// AllowFabrication Payload: 64 bits Hash, 64 Bits Uint64 Allowed Node, Count Sets Maximum Parts

// TODO implement Announcement / Cancellation!
// Announcing and Canceling protects against Network / Power Glitches to prevent production of excess parts
// AnnounceFabrication Payload: 64bits Hash, Count Number of Parts intended to produce
// CancelFabrication Payload: 64bits Hash, Count Number of Parts not produced - must correspond to earlier announce request

type RequestType int

const (
	SystemReserved RequestType = iota
	AddFile
	AllowFabrication
	AnnounceFabrication
	CancelFabrication
)

func (t RequestType) String() string {
	return [...]string{"SystemReserved", "AddFile", "AllowFabrication", "AnnounceFabrication", "CancelFabrication"}[t]
}

// ToBytes returns a byte array representation of the request
func (txn Request) ToBytes() []byte {
	rawTxn, err := asn1.Marshal(txn)
	if err != nil {
		panic(err)
	}
	return rawTxn
}

func requestFromBytes(req []byte) *Request {
	var r Request
	rest, err := asn1.Unmarshal(req, &r)
	if len(rest) > 0 {
		panic("unexpected trailing data")
	}
	if err != nil {
		panic(err)
	}
	return &r
}

type batch struct {
	Requests [][]byte
}

func (b batch) toBytes() []byte {
	rawBlock, err := asn1.Marshal(b)
	if err != nil {
		panic(err)
	}
	return rawBlock
}

func batchFromBytes(rawBlock []byte) *batch {
	var block batch
	asn1.Unmarshal(rawBlock, &block)
	return &block
}

// AppRecord represents a committed batch and metadata
type AppRecord struct {
	Batch    *batch
	Metadata []byte
}
