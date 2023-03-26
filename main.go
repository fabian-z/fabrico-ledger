package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	selfID   *uint64
	nodeName string
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var flagPeers arrayFlags

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	selfID = flag.Uint64("id", 1, "id number")
	flag.Var(&flagPeers, "peers", "Set peers to add without discovery")
	flag.Parse()
	nodeName = "node" + strconv.FormatUint(*selfID, 10)
}

type TLSPaths struct {
	NodeCertificate string
	NodeKey         string
	CaCertificate   string
}

func main() {
	log.Println("Starting with ID", selfID)

	output, err := os.OpenFile("/tmp/"+strconv.FormatUint(*selfID, 10), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	tmpdir, _ := ioutil.TempDir("", "app-"+nodeName)

	tlsPaths := TLSPaths{
		NodeCertificate: path.Join("res", "ca", nodeName+".crt"),
		NodeKey:         path.Join("res", "ca", nodeName+".key"),
		CaCertificate:   path.Join("res", "ca", "ca.crt"),
	}

	node := newNode(NodeID(*selfID), tmpdir, tlsPaths, true, 10)

	// Allow for initial peer discovery..
	time.Sleep(10 * time.Second)

	err = node.Consensus.Start()
	if err != nil {
		panic(err)
	}

	apiSrv := APIServer{
		Node:                node.Node,
		FabricationEndpoint: "localhost:9001",
	}

	apiPort := 8000 + *selfID
	go apiSrv.ServeHTTP(":" + strconv.Itoa(int(apiPort)))

	go func() {
		time.Sleep(10 * time.Second)
		for i := 1; i < 3; i++ {
			reqID, err := uuid.NewRandom()
			if err != nil {
				log.Fatal(err)
			}
			node.Submit(Request{ID: reqID.String(), ClientID: fmt.Sprintf("node-%v", *selfID)})
			time.Sleep(1 * time.Second)
		}
	}()

	for delivery := range node.Node.app.Delivered {
		node.Node.app.logger.Debug("Delivered message: ", delivery)
		for _, v := range delivery.Batch.Requests {
			req := requestFromBytes(v)
			fmt.Fprintf(output, "%v - %v\n", req.ID, req.ClientID)
			log.Printf("Received delivered Request ID %v from ClientID %v\n", req.ID, req.ClientID)
		}
	}

}
