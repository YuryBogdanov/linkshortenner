package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const urlParameterShortenedID = "id"

func handleNewLinkRegistration(w http.ResponseWriter, r *http.Request) {
	if url, err := io.ReadAll(r.Body); err == nil {
		linkID, err := makeAndStoreShortURL(string(url))
		if err != nil {
			handleError(w)
			return
		}
		resultLink := getShortenedLink(r, linkID)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resultLink))
	} else {
		handleError(w)
	}
}

func handleExistingLinkRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path[1:]
	// query := chi.URLParam(r, urlParameterShortenedID)
	if link, exists := links[query]; exists {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		handleError(w)
	}
}

func getShortenedLink(r *http.Request, linkID string) string {
	return "http://" + r.Host + "/" + linkID
}

func handleError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	router := chi.NewRouter()
	router.Get("/{id}", handleExistingLinkRequest)
	router.Post("/", handleNewLinkRegistration)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
