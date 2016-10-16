package service

import (
	"io/ioutil"
    "github.com/elgs/jsonql"
)

func FindComicList(path, q string) (ComicList, error) {
	file, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    data := string(file)
    parser, err := jsonql.NewStringQuery(data)
    if err != nil {
        return nil, err
    }
    result, err := parser.Query(q)
    if err != nil {
        return nil, err
    }
    return NewComicList(result)
}