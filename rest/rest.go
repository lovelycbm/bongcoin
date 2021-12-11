package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovelycbm/bongcoin/blockchain"
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


func blocks(rw http.ResponseWriter, r *http.Request) {	
	switch r.Method {
	case "GET": 		
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":					
		blockchain.BlockChain().AddBlock()		
		rw.WriteHeader(http.StatusCreated)
	}
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
		
	}
	fmt.Println(data)
	// rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
	
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

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)		
	address := vars["address"]
	// url에 포함되는 쿼리 받는 법
	total := r.URL.Query().Get("total")
	// if total == "true" {
	// 	err := json.NewEncoder(rw).Encode(blockchain.BlockChain().BalanceByAddress(address))
	// 	utils.HandleError(err)	
	// } else {
	
	// }
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

func mempool(rw http.ResponseWriter, r *http.Request) {	
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {	
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err:= blockchain.Mempool.AddTx(payload.To, payload.Amount)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return 
	}
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request){
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address:address})
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d",aPort)
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}