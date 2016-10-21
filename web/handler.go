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

type elementHandler func(json *jsonql.JSONQL, id string) (jsonAble, error)
type listHandler func(json *jsonql.JSONQL) (jsonAble, error)

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
    router.GET("/api/comics", listHandle(listComics, comics))
    router.GET("/api/comics/:id", elementHandle(getComic, comics))
    router.GET("/api/phases", listHandle(listPhases, phases))
    router.GET("/api/phases/:id", elementHandle(getPhase, phases))
    router.GET("/api/issues", listHandle(listIssuesPhases, issuesPhases))
    router.GET("/api/issues/:id", elementHandle(getIssuesPhase, issuesPhases))
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

func listHandle(handle listHandler, json *jsonql.JSONQL) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result, err := handle(json)
		writeResponse(w, result, err)
	}
}

func elementHandle(handle elementHandler, json *jsonql.JSONQL) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		result, err := handle(json, id)
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

func listComics(comics *jsonql.JSONQL) (jsonAble, error) {
	return service.ListComics(comics)
}

func getComic(comics *jsonql.JSONQL, id string) (jsonAble, error) {
	return service.FindComic(comics, id)
}

func listPhases(phases *jsonql.JSONQL) (jsonAble, error) {
	return service.ListPhases(phases)
}

func getPhase(phases *jsonql.JSONQL, id string) (jsonAble, error) {
	return service.FindPhase(phases, id)
}

func listIssuesPhases(issuesPhases *jsonql.JSONQL) (jsonAble, error) {
	return service.ListIssuesPhases(issuesPhases)
}

func getIssuesPhase(issuesPhases *jsonql.JSONQL, id string) (jsonAble, error) {
	return service.FindIssuesPhase(issuesPhases, id)
}
