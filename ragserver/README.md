# ragserver

*RAG stands for Retrieval Augmented Generation*

Demos of implementing a "RAG Server" in Go, using [Google
AI](https://ai.google.dev/) for embeddings and language models and
[Weaviate](https://weaviate.io/) as a vector database.


## How it works

The server we're developing is a standard Go HTTP server, listening on a local
port. See the next section for the request schema for this server. It supports
adding new documents to its context, and getting queries that would use this
context.

Weaviate has the be installed locally; the easiest way to do so is by using
`docker-compose` as described in the Usage section.

## Server request schema

```
/add/: POST {"documents": [{"text": "..."}, {"text": "..."}, ...]}
  response: OK (no body)

/query/: POST {"content": "..."}
  response: model response as a string
```

## Server variants

* `ragserver`: uses the Google AI Go SDK directly for LLM calls and embeddings,
  and the Weaviate Go client library directly for interacting with Weaviate.
* `ragserver-langchaingo`: uses [LangChain for Go](https://github.com/tmc/langchaingo)
  to interact with Weaviate and Google's LLM and embedding models.
* `ragserver-genkit`: uses [Genkit Go](https://firebase.google.com/docs/genkit-go/get-started-go)
  to interact with Weaviate and Google's LLM and embedding models.

## Usage

* In terminal window 1, `cd tests` and run `docker-compose up`;
  This will start the weaviate service in the foreground.
* In terminal window 2, run `GEMINI_API_KEY=... go run .` in the tested
  `ragserver` directory.
* In terminal window 3, we can now run scripts to clear/populate the
  weaviate DB and interact with `ragserver`. The following instructions are
  for terminal window 3.

Run `cd tests`; then we can clear out the weaviate DB with
`./weaviate-delete-objects.sh`. To add documents to the DB through `ragserver`,
run `./add-documents.sh`. For a sample query, run `./query.sh`
Adjust the contents of these scripts as needed.

## Environment variables

* `SERVERPORT`: the port this server is listening on (default 9020)
* `WVPORT`: the port Weaviate is listening on (default 9035)
* `GEMINI_API_KEY`: API key for the Gemini service at https://ai.google.dev
