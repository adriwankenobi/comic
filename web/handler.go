package api

import (
	"fmt"
	"io/ioutil"
    "net/http"
    "github.com/elgs/jsonql"
    "github.com/julienschmidt/httprouter"
    "github.com/adriwankenobi/comic/service"
)

const (
	ComicsFile = "comics.json"
	PhasesFile = "phases.json"
	IssuesPhasesFile = "issues-phases.json"
)

type jsonAble interface {
    ToJson() ([]byte, error)
    IsEmpty() bool
}

type handler func(json *jsonql.JSONQL, p httprouter.Params) (jsonAble, error)

func init() {
	// Read files
	comics, err := getJson(ComicsFile)
    if err != nil {
        return
    }
    phases, err := getJson(PhasesFile)
    if err != nil {
        return
    }
    issuesPhases, err := getJson(IssuesPhasesFile)
    if err != nil {
        return
    }
	
	// Start server
	router := httprouter.New()
	
	// Get all comics
    router.GET("/api/comics", jsonHandle(func(comics *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.ListComics(comics)
	}, comics))
    
    // Get this comic
    router.GET("/api/comics/:id", jsonHandle(func(comics *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.FindComicByID(comics, p.ByName("id"))
	}, comics))
    
    // Get all phases
    router.GET("/api/phases", jsonHandle(func(phases *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.ListPhases(phases)
	}, phases))
    
    // Get this phase
    router.GET("/api/phases/:id", jsonHandle(func(phases *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.FindPhaseByID(phases, p.ByName("id"))
	}, phases))
    
    // Get all first issues from all phases
    router.GET("/api/fissues", jsonHandle(func(issuesPhases *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.ListFirstIssues(issuesPhases)
	}, issuesPhases))
    
    // Get all first issues from this phase
    router.GET("/api/fissues/:id", jsonHandle(func(issuesPhases *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.FindFirstIssuesByPhaseID(issuesPhases, p.ByName("id"))
	}, issuesPhases))
    
    // Get all issues from this phase
    router.GET("/api/phases/:id/issues", jsonHandle(func(comics *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.ListComicsByPhaseID(comics, p.ByName("id"))
	}, comics))
    
    // Get all issues from this comic from this phase
    router.GET("/api/phases/:id/issues/:sortid", jsonHandle(func(comics *jsonql.JSONQL, p httprouter.Params) (jsonAble, error) {
		return service.ListComicsByPhaseAndSortIDs(comics, p.ByName("id"), p.ByName("sortid"))
	}, comics))
    
    http.Handle("/", router)
}

func getJson(file string) (*jsonql.JSONQL, error) {
	content, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, fmt.Errorf("Error reading file '%s'", file)
    }
    json, err := jsonql.NewStringQuery(string(content))
    if err != nil {
        return nil, fmt.Errorf("Error parsing content from file '%s'", file)
    }
    return json, nil
}

func jsonHandle(handle handler, json *jsonql.JSONQL) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result, err := handle(json, p)
		writeResponse(w, result, err)
	}
}

func writeResponse(w http.ResponseWriter, j jsonAble, err error) {
	if err != nil {
		w.Header().Set("Error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if j.IsEmpty() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bytes, err := j.ToJson()
	if err != nil {
		w.Header().Set("Error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
