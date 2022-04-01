package main

import (
	"fmt"
	"io"
	"net/http"
	"context"
	mrand "math/rand"
	"time"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-core/host"
	ps "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/crypto"

)

var path_getid string = "/attend/getid"
var url_pref string = "http://dip.popiang.com:8080"

var h3,h2 host.Host;
var id2 string;
var sndID,rcvID string;

func setRecv(c1 chan string) {

	fmt.Printf("Start Receiver...%s %s\n", url_pref, path_getid)

	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	prv1, _, _ := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	h3, _ = libp2p.New(context.Background(), libp2p.Identity(prv1),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9003"),
		libp2p.EnableRelay())

	rcvID0 := h3.ID().Pretty()
	fmt.Printf("rcvID[%s]\n", rcvID0)
	resp, _ := http.Get(url_pref+path_getid+"?rcvID="+rcvID0)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	sndID = string(body)

	addr, _ := ma.NewMultiaddr("/dns4/dip.popiang.com/tcp/9002/p2p/"+sndID)
	pe, _ := ps.AddrInfoFromP2pAddr(addr)
	if err := h3.Connect(context.Background(), *pe); err != nil { panic(err) }

	h3.SetStreamHandler("/cats", func(s network.Stream) {
		c1 <- "Data came"
//		fmt.Println("H3 stream started!")
		s.Close()
	})

}

func main() {
	c1 := make(chan string)
	setRecv(c1);
	for ;; {
		select {
		case dt := <-c1 :
			fmt.Printf("Got data: %s\n", dt)
		case <-time.After(10*time.Second) :
			fmt.Print(".")
		}

	}
	time.Sleep(100 * time.Second)
}

