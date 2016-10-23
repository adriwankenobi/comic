package web

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"github.com/elgs/jsonql"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

var htmlFiles = []string{
	"template.html",
	"tabs.html",
	"tab-li.html",
	"tab-content.html",
	"tab-content-intro.html",
	"tab-content-phase.html",
	"clear-fix.html",
	"issues.html",
	"issue-content.html",
	"not-found.html",
}

type jsonHandler func(p httprouter.Params) (service.JsonAble, error)
type webHandler func(p httprouter.Params) (string, error)

type webContent map[string]string
type jsonContent map[string]*jsonql.JSONQL

var c webContent
var j jsonContent

func init() {
	// Read files
	j, err := readJsonFiles(service.Datastore)
	if err != nil {
		return
	}

	c, err = readWebFiles(htmlFiles)
	if err != nil {
		return
	}

	// Start server
	router := httprouter.New()

	// API

	// Get all comics
	router.GET("/api/comics", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.ListComics(j["comics"])
	}))

	// Get this comic
	router.GET("/api/comics/:id", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.FindComicByID(j["comics"], p.ByName("id"))
	}))

	// Get all phases
	router.GET("/api/phases", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.ListPhases(j["phases"])
	}))

	// Get this phase
	router.GET("/api/phases/:id", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.FindPhaseByID(j["phases"], p.ByName("id"))
	}))

	// Get all first issues from all phases
	router.GET("/api/fissues", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.ListFirstIssues(j["fissues"])
	}))

	// Get all first issues from this phase
	router.GET("/api/fissues/:id", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.FindFirstIssuesByPhaseID(j["fissues"], p.ByName("id"))
	}))

	// Get all issues from this phase
	router.GET("/api/phases/:id/issues", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.ListComicsByPhaseID(j["comics"], p.ByName("id"))
	}))

	// Get all issues from this comic from this phase
	router.GET("/api/phases/:id/issues/:sortid", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.ListComicsByPhaseAndSortIDs(j["comics"], p.ByName("id"), p.ByName("sortid"))
	}))

	// WEB

	// Index -> Get all first issues from all phases
	router.GET("/", webHandle(func(p httprouter.Params) (string, error) {
		issues, err := service.ListFirstIssues(j["fissues"])
		if err != nil {
			return "", err
		}
		return getIndexPage(issues)
	}))

	// Issues -> Get all issues from this comic from this phase
	router.GET("/phases/:id/issues/:sortid", webHandle(func(p httprouter.Params) (string, error) {
		issues, err := service.ListComicsByPhaseAndSortIDs(j["comics"], p.ByName("id"), p.ByName("sortid"))
		if err != nil {
			return "", err
		}
		return getIssuesPage(issues)
	}))

	http.Handle("/", router)
}

// File readers
func readJsonFiles(in service.DatastoreType) (jsonContent, error) {
	m := make(jsonContent)
	for key, _ := range in {
		bytes, err := ioutil.ReadFile(fmt.Sprintf("%s.json", key))
		if err != nil {
			return m, err
		}
		json, err := jsonql.NewStringQuery(string(bytes))
		if err != nil {
			return m, err
		}
		m[key] = json
	}
	return m, nil
}

func readWebFiles(in []string) (webContent, error) {
	m := make(webContent)
	for _, e := range in {
		bytes, err := ioutil.ReadFile(e)
		if err != nil {
			return m, err
		}
		tag := strings.Split(e, ".")[0]
		m[tag] = string(bytes)
	}
	return m, nil
}

// Handlers
func jsonHandle(handle jsonHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result, err := handle(p)
		writeJsonResponse(w, result, err)
	}
}

func webHandle(handle webHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result, err := handle(p)
		writeResponse(w, result, err)
	}
}

// Response Writers
func writeJsonResponse(w http.ResponseWriter, j service.JsonAble, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	if j.IsEmpty() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bytes, err := j.ToJson()
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func writeResponse(w http.ResponseWriter, s string, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	bytes := []byte(s)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Error", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}
