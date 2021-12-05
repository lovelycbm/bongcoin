package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/lovelycbm/bongcoin/explorer"
	"github.com/lovelycbm/bongcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to Bongcoin\n\n")
	fmt.Printf("Please use following flags: \n\n")
	fmt.Printf("-port:	Set port of the server\n")
	fmt.Printf("-mode:	Set mode of the server ('rest' or 'html')\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Set mode of the server ('rest' or 'html')")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}

}