package service

import (
	"github.com/elgs/jsonql"
	"sort"
)

// Menu
type Menu struct {
	URI          string
	IsEssentials bool
	Phases       *NamableList
	Events       *NamableList
	Characters   *NamableList
}

// Get menu
func GetMenu(phases *jsonql.JSONQL, events *jsonql.JSONQL, characters *jsonql.JSONQL) (Menu, error) {
	m := Menu{}
	phaseList, err := ListNamables(phases)
	if err != nil {
		return m, err
	}
	m.Phases = phaseList
	eventList, err := ListNamables(events)
	if err != nil {
		return m, err
	}
	m.Events = eventList
	charsList, err := ListNamables(characters)
	if err != nil {
		return m, err
	}
	sort.Sort(ByName(*charsList))
	m.Characters = charsList
	return m, nil
}

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

// Find namables
func FindNamableByID(namables *jsonql.JSONQL, id string) (*Namable, error) {
	list, err := FindNamableList(namables, "id='"+id+"'")
	if err != nil {
		return &Namable{}, err
	}
	if len(*list) <= 0 {
		return &Namable{}, nil
	}
	return &(*list)[0], nil
}

func ListNamables(namables *jsonql.JSONQL) (*NamableList, error) {
	return FindNamableList(namables, "id!=''")
}

// Find first issues
func FindFirstIssuesByID(fissues *jsonql.JSONQL, id string) (*Fissues, error) {
	list, err := FindFissuesList(fissues, "namable.id='"+id+"'")
	if err != nil {
		return &Fissues{}, err
	}
	if len(*list) <= 0 {
		return &Fissues{}, nil
	}
	return &(*list)[0], nil
}

func ListFirstIssues(fissues *jsonql.JSONQL) (*FissuesList, error) {
	return FindFissuesList(fissues, "namable.id!=''")
}

// Utils
func FindComicList(comics *jsonql.JSONQL, q string) (*ComicList, error) {
	result, err := comics.Query(q)
	if err != nil {
		return nil, err
	}
	return NewComicList(result)
}

func FindNamableList(namables *jsonql.JSONQL, q string) (*NamableList, error) {
	result, err := namables.Query(q)
	if err != nil {
		return nil, err
	}
	return NewNamableList(result)
}

func FindFissuesList(fissues *jsonql.JSONQL, q string) (*FissuesList, error) {
	result, err := fissues.Query(q)
	if err != nil {
		return nil, err
	}
	return NewFissuesList(result)
}
