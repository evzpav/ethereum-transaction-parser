package api

import (
	"encoding/json"
	"ethereum-parser/models"
	"fmt"
	"io"
	"log"
	"net/http"
)

type api struct {
	port   int
	parser models.Parser
}

func New(port int, parser models.Parser) *api {
	a := &api{
		port:   port,
		parser: parser,
	}

	a.Routes()

	return a
}

func (a *api) Serve() error {
	log.Printf("Server started at %d", a.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", a.port), nil)
}

func (a *api) Routes() {
	http.HandleFunc("/transactions", a.GetTransactions)
	http.HandleFunc("/subscribe", a.Subscribe)
	http.HandleFunc("/current-block", a.GetCurrentBlock)
}

func (a *api) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res := map[string]interface{}{
		"currentBlock": a.parser.GetCurrentBlock(),
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}

}

func (a *api) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	address := r.URL.Query().Get("address")

	if address == "" {
		http.Error(w, "address query param is required", http.StatusBadRequest)
		return
	}
	res := map[string]interface{}{
		"transactions": a.parser.GetTransactions(address),
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}

}

func (a *api) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusInternalServerError)
		return
	}

	type SubscribeRequest struct {
		Address string `json:"address"`
	}

	var sub SubscribeRequest

	err = json.Unmarshal(bs, &sub)
	if err != nil {
		http.Error(w, "error unmarshaling request body", http.StatusInternalServerError)
		return
	}

	if sub.Address == "" {
		http.Error(w, "address query param is required", http.StatusBadRequest)
		return
	}

	if a.parser.Subscribe(sub.Address) {

		res := map[string]interface{}{
			"message": fmt.Sprintf("%s subscribed", sub.Address),
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "already subscribed", http.StatusBadRequest)
}
