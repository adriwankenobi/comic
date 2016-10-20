package service

import (
    "github.com/elgs/jsonql"
)

func FindComic(comics *jsonql.JSONQL, id string) (*Comic, error) {
	list, err := FindComicList(comics, "id='" + id + "'")
	if err != nil {
		return &Comic{}, err
	}
	if len(*list) <= 0 {
		return &Comic{}, nil
	}
	return &(*list)[0], nil
}

func ListComics(comics *jsonql.JSONQL) (*ComicList, error) {
	return FindComicList(comics, "title='Inhumans'") // TODO
}

func FindComicList(comics *jsonql.JSONQL, q string) (*ComicList, error) {
    result, err := comics.Query(q)
    if err != nil {
        return nil, err
    }
    return NewComicList(result)
}

func FindPhase(phases *jsonql.JSONQL, id string) (*Phase, error) {
	list, err := FindPhaseList(phases, "id='" + id + "'")
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

func FindPhaseList(phases *jsonql.JSONQL, q string) (*PhaseList, error) {
    result, err := phases.Query(q)
    if err != nil {
        return nil, err
    }
    return NewPhaseList(result)
}