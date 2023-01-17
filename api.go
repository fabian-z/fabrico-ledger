package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-uuid"
)

type APIServer struct {
	Node                *Node
	FabricationEndpoint string
}

type NodeStatus struct {
	// TODO Extend this status message
	NodeID NodeID
}

type AvailableData struct {
	OriginatingID NodeID
	FileHash      string // hex string representing FabricationDataHash value
	Remaining     int
}

func (a *APIServer) ServeHTTP(endpoint string) error {
	http.HandleFunc("/api/status", a.NodeStatus)
	http.HandleFunc("/api/availabledata", a.AvailableData)

	http.HandleFunc("/api/addfile", a.AddFile)
	http.HandleFunc("/api/fabricate", a.Fabricate)

	http.Handle("/", http.FileServer(http.Dir("ui/dist/")))

	return http.ListenAndServe(endpoint, nil)
}

func (a *APIServer) NodeStatus(w http.ResponseWriter, _ *http.Request) {
	status := &NodeStatus{
		NodeID: a.Node.id,
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(status)

	if err != nil {
		a.Node.app.logger.Error(err)
	}
}

func (a *APIServer) AvailableData(w http.ResponseWriter, _ *http.Request) {
	var available []*AvailableData

	a.Node.cb.aggregationLock.Lock()
	defer a.Node.cb.aggregationLock.Unlock()

	for hash, allowed := range a.Node.cb.allowedNodes {
		for _, allow := range allowed {
			if a.Node.id == allow.NodeID && allow.RemainingCount > 0 {
				available = append(available, &AvailableData{
					OriginatingID: a.Node.cb.knownFiles[hash],
					FileHash:      fmt.Sprintf("%x", hash),
					Remaining:     allow.RemainingCount,
				})
			}
		}
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(available)

	if err != nil {
		a.Node.app.logger.Error(err)
	}
}

func (a *APIServer) AddFile(w http.ResponseWriter, req *http.Request) {

	err := req.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		a.Node.app.logger.Error(err)
		return
	}

	partCountList := req.MultipartForm.Value["partCount"]
	if len(partCountList) != 1 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	partCount, err := strconv.ParseUint(partCountList[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	nodeListForm := req.MultipartForm.Value["selectNode"]
	if len(nodeListForm) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var nodeList []NodeID
	for _, v := range nodeListForm {
		nodeID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		nodeList = append(nodeList, NodeID(nodeID))
	}

	// Valid request

	file := req.MultipartForm.File["uploadFile"]

	if len(file) != 1 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	uploadFile, err := file[0].Open()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	fileData, err := io.ReadAll(uploadFile)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// TODO factor out the following routine for clarity?

	hash, err := a.Node.app.Store.StoreData(fileData)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	binaryOriginID := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryOriginID, uint64(a.Node.id))

	var payload []byte
	payload = append(payload, hash[:]...)
	payload = append(payload, binaryOriginID...)

	reqID, err := uuid.GenerateUUID()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	a.Node.app.Submit(Request{
		ClientID: fmt.Sprintf("node-%v", a.Node.id),
		ID:       reqID,
		Type:     AddFile,
		Payload:  payload,
	})

	for _, node := range nodeList {
		binaryAllowedID := make([]byte, 8)
		binary.LittleEndian.PutUint64(binaryAllowedID, uint64(node))

		var payload []byte
		payload = append(payload, hash[:]...)
		payload = append(payload, binaryAllowedID...)

		reqID, err = uuid.GenerateUUID()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		a.Node.app.Submit(Request{
			ClientID: fmt.Sprintf("node-%v", a.Node.id),
			ID:       reqID,
			Type:     AllowFabrication,
			Payload:  payload,
			Count:    int(partCount),
		})
	}

	fmt.Fprintf(w, "Stored data with hash %x, payload %x for nodes %v", hash, payload, nodeList)

}

func (a *APIServer) Fabricate(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		a.Node.app.logger.Error(err)
		return
	}

	partSelectList := req.MultipartForm.Value["partSelect"]
	if len(partSelectList) != 1 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hash, err := hex.DecodeString(partSelectList[0])
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var address FabricationDataHash
	if len(hash) != len(address) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	copy(address[:], hash)

	fabricationData, err := a.Node.app.Store.GetData(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Data not found for hash %x", address), http.StatusNotFound)
		return
	}

	fabricationConn, err := net.Dial("tcp", a.FabricationEndpoint)
	if err != nil {
		http.Error(w, "Communication error", http.StatusInternalServerError)
		return
	}

	n, err := fabricationConn.Write(fabricationData)
	if n != len(fabricationData) || err != nil {
		http.Error(w, "Communication error", http.StatusInternalServerError)
		return
	}

}
