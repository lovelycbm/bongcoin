package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovelycbm/bongcoin/blockchain"
	"github.com/lovelycbm/bongcoin/p2p"
	"github.com/lovelycbm/bongcoin/utils"
	"github.com/lovelycbm/bongcoin/wallet"
)

var port string
// const port string = ":4000"
// const url string ="http://localhost"
type url string

func (u url) MarshalText() ([]byte,error){
	url:= fmt.Sprintf("http://localhost%s%s",port,u)
	return []byte(url),nil
}

type urlDescription struct {
	URL url `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string  `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int `json:"balance"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To  string 
	Amount int 
}

type addPeerPayload struct{
	Address, Port string
}



func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All blocks",			
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See Status of the blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get Txouts for an address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
		
	}
	fmt.Println(data)
	// rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
	
}

func blocks(rw http.ResponseWriter, r *http.Request) {	
	switch r.Method {
	case "GET": 		
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":					
		newBlock := blockchain.BlockChain().AddBlock()		
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
} 

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)		
	hash := vars["hash"]
	
	block, err:= blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	
	// utils.HandleError(err)
	if err == blockchain.ErrBlockNotFound { 
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	// json.NewEncoder(rw).Encode(blockchain.BlockChain())
	blockchain.Status(blockchain.BlockChain(),rw)
}

func mempool(rw http.ResponseWriter, r *http.Request) {		
	blockchain.MemStatus(blockchain.Mempool(),rw)
	// blockchain.MempoolReplace(blockchain.Mempool().Txs,rw)
	// utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))

}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)		
	address := vars["address"]
	// url에 포함되는 쿼리 받는 법
	total := r.URL.Query().Get("total")	
		
	switch total {
		case "true": 
			amount := blockchain.BalanceByAddress(address,blockchain.BlockChain())			
			err := json.NewEncoder(rw).Encode(balanceResponse{address,amount})
			utils.HandleError(err)
		default:
			err := json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address,blockchain.BlockChain()))
			utils.HandleError(err)
	}
	
}


func transactions(rw http.ResponseWriter, r *http.Request) {	
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	tx, err:= blockchain.Mempool().AddTx(payload.To, payload.Amount)

	
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return 
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request){
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address:address})
}

func peers(rw http.ResponseWriter, r *http.Request){
	switch r.Method {
		
		case "POST":
			var payload addPeerPayload
			utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
			p2p.AddPeer(payload.Address, payload.Port,port[1:],true)
			rw.WriteHeader(http.StatusOK)
		case "GET":
			json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))

	}
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d",aPort)
	router.Use(jsonContentTypeMiddleWare,loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET","POST")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}