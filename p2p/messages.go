package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/lovelycbm/bongcoin/blockchain"
	"github.com/lovelycbm/bongcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind MessageKind
	// byte로 하면 다양한 형태의 데이터를 가지고 올수 있음.
	Payload []byte
}


func makeMessage(kind MessageKind, payload interface{}) []byte {
	m:= Message{
		Kind: kind,
		Payload: utils.ToJSON(payload),
	}			
	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	b, err:= blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleError(err)
	m := makeMessage(MessageNewestBlock,b)
	p.inbox <- m
} 

func handleMsg(m *Message, p *peer){
	fmt.Printf("Peer : %s, Sent a message with kind of : %d\n",p.key,m.Kind)
	switch m.Kind {
		case MessageNewestBlock:
			var payload blockchain.Block			
			utils.HandleError(json.Unmarshal(m.Payload, &payload))
			fmt.Println(payload)
		
	}
}