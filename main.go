package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/shipyardchain/hama_hyperledger-fabric-400/hyperledger"
	mux "github.com/shipyardchain/hama_mux"
)

type Message struct {
	key   string
	value string
}

// Start & Test
func main() {
	hyperledger.StartFabric()

	hyperledger.WriteTrans("1", "bitcoin BTC111111111")
	hyperledger.WriteTrans("2", "eos")
	hyperledger.WriteTrans("3", "hyperledger HPL")
	hyperledger.WriteTrans("4", "ethereum ETH")

	time.Sleep(1 * time.Second)

	result1 := hyperledger.GetTrans("1")
	result2 := hyperledger.GetTrans("2")
	result3 := hyperledger.GetTrans("3")
	result4 := hyperledger.GetTrans("4")

	fmt.Printf("key1 : %s \n", result1)
	fmt.Printf("key2 : %s \n", result4)
	fmt.Printf("key3 : %s \n", result3)
	fmt.Printf("key4 : %s \n", result2)

	//web_server_run()
}

func web_server_run() error {
	mux := makeMuxRouter()
	httpPort := "8080"
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", getLetter).Methods("GET")
	muxRouter.HandleFunc("/", writeLetter).Methods("POST")
	return muxRouter
}

func getLetter(w http.ResponseWriter, r *http.Request) {

	letters := hyperledger.GetTrans("1")
	bytes, err := json.MarshalIndent(letters, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func writeLetter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var msg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	result := hyperledger.WriteTrans(msg.key, msg.value)
	respondWithJSON(w, r, http.StatusCreated, result)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
