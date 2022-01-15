package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHTTPServer(address string) *http.Server {
	httpServer := newHTTPServer()
	router := mux.NewRouter()
	router.HandleFunc("/", httpServer.handleProduce).Methods("POST")
	router.HandleFunc("/", httpServer.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var produceRequest ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&produceRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := s.Log.Append(produceRequest.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	produceResponse := ProduceResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(produceResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var consumeRequest ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&consumeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	record, err := s.Log.Read(consumeRequest.Offset)
	if err == ErrOffSetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	consumeResponse := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(consumeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
