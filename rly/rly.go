package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	"os"
	"os/signal"
	"syscall"
	"github.com/libp2p/go-libp2p-core/network"
	"net"
	"bytes"
	"strings"
	"time"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var ctx context.Context
var node2 host.Host
var rlId string

var node1 host.Host

var ouID string

var out0 *peer.AddrInfo

var node2addr = "/ip4/0.0.0.0/tcp/10002"

var id2addr = make(map[string]*peer.AddrInfo)

func main() {

	ctx = context.Background()

	bPrv1, _ := ioutil.ReadFile("key1.prv")
	sPrv1 := string(bPrv1)
	bPrv1,_ = base64.StdEncoding.DecodeString(sPrv1)
	prv1, _ := crypto.UnmarshalPrivateKey(bPrv1);

	bPrv2, _ := ioutil.ReadFile("key2.prv")
	sPrv2 := string(bPrv2)
	bPrv2,_ = base64.StdEncoding.DecodeString(sPrv2)
	prv2, _ := crypto.UnmarshalPrivateKey(bPrv2);
	bId2, _ := ioutil.ReadFile("key2.pid")
	rlId = string(bId2)

	bId3, _ := ioutil.ReadFile("key3.pid")
	ouID = string(bId3)

	node1, _ = libp2p.New(ctx,
		libp2p.Identity(prv1),
		libp2p.EnableRelay(),
		libp2p.Ping(false),
	)

	node2, _ = libp2p.New(ctx,
		libp2p.Identity(prv2),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ListenAddrStrings(node2addr),
		libp2p.Ping(false),
	)
	if node2!=nil { fmt.Println("started") }

	rlya, _ := ma.NewMultiaddr(node2addr+"/ipfs/"+rlId)
	rlyb, _ := peer.AddrInfoFromP2pAddr(rlya)

	outa, _ := ma.NewMultiaddr("/p2p/" + rlId + "/p2p-circuit/ipfs/" + ouID)
	out0, _ = peer.AddrInfoFromP2pAddr(outa)
fmt.Println("ID3", ouID)
	if err := node1.Connect(ctx, *rlyb); err != nil { panic(err) }

	node1.SetStreamHandler("job0", func(s network.Stream) {
		buf := make([]byte,200)
		n,_ := s.Read(buf)
		if n>0 {
			id := string(buf[:n])
			go func() {
				b1, _ := ma.NewMultiaddr(id)
				b2, _ := peer.AddrInfoFromP2pAddr(b1)
fmt.Println("\nADDR:["+id+"]",b2)
				id2addr[id] = b2
			}()
		}
		s.Close()
	})

	l, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil { panic("err1") }
	defer l.Close()
fmt.Println("STARTED")
	for {
		in, err := l.Accept()
//fmt.Println("ACCEPT")
		if err != nil { panic("err3") }
		go income(in)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

func income(ou net.Conn) {

	var buf [2048]byte
	n, err := ou.Read(buf[0:])
	if err != nil { return }

	addr := out0
	data := buf[:n]
	id0 := ""
//fmt.Println("getin:", string(data))
	i1 := bytes.LastIndex(data, []byte("Cookie:"))
	if i1>=0 {
fmt.Println("cookie:")
		i2 :=bytes.Index(data[i1:], []byte("\r\n"))
		if i2>=0 {
fmt.Println("return:")
			ck := string(data[i1+8:i1+i2])
			cks := strings.Split(ck, ";")
fmt.Println("count", len(cks))
			for _, s := range cks {
				s = strings.TrimSpace(s)
fmt.Println("s", s)
				if strings.HasPrefix(s, "podi-id=") {
					rs := s[8:]
fmt.Println("rs", rs)
					a1 := id2addr[rs]
					if a1!= nil {
						id0 = rs
fmt.Println("GOTO", a1)
						addr = a1
					}
				}
			}
		}
	}

	rl, er1 := node1.NewStream(ctx, addr.ID, "job1")
	if er1 != nil {
		if len(id0)>0 {
			delete(id2addr, id0)
		}
		return
	}

	n, err = rl.Write(buf[0:n])
	if err != nil { return }

	ch := make(chan bool, 0)
	go monitor(rl, ou, ch)
	go st2co(rl, ou, ch)
	go co2st(rl, ou, ch)
}

func monitor(rl network.Stream, ou net.Conn, ch chan bool) {
    end := false
    for !end {
        select {
        case <-ch:
        case <-time.After(1 * time.Second):
            rl.Close()
            ou.Close()
            end = true
        }
    }
}

func st2co(rl network.Stream, ou net.Conn, ch chan bool) {
    var buf [256]byte
    for {
        n, err := rl.Read(buf[0:])
        if err != nil { break }
        n, err = ou.Write(buf[0:n])
        if err != nil { break }
        ch <- true
    }
}

func co2st(rl network.Stream, ou net.Conn, ch chan bool) {
    var buf [256]byte
    for {
        n, err := ou.Read(buf[0:])
        if err != nil { break }
        n, err = rl.Write(buf[0:n])
        if err != nil { break }
        ch <- true
    }
}


