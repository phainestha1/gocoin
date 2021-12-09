package cli

import (
	"flag"
	"fmt"
	"gocoin/explorer"
	"gocoin/rest"
	"os"
	"runtime"
)


func usage() {
	fmt.Printf("🔥 Golang Practice 🔥 \n\n")
	fmt.Printf("Please use the following commands:\n\n")
	fmt.Printf("-port: Set the port of the server\n")
	fmt.Printf("-mode: hoose 'html' or 'rest'\n\n")
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	hport := flag.Int("hport", 4000, "Set port of the html server")
	rport := flag.Int("rport", 3000, "Set port of the rest server")
	mode := flag.String("mode", "run", "Choose between 'html' and 'rest' or 'run' for both")
	
	flag.Parse()
	
	switch *mode {
	case "run":
		go rest.Start(*rport)
		explorer.Start(*hport)
	case "rest":
		rest.Start(*rport)
	case "html":
		explorer.Start(*hport)
	default:
		usage()
	}
}