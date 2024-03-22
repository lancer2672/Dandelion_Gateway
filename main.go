package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lancer2672/Dandelion_Gateway/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("URL PATH", r.URL.Path)
		switch {
		case strings.HasPrefix(r.URL.Path, "/notification"):
			http.Redirect(w, r, config.NotificationServiceAddress+r.URL.Path, http.StatusFound)
		case strings.HasPrefix(r.URL.Path, "/movie"):
			http.Redirect(w, r, config.MovieGRPCAddress+r.URL.Path, http.StatusFound)
		default:
			http.Redirect(w, r, config.MainServiceAddress+r.URL.Path, http.StatusFound)
		}
	})

	log.Fatal(http.ListenAndServe(config.GatewayAddress, nil))
}
