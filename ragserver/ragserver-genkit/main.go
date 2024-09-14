// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command ragserver is an HTTP server that implements RAG (Retrieval
// Augmented Generation) using the Gemini model and Weaviate, which
// are accessed using the Genkit package. See the accompanying README file for
// additional details.
package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/googleai"
	"github.com/firebase/genkit/go/plugins/weaviate"
)

const generativeModelName = "gemini-1.5-flash"
const embeddingModelName = "text-embedding-004"

// This is a standard Go HTTP server. Server state is in the ragServer struct.
// The `main` function connects to the required services (Weaviate and Google
// AI), initializes the server state and registers HTTP handlers.
func main() {
	ctx := context.Background()
	err := googleai.Init(ctx, &googleai.Config{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	wvConfig := &weaviate.ClientConfig{
		Scheme: "http",
		Addr:   "localhost:" + cmp.Or(os.Getenv("WVPORT"), "9035"),
	}
	_, err = weaviate.Init(ctx, wvConfig)
	if err != nil {
		log.Fatal(err)
	}

	classConfig := &weaviate.ClassConfig{
		Class:    "Document",
		Embedder: googleai.Embedder(embeddingModelName),
	}
	indexer, retriever, err := weaviate.DefineIndexerAndRetriever(ctx, *classConfig)
	if err != nil {
		log.Fatal(err)
	}

	model := googleai.Model(generativeModelName)
	if model == nil {
		log.Fatal("unable to set up gemini-1.5-flash model")
	}

	server := &ragServer{
		ctx:       ctx,
		indexer:   indexer,
		retriever: retriever,
		model:     model,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /add/", server.addDocumentsHandler)
	mux.HandleFunc("POST /query/", server.queryHandler)

	port := cmp.Or(os.Getenv("SERVERPORT"), "9020")
	address := "localhost:" + port
	log.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, mux))
}

type ragServer struct {
	ctx       context.Context
	indexer   ai.Indexer
	retriever ai.Retriever
	model     ai.Model
}

func (rs *ragServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type document struct {
		Text string
	}
	type addRequest struct {
		Documents []document
	}
	ar := &addRequest{}
	err := readRequestJSON(req, ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert request documents into Weaviate documents used for embedding.
	var wvDocs []*ai.Document
	for _, doc := range ar.Documents {
		wvDocs = append(wvDocs, ai.DocumentFromText(doc.Text, nil))
	}

	// Index the requested documents.
	err = ai.Index(rs.ctx, rs.indexer, ai.WithIndexerDocs(wvDocs...))
	if err != nil {
		http.Error(w, fmt.Errorf("indexing: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (rs *ragServer) queryHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type queryRequest struct {
		Content string
	}
	qr := &queryRequest{}
	err := readRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the most similar documents using the retriever.
	resp, err := ai.Retrieve(rs.ctx,
		rs.retriever,
		ai.WithRetrieverDoc(ai.DocumentFromText(qr.Content, nil)),
		ai.WithRetrieverOpts(&weaviate.RetrieverOptions{
			Count: 3,
		}))
	if err != nil {
		http.Error(w, fmt.Errorf("retrieval: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	var docsContents []string
	for _, d := range resp.Documents {
		docsContents = append(docsContents, d.Content[0].Text)
	}

	// Create a RAG query for the LLM with the most relevant documents as
	// context.
	ragQuery := fmt.Sprintf(ragTemplateStr, qr.Content, strings.Join(docsContents, "\n"))
	genResp, err := ai.Generate(rs.ctx, rs.model, ai.WithTextPrompt(ragQuery))
	if err != nil {
		log.Printf("calling generative model: %v", err.Error())
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	if len(genResp.Candidates) != 1 {
		log.Printf("got %v candidates, expected 1", len(genResp.Candidates))
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	renderJSON(w, genResp.Text())
}

const ragTemplateStr = `
I will ask you a question and will provide some additional context information.
Assume this context information is factual and correct, as part of internal
documentation.
If the question relates to the context, answer it using the context.
If the question does not relate to the context, answer it as normal.

For example, let's say the context has nothing in it about tropical flowers;
then if I ask you about tropical flowers, just answer what you know about them
without referring to the context.

For example, if the context does mention minerology and I ask you about that,
provide information from the context along with general knowledge.

Question:
%s

Context:
%s
`
