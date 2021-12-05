package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lovelycbm/bongcoin/blockchain"
	"github.com/lovelycbm/bongcoin/utils"
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

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}



func blocks(rw http.ResponseWriter, r *http.Request) {	
	switch r.Method {
	case "GET": 
		// rw.Header().Add("Content-Type","application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))		
		//fmt.Println(addBlockBody)
		blockchain.GetBlockChain().AddBlock(addBlockBody.Message)
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
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A block",
		},
	}
	fmt.Println(data)
	// rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
	
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars)
	height := vars["height"]
	// id convert string to int 
	id, err := strconv.Atoi(height)
	utils.HandleError(err)

	encoder := json.NewEncoder(rw)

	block, err:= blockchain.GetBlockChain().GetBlock(id)
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

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d",aPort)
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}