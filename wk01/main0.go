package main

import (
	"fmt"
	"io/ioutil"
	"strings"
//	"io"
	"log"
	"bufio"
	"net/http"
	"context"
	"github.com/gorilla/mux"
//	"encoding/base64"
	mrand "math/rand"
//	"crypto/rand"
	"time"
//	"crypto/rsa"

	"github.com/libp2p/go-libp2p"
//	"github.com/libp2p/go-libp2p-core/network"
//	"github.com/libp2p/go-libp2p-core/peer"
	circuit "github.com/libp2p/go-libp2p-circuit"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p-core/host"
	ps "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/crypto"

)

var a1,a2,a3 config.Option;
var h1,h2,h3 host.Host;
var id1, id2, id3 string;
var sndID,rcvID,attID,rdfID string;

var path_getid string = "/attend/getid"
var path_home string = "/attend/{attID}/{rdfID}"
var path_send string = "/attend/send"

var mode,lat,lng,pic,picuri string

func attendHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	attID := vars["attID"];
	rdfID := vars["rdfID"];

	fmt.Printf("attID:%s\n", attID)
	fmt.Printf("rdfID:%s\n", rdfID)

	w.Header().Set("Content-Type", "text/html")
	body, _ := ioutil.ReadFile("attend.html")
	sbody := string(body)
	fmt.Fprintf(w, sbody, attID, rdfID)
}

func attendSend(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GOT DATA")
	r.ParseForm()

	attID = r.Form.Get("attID")
	rdfID = r.Form.Get("rdfID")
	mode = r.Form.Get("mode")
	lat = r.Form.Get("lat")
	lng = r.Form.Get("lng")
	pic = r.Form.Get("pic")

	fmt.Printf("ATTID: %s\n",attID)
	fmt.Printf("RDFID: %s\n",rdfID)
	fmt.Printf("MODE: %s\n",mode)
	fmt.Printf("LAT: %s\n",lat)
	fmt.Printf("LONG: %s\n",lng)
//	fmt.Printf("PIC: %s\n",pic)

	pic = strings.ReplaceAll(pic," ","+")
	picuri = "data:image/jpeg;base64,"+pic

	/*
	bin2, _ := base64.StdEncoding.DecodeString(pic)
	ioutil.WriteFile("webcam.jpg", bin2, 0644)
	*/

	/*
	fmt.Printf("LEN: %d\n",len(pic))
	pic = "<img src='"+pic+"'>"
	picuri := []byte(pic)
	ioutil.WriteFile("webcam.html", picuri, 0644)
	*/

	rcvID0 := attIdMap[attID]
//	fmt.Printf("RECHECK %s == %s\n", rcvID0, rcvID)

//	if len(rcvID)>0 {
	if len(rcvID0)>0 {
		rcvAddr, _ := ma.NewMultiaddr("/dns4/dip.popiang.com/tcp/9002/p2p/"+rcvID0)
		rcvPeer, _ := ps.AddrInfoFromP2pAddr(rcvAddr)
		s, er := h2.NewStream(context.Background(), rcvPeer.ID, "/cats")
		if er==nil {
			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			rw.WriteString(fmt.Sprintf("att=%s\n", attID))
			rw.WriteString(fmt.Sprintf("rdf=%s\n", rdfID))
			rw.WriteString(fmt.Sprintf("mod=%s\n", mode))
			rw.WriteString(fmt.Sprintf("lat=%s\n", lat))
			rw.WriteString(fmt.Sprintf("lng=%s\n", lng))
			rw.WriteString(fmt.Sprintf("pic=%s\n", picuri))
			rw.Flush()
			s.Read(make([]byte, 1))
		}
	}

}

/*
func enterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if len(rcvID)>0 {
		rcvAddr, _ := ma.NewMultiaddr("/dns4/dip.popiang.com/tcp/9002/p2p/"+rcvID)
		rcvPeer, _ := ps.AddrInfoFromP2pAddr(rcvAddr)
		s, er := h2.NewStream(context.Background(), rcvPeer.ID, "/cats")
		if er==nil {
			s.Read(make([]byte, 1))
		}
	}
	fmt.Fprintf(w, "ENTER " +"<a href='home'>HOME</a>");
}
*/

var attIdMap map[string]string;

func getidHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	rcvID2, _ := r.URL.Query()["rcvID"]
	fmt.Printf("recv id: %s\n", rcvID2)
	attID2, _ := r.URL.Query()["ATTENDID"]
	fmt.Printf("ATTENDID: %s\n", attID2)

	attID = attID2[0]
	rcvID = rcvID2[0]
	attIdMap[attID] = rcvID2[0]

	fmt.Fprintf(w, id2);
}

/*
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FAVICON REQ");
	w.Header().Set("Content-Type", "image/x-icon")
	body, _ := ioutil.ReadFile("favicon.ico")
	fmt.Fprintf(w, "%s", body)
}

func leaveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "LEAVE " +"<a href='home'>HOME</a>");
}
*/

func setRelay() {
	var prv1 crypto.PrivKey
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	prv1, _, _ = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	h2, _ = libp2p.New(context.Background(), libp2p.Identity(prv1),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9002"),
		libp2p.EnableRelay(circuit.OptHop))
	id2 = h2.ID().Pretty()

	fmt.Println()
	fmt.Printf("ID2 : %s\n",h2.ID());
}

func main() {

	attIdMap = make(map[string]string)

	setRelay();

	r := mux.NewRouter()
	r.HandleFunc(path_getid, getidHandler)
	r.HandleFunc(path_home, attendHandler)
	r.HandleFunc(path_send, attendSend)
	http.Handle("/", r)

	/*
	http.HandleFunc(path_getid, getidHandler)
	http.HandleFunc(path_home, attendHandler)
	http.HandleFunc(path_send, attendSend)
	*/
	/*
	http.HandleFunc(path_enter, enterHandler)
	http.HandleFunc(path_leave, leaveHandler)
	http.HandleFunc("/attend/favicon.ico", faviconHandler)
	*/

	fmt.Println("start server ");
	log.Fatal(http.ListenAndServe(":8080", nil))
}

