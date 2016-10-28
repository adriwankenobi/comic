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

func ListComicsBySortID(comics *jsonql.JSONQL, sortid string) (*ComicList, error) {
	return FindComicList(comics, "sortid='"+sortid+"'")
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
func FindFirstIssuesByPhaseID(fissues *jsonql.JSONQL, id string) (*Fissues, error) {
	list, err := FindFissuesList(fissues, "phase.id='"+id+"'")
	if err != nil {
		return &Fissues{}, err
	}
	if len(*list) <= 0 {
		return &Fissues{}, nil
	}
	return &(*list)[0], nil
}

func ListFirstIssues(fissues *jsonql.JSONQL) (*FissuesList, error) {
	return FindFissuesList(fissues, "phase.id!=''")
}

// Find events
func FindEventByID(events *jsonql.JSONQL, id string) (*Event, error) {
	list, err := FindEventList(events, "id='"+id+"'")
	if err != nil {
		return &Event{}, err
	}
	if len(*list) <= 0 {
		return &Event{}, nil
	}
	return &(*list)[0], nil
}

func ListEvents(events *jsonql.JSONQL) (*EventList, error) {
	return FindEventList(events, "id!=''")
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

func FindFissuesList(fissues *jsonql.JSONQL, q string) (*FissuesList, error) {
	result, err := fissues.Query(q)
	if err != nil {
		return nil, err
	}
	return NewFissuesList(result)
}

func FindEventList(events *jsonql.JSONQL, q string) (*EventList, error) {
	result, err := events.Query(q)
	if err != nil {
		return nil, err
	}
	return NewEventList(result)
}
