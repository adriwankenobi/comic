package api

import (
    "net/http"
    "github.com/adriwankenobi/comic/service"
)

const data = "data.json"

func init() {
    http.HandleFunc("/api", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
    	w.WriteHeader(http.StatusForbidden)
    	return
	}
    param := r.URL.Query()["q"]
    q := "id!=''" // default search, will return everything
    if len(param) > 0 {
	    q = param[0]
    }
    result, err := service.FindComicList(data, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := result.ToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Query", q)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
