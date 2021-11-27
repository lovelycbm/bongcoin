package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/lovelycbm/bongcoin/blockchain"
)

const port string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// express에 response.send를 이런 방식으로 한다면 이해가 되는듯?
	// fmt.Fprint(rw, "Hello from home")
	tmpl := template.Must(template.ParseFiles("templates/pages/home.gohtml"))

	data := homeData{"bong coin home", blockchain.GetBlockChain().AllBlocks()}
	tmpl.Execute(rw, data)

}

func main() {
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
	// chain := blockchain.GetBlockChain()
	// chain.AddBlock("Second Block")
	// chain.AddBlock("Third Block")
	// chain.AddBlock("Fourth Block")
	// for _, block := range chain.AllBlocks() {
	// 	fmt.Printf("Data : %s\n", block.Data)
	// 	fmt.Printf("Hash : %s\n", block.Hash)
	// 	fmt.Printf("Prev Hash : %s\n", block.PrevHash)
	// }
}
