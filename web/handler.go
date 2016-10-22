package api

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"github.com/elgs/jsonql"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

type jsonAble interface {
	ToJson() ([]byte, error)
	IsEmpty() bool
}

type jsonHandler func(p httprouter.Params) (jsonAble, error)
type webHandler func(p httprouter.Params) (string, error)

type webContent map[string]string
type jsonContent map[string]*jsonql.JSONQL

var c webContent
var j jsonContent

func init() {
	// Read files
	j, err := readJsonFiles([]string{
		"comics.json",
		"phases.json",
		"issues-phases.json",
	})
	if err != nil {
		return
	}
	
	c, err = readWebFiles([]string{
		"template.html",
		"tabs.html",
		"tab-li.html",
		"tab-content.html",
		"tab-content-intro.html",
		"tab-content-phase.html",
		"clear-fix.html",
		"issues.html",
		"issue-content.html",
	})
	if err != nil {
		return
	}

	// Start server
	router := httprouter.New()

	// API

	// Get all comics
	router.GET("/api/comics", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.ListComics(j["comics"])
	}))

	// Get this comic
	router.GET("/api/comics/:id", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.FindComicByID(j["comics"], p.ByName("id"))
	}))

	// Get all phases
	router.GET("/api/phases", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.ListPhases(j["phases"])
	}))

	// Get this phase
	router.GET("/api/phases/:id", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.FindPhaseByID(j["phases"], p.ByName("id"))
	}))

	// Get all first issues from all phases
	router.GET("/api/fissues", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.ListFirstIssues(j["issues-phases"])
	}))

	// Get all first issues from this phase
	router.GET("/api/fissues/:id", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.FindFirstIssuesByPhaseID(j["issues-phases"], p.ByName("id"))
	}))

	// Get all issues from this phase
	router.GET("/api/phases/:id/issues", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.ListComicsByPhaseID(j["comics"], p.ByName("id"))
	}))

	// Get all issues from this comic from this phase
	router.GET("/api/phases/:id/issues/:sortid", jsonHandle(func(p httprouter.Params) (jsonAble, error) {
		return service.ListComicsByPhaseAndSortIDs(j["comcis"], p.ByName("id"), p.ByName("sortid"))
	}))

	// WEB

	// Index -> Get all first issues from all phases
	router.GET("/", webHandle(func(p httprouter.Params) (string, error) {
		return getIndexPage(j["issues-phases"])
	}))

	// Issues -> Get all issues from this comic from this phase
	router.GET("/phases/:id/issues/:sortid", webHandle(func(p httprouter.Params) (string, error) {
		return getIssuesPage(j["comics"], p.ByName("id"), p.ByName("sortid"))
	}))

	http.Handle("/", router)
}

// File readers
func readJsonFiles(in []string) (jsonContent, error) {
	m := make(jsonContent)
	for _, e := range in {
		bytes, err := ioutil.ReadFile(e)
		if err != nil {
			return m, err
		}
		json, err := jsonql.NewStringQuery(string(bytes))
		if err != nil {
			return m, err
		}
		tag := strings.Split(e, ".")[0]
		m[tag] = json
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
func writeJsonResponse(w http.ResponseWriter, j jsonAble, err error) {
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

// Web compilers
func getIndexPage(issuesPhases *jsonql.JSONQL) (string, error) {
	issues, err := service.ListFirstIssues(issuesPhases)
	if err != nil {
		return "", err
	}
	introID := "intro"
	issuesLi := fmt.Sprintf(c["tab-li"], "active", introID, introID, introID, "Intro")
	issuesContent := fmt.Sprintf(c["tab-content"], "active", introID, introID, introID, c["tab-content-intro"])
	for _, e := range *issues {
		phaseID := fmt.Sprintf("phase%s", e.Phase.ID)
		li := fmt.Sprintf(c["tab-li"], "", phaseID, phaseID, phaseID, e.Phase.Name)
		issuesLi = fmt.Sprintf("%s%s", issuesLi, li)
		conPhase := ""
		for _, i := range e.List {
			year := ""
			if i.Date != "" {
				year = i.Date[:4]
			}
			conIssue := fmt.Sprintf(c["tab-content-phase"], e.Phase.ID, i.SortID, i.Pic, i.Title, e.Phase.ID, i.SortID, i.Title, year, "Protagonist")
			conPhase = fmt.Sprintf("%s%s", conPhase, conIssue)
		}
		con := fmt.Sprintf(c["tab-content"], "", phaseID, phaseID, phaseID, conPhase)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["tabs"], issuesLi, issuesContent)
	content = fmt.Sprintf(c["template"], content)
	return content, nil
}

func getIssuesPage(comics *jsonql.JSONQL, phaseid, sortid string) (string, error) {
	issues, err := service.ListComicsByPhaseAndSortIDs(comics, phaseid, sortid)
	if err != nil {
		return "", err
	}
	issuesContent := ""
	for _, e := range *issues {
		name := fmt.Sprintf("%s vol. %v #%v", e.Collection, e.Vol, e.Num)
		essential := "NO"
		if e.Essential {
			essential = "YES"
		}
		displayEvent := "block"
		if e.Event == "" {
			displayEvent = "none"
		}
		displayComments := "block"
		if e.Comments == "" {
			displayComments = "none"
		}
		con := fmt.Sprintf(c["issue-content"], name, e.PhaseID, e.SortID, e.Pic, name, 
			e.Collection, 
			e.Vol, 
			e.Num, 
			e.Date, 
			e.Universe, 
			e.PhaseName,
			displayEvent,
			e.Event,
			essential, 
			e.Characters, 
			e.Creators,
			displayComments,
			e.Comments,
		)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["issues"], phaseid, sortid, (*issues)[0].Title, issuesContent)
	content = fmt.Sprintf(c["template"], content)
	return content, nil
}
