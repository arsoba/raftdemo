package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"raftdemo/service"
	"raftdemo/store"
)

var httpAddr string
var raftAddr string
var leaderAddr string
var nodeID string
var raftDir string

func init() {
	flag.StringVar(&httpAddr, "haddr", "", "Set the HTTP bind address")
	flag.StringVar(&raftAddr, "raddr", "", "Set Raft bind address")
	flag.StringVar(&leaderAddr, "join", "", "Set join address, if any")
	flag.StringVar(&nodeID, "id", "", "Node ID")
	flag.StringVar(&raftDir, "dir", "", "Raft Dir")
}

func main() {
	flag.Parse()

	os.MkdirAll(raftDir, 0700)

	s := store.New()
	s.RaftDir = raftDir
	s.RaftBind = raftAddr
	if err := s.Open(leaderAddr == "", nodeID); err != nil {
		log.Fatalf("Failed to open store: %s", err.Error())
	}

	h := service.New(httpAddr, s)
	if err := h.Start(); err != nil {
		log.Fatalf("Failed to start HTTP service: %s", err.Error())
	}

	if leaderAddr != "" {
		if err := joinToCluster(leaderAddr, raftAddr, nodeID); err != nil {
			log.Fatalf("Failed to join node at %s: %s", leaderAddr, err.Error())
		}
	}

	log.Println("Started successfully")

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("Exiting")
}

func joinToCluster(leaderAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", leaderAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
