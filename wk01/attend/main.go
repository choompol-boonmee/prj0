package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"log"
	"bufio"
	"net/http"
	"context"
	"github.com/gorilla/mux"
	mrand "math/rand"
	"time"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	ma "github.com/multiformats/go-multiaddr"
//	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p-core/host"
	ps "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/crypto"

)

//var a1,a2,a3 config.Option;
var h2 host.Host;
var id2 string;
var sndID,rcvID,attID,rdfID string;

var path_getid string = "/attend/getid"
var path_home string = "/attend/{attID}/{rdfID}/{fname}/{lname}";
var path_home2 string = "/attend2/{attID}/{rdfID}/{fname}/{lname}";
var path_send string = "/attend/send";

var mode,lat,lng,pic,picuri string

func attendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	vars := mux.Vars(r)
	attID := vars["attID"];
	rdfID := vars["rdfID"];
	fname := vars["fname"];
	lname := vars["lname"];

	fmt.Printf("attID:%s\n", attID)
	fmt.Printf("rdfID:%s\n", rdfID)
	fmt.Printf("fname:%s\n", fname)
	fmt.Printf("lname:%s\n", lname)

	w.Header().Set("Content-Type", "text/html")
	body, _ := ioutil.ReadFile("attend.html")
	sbody := string(body)
	fmt.Fprintf(w, sbody, attID, rdfID, fname, lname)
}

func attend2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	vars := mux.Vars(r)
	attID := vars["attID"];
	rdfID := vars["rdfID"];
	fname := vars["fname"];
	lname := vars["lname"];

	fmt.Printf("attID:%s\n", attID)
	fmt.Printf("rdfID:%s\n", rdfID)
	fmt.Printf("fname:%s\n", fname)
	fmt.Printf("lname:%s\n", lname)

	w.Header().Set("Content-Type", "text/html")
	body, _ := ioutil.ReadFile("attend2.html")
	sbody := string(body)
	fmt.Fprintf(w, sbody, attID, rdfID, fname, lname)
}

func attendSend(w http.ResponseWriter, r *http.Request) {
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

	pic = strings.ReplaceAll(pic," ","+")
	picuri = "data:image/jpeg;base64,"+pic

	rcvID0 := attIdMap[attID]

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

			rs, _ := rw.ReadString('\n')

			fmt.Printf("RESULT=%s\n",rs);
			if(strings.HasPrefix(rs,"OK:")) {
				time := string([]rune(rs)[3:])
				fmt.Printf(">>>>>>>> SUCCESS %s\n", time);
//				string([]rune(rs)[3:len(att)-1])
				fmt.Fprintf(w, time)
//				fmt.Fprintf(w, "SUCCESS")
			} else {
				fmt.Printf(">>>>>>>> FAILED !!!\n");
				fmt.Fprintf(w, "FAILED")
			}

//			s.Read(make([]byte, 1))
		}
	}

}

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
	r.HandleFunc(path_home2, attend2Handler)
	r.HandleFunc(path_send, attendSend)
	http.Handle("/", r)

	fmt.Println("start server ");
	log.Fatal(http.ListenAndServe(":8080", nil))
}

