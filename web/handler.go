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

type jsonHandler func(p httprouter.Params) (service.JsonAble, error)
type webHandler func(p httprouter.Params) (string, error)

type webContent map[string]string
type jsonContent map[string]*jsonql.JSONQL

var c webContent
var j jsonContent

func init() {
	// Read files
	j, err := readJsonFiles()
	if err != nil {
		return
	}
	c, err = readWebFiles()
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
		return service.ListComics(j[fmt.Sprintf("comics-phase-%s", p.ByName("id"))])
	}))

	// Get all issues from this comic from this phase
	router.GET("/api/phases/:id/issues/:sortid", jsonHandle(func(p httprouter.Params) (service.JsonAble, error) {
		return service.ListComicsBySortID(j[fmt.Sprintf("comics-phase-%s", p.ByName("id"))], p.ByName("sortid"))
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
		issues, err := service.ListComicsBySortID(j[fmt.Sprintf("comics-phase-%s", p.ByName("id"))], p.ByName("sortid"))
		if err != nil {
			return "", err
		}
		return getIssuesPage(issues)
	}))
	
	// About
	router.GET("/about", webHandle(func(p httprouter.Params) (string, error) {
		return getAboutPage(), nil
	}))

	http.Handle("/", router)
}

// File readers
func readJsonFiles() (jsonContent, error) {
	m := make(jsonContent)
	folder := "data"
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return m, err
	}
	for _, f := range files {
		split := strings.Split(f.Name(), ".")
		if f.IsDir() || len(split) != 2 || split[1] != "json" {
			break
		}
		bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folder, f.Name()))
		if err != nil {
			return m, err
		}
		json, err := jsonql.NewStringQuery(string(bytes))
		if err != nil {
			return m, err
		}
		m[split[0]] = json
	}
	return m, nil
}

func readWebFiles() (webContent, error) {
	m := make(webContent)
	folder := "html"
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return m, err
	}
	for _, f := range files {
		split := strings.Split(f.Name(), ".")
		if f.IsDir() || len(split) != 2 || split[1] != "html" {
			break
		}
		bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folder, f.Name()))
		if err != nil {
			return m, err
		}
		m[split[0]] = string(bytes)
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
