package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lovelycbm/bongcoin/blockchain"
	"github.com/lovelycbm/bongcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
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
	fmt.Printf("Sending newest block to %s\n",p.key)
	b, err:= blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleError(err)
	m := makeMessage(MessageNewestBlock,b)
	p.inbox <- m
} 

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest,nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer){
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.BlockChain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block ,p *peer){
	m := makeMessage(MessageNewBlockNotify,b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx ,p *peer){
	m := makeMessage(MessageNewTxNotify,tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}


func handleMsg(m *Message, p *peer){
	// fmt.Printf("Peer : %s, Sent a message with kind of : %d\n",p.key,m.Kind)
	switch m.Kind {
		case MessageNewestBlock:
			fmt.Printf("Received the newest block from %s\n",p.key)
			var payload blockchain.Block			
			utils.HandleError(json.Unmarshal(m.Payload, &payload))
			b, err:= blockchain.FindBlock(blockchain.BlockChain().NewestHash)
			utils.HandleError(err)
			
			if payload.Height >= b.Height{
				fmt.Printf("Requesting all blocks from %s\n",p.key)
				requestAllBlocks(p)
				// port 3000이 먼저 실행되는것부터 생각 
				// 위의 파라미터 p는 포트 4000임.
				// request all the blocks from 4000
			} else {
				// send 4000 our block
				fmt.Printf("Sending newest block to %s\n",p.key)
				sendNewestBlock(p)
			}
		case MessageAllBlocksRequest:
			fmt.Printf("%s wants all the blocks\n",p.key)
			sendAllBlocks(p)
		case MessageAllBlocksResponse:
			fmt.Printf("Received all the blocks from %s\n",p.key)
			var payload []*blockchain.Block			
			utils.HandleError(json.Unmarshal(m.Payload,&payload))
			blockchain.BlockChain().Replace(payload)
		case MessageNewBlockNotify:
			var payload *blockchain.Block			
			utils.HandleError(json.Unmarshal(m.Payload,&payload))
			blockchain.BlockChain().AddPeerBlock(payload)
		case MessageNewTxNotify:
			var payload *blockchain.Tx
			utils.HandleError(json.Unmarshal(m.Payload,&payload))
			blockchain.Mempool().AddPeerTx(payload)
			// blockchain.BlockChain().AddPeerTx(payload)	
		case MessageNewPeerNotify:
			var payload string
			utils.HandleError(json.Unmarshal(m.Payload,&payload))
			fmt.Println(payload)
			parts := strings.Split(payload, ":")
			AddPeer(parts[0], parts[1],parts[2],false)

	}
}
