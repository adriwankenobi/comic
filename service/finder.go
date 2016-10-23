package service

import (
	"github.com/elgs/jsonql"
)

// Find comics
func FindComicByID(comics *jsonql.JSONQL, id string) (*Comic, error) {
	list, err := FindComicList(comics, "id='"+id+"'")
	if err != nil {
		return &Comic{}, err
	}
	if len(*list) <= 0 {
		return &Comic{}, nil
	}
	return &(*list)[0], nil
}

func ListComics(comics *jsonql.JSONQL) (*ComicList, error) {
	// HEAVY
	return FindComicList(comics, "id!=''")
}

func ListComicsByPhaseID(comics *jsonql.JSONQL, id string) (*ComicList, error) {
	return FindComicList(comics, "phaseid='"+id+"'")
}

func ListComicsByPhaseAndSortIDs(comics *jsonql.JSONQL, id, sortid string) (*ComicList, error) {
	return FindComicList(comics, "phaseid='"+id+"' && sortid='"+sortid+"'")
}

// Find phases
func FindPhaseByID(phases *jsonql.JSONQL, id string) (*Phase, error) {
	list, err := FindPhaseList(phases, "id='"+id+"'")
	if err != nil {
		return &Phase{}, err
	}
	if len(*list) <= 0 {
		return &Phase{}, nil
	}
	return &(*list)[0], nil
}

func ListPhases(phases *jsonql.JSONQL) (*PhaseList, error) {
	return FindPhaseList(phases, "id!=''")
}

// Find first issues
func FindFirstIssuesByPhaseID(issuesPhases *jsonql.JSONQL, id string) (*IssuesPhase, error) {
	list, err := FindIssuesPhaseList(issuesPhases, "phase.id='"+id+"'")
	if err != nil {
		return &IssuesPhase{}, err
	}
	if len(*list) <= 0 {
		return &IssuesPhase{}, nil
	}
	return &(*list)[0], nil
}

func ListFirstIssues(issuesPhases *jsonql.JSONQL) (*IssuesPhaseList, error) {
	return FindIssuesPhaseList(issuesPhases, "phase.id!=''")
}

// Utils
func FindComicList(comics *jsonql.JSONQL, q string) (*ComicList, error) {
	result, err := comics.Query(q)
	if err != nil {
		return nil, err
	}
	return NewComicList(result)
}

func FindPhaseList(phases *jsonql.JSONQL, q string) (*PhaseList, error) {
	result, err := phases.Query(q)
	if err != nil {
		return nil, err
	}
	return NewPhaseList(result)
}

func FindIssuesPhaseList(issuesPhases *jsonql.JSONQL, q string) (*IssuesPhaseList, error) {
	result, err := issuesPhases.Query(q)
	if err != nil {
		return nil, err
	}
	return NewIssuesPhaseList(result)
}
