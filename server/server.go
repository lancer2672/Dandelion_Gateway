package server

import (
	"log"
	"net/http"

	"github.com/lancer2672/Dandelion_Gateway/internal/utils"
)

func RunServer() {

	http.HandleFunc("/ws", handleWebSocket)
	r := setupRouter()
	http.Handle("/", r)
	http.ListenAndServe(utils.ConfigIns.GatewayAddress, nil)
	log.Println("Server started at:", utils.ConfigIns.GatewayAddress)
}

// func rewriteBody(resp *http.Response) (err error) {
// 	b, err := io.ReadAll(resp.Body) //Read html
// 	if err != nil {
// 		return err
// 	}
// 	err = resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("RewriteBody", b)
// 	// b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1) // replace html
// 	// body := io.NopCloser(bytes.NewReader(b))
// 	// resp.Body = body
// 	// resp.ContentLength = int64(len(b))
// 	// resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
// 	return nil
// }
