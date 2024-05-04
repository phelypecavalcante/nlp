package main

import (
	"encoding/json"
	"expvar"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/phelypecavalcante/nlp"
	stemmer "github.com/phelypecavalcante/nlp/stremmer"
)

var (
	numTok = expvar.NewInt("tokenize.calls")
)

func main() {
	// Create server
	logger := log.New(log.Writer(), "[nlpd] ", log.Flags()|log.Lshortfile)
	s := Server{
		logger: logger, // dependency injection
	}
	// routing
	// /health is an exact match
	// /health/ is a prefix match
	/*
		http.HandleFunc("/health", healthHandler)
		http.HandleFunc("/tokenize", tokenizeHandler)
	*/

	r := mux.NewRouter()
	r.HandleFunc("/heath", s.healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/tokenize", s.tokenizeHandler).Methods(http.MethodPost)
	r.HandleFunc("/stem/{word}", s.stemHandler).Methods(http.MethodGet)

	http.Handle("/", r)
	// run server
	addr := os.Getenv("NLPD_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	s.logger.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("error: %s", err)
	}

}

type Server struct {
	logger *log.Logger
}

func (s *Server) stemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	stem := stemmer.Stem(word)
	fmt.Fprintln(w, stem)
}

// exercise: Write a tokenizeHandler that will read the text from the request
// body and return JSON in the format "{"tokens": ["who", "on", "first"]}"

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Run a health check
	fmt.Fprintln(w, "OK")
}

func (s *Server) tokenizeHandler(w http.ResponseWriter, r *http.Request) {
	/* Before Gorilla
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}
	*/

	numTok.Add(1)

	// step 1: Get, conver & validate the data
	defer r.Body.Close()
	rdr := io.LimitReader(r.Body, 1_000_000)
	data, err := io.ReadAll(rdr)
	if err != nil {
		http.Error(w, "can't read", http.StatusBadRequest)
		return
	}

	if len(data) == 0 {
		s.logger.Printf("error: can't read - %s", err)
		http.Error(w, "missing data", http.StatusBadRequest)
		return
	}

	text := string(data)
	// Step 2: Work
	tokens := nlp.Tokenize(text)

	// Step 3:  Encode & emit output
	resp := map[string]any{
		"tokens": tokens,
	}
	data, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, "can't encode", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
