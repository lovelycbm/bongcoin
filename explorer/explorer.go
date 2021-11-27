package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/lovelycbm/bongcoin/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// express에 response.send를 이런 방식으로 한다면 이해가 되는듯?
	// fmt.Fprint(rw, "Hello from home")
	data := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}

}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "/pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "/partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
