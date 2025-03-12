package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a JSON RPC response
		response := RPCResponse{
			Jsonrpc: "2.0",
			ID:      2,
			Result:  "0x123456",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Override the RPC URL to use the mock server
	rpcURL = server.URL

	// Create a request to the /blocknumber endpoint
	req, err := http.NewRequest("GET", "/blocknumber", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getBlockNumber)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response body
	var actualResponse RPCResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Define the expected response
	expectedResponse := RPCResponse{
		Jsonrpc: "2.0",
		ID:      2,
		Result:  "0x123456",
	}

	// Compare the actual and expected responses
	if actualResponse.Jsonrpc != expectedResponse.Jsonrpc ||
		actualResponse.ID != expectedResponse.ID ||
		actualResponse.Result != expectedResponse.Result {
		t.Errorf("handler returned unexpected body: got %+v want %+v", actualResponse, expectedResponse)
	}
}

func TestGetBlockByNumber(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a JSON RPC response
		response := RPCResponse{
			Jsonrpc: "2.0",
			ID:      2,
			Result: map[string]interface{}{
				"number": "0x134e82a",
				"hash":   "0xabc123",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Override the RPC URL to use the mock server
	rpcURL = server.URL

	// Create a request to the /block endpoint
	req, err := http.NewRequest("GET", "/block?number=0x134e82a", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getBlockByNumber)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response body
	var actualResponse RPCResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Define the expected response
	expectedResponse := RPCResponse{
		Jsonrpc: "2.0",
		ID:      2,
		Result: map[string]interface{}{
			"number": "0x134e82a",
			"hash":   "0xabc123",
		},
	}

	// Compare the actual and expected responses
	if actualResponse.Jsonrpc != expectedResponse.Jsonrpc ||
		actualResponse.ID != expectedResponse.ID ||
		actualResponse.Result.(map[string]interface{})["number"] != expectedResponse.Result.(map[string]interface{})["number"] ||
		actualResponse.Result.(map[string]interface{})["hash"] != expectedResponse.Result.(map[string]interface{})["hash"] {
		t.Errorf("handler returned unexpected body: got %+v want %+v", actualResponse, expectedResponse)
	}
}