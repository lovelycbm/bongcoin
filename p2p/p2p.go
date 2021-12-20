package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lovelycbm/bongcoin/blockchain"
	"github.com/lovelycbm/bongcoin/utils"
)


var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port : 3000 will upgrade the request from :4000
	
	// 여기는 포트 3000 으로 시작한것으로 아래 conn은 4000임.

	ip:= utils.Splitter(r.RemoteAddr,":",0)
	openPort := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool { 
		return openPort != "" && ip!=""
	}
	fmt.Printf("\n%s wants to upgrade\n",openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)		
	utils.HandleError(err)
	
	initPeer(conn,ip,openPort)
	// time.Sleep(time.Second * 10)
	// peer.inbox<- []byte("hello form 3000!")
	
}

func AddPeer(address, port, openPort string,  broadcast bool) {
	// Port :4000 is requesting an upgrade from the port :3000
	// 여기는 포트 4000 으로 시작한것 아래 conn 은 3000임.
	fmt.Printf("\n%s wants to connect to port %s\n",openPort,port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleError(err)
	p := initPeer(conn, address, port)

	if broadcast{
		broadcastNewpeer(p)
		return 
	}	
	sendNewestBlock(p)
	// 그 무한루프하고 있는 inbox에다가 채널로 메세지를 보내면
	// 무한 루프중인곳에서 메세지를 받아서(write) 출력 (read)
	// peer.inbox<- []byte("hello form port 4000!")
	// time.Sleep(time.Second * 10)
	// conn.WriteMessage(websocket.TextMessage,[]byte("hello form port 4000!"))
}

func BroadcastNewBlock(b *blockchain.Block){
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		notifyNewBlock(b,p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _,p := range Peers.v {
		notifyNewTx(tx,p)
	}
}

func broadcastNewpeer(newPeer *peer) {
	for key, p := range Peers.v{
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload,p)
		}
	}
}