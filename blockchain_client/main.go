package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var rpcURL = "https://polygon-rpc.com/"

type RPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func makeRPCRequest(client *http.Client, method string, params []interface{}) ([]byte, error) {
	request := RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      2,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(rpcURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func getBlockNumber(w http.ResponseWriter, r *http.Request) {
	response, err := makeRPCRequest(http.DefaultClient, "eth_blockNumber", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(response, &rpcResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rpcResponse)
}

func getBlockByNumber(w http.ResponseWriter, r *http.Request) {
	blockNumber := r.URL.Query().Get("number")
	if blockNumber == "" {
		http.Error(w, "Missing block number", http.StatusBadRequest)
		return
	}

	response, err := makeRPCRequest(http.DefaultClient, "eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(response, &rpcResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rpcResponse)
}


func main() {
	http.HandleFunc("/blocknumber", getBlockNumber)
	http.HandleFunc("/block", getBlockByNumber)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}