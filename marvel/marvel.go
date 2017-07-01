package marvel

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL              = "http://gateway.marvel.com/v1/public"
	marvelDateFormat     = "2006-01-02T15:04:05-0700"
	marvelResponseFormat = "2006-01-02"
)

type MarvelAPI struct {
	publicKey  string
	privateKey string
}

type MarvelResponse struct {
	Date       string
	Pic        string
	Creators   string
	Characters string
}

type dateResponse struct {
	Date     string `json:"date"`
	TypeDate string `json:"type"`
}

type datesResponse []dateResponse

func (d *datesResponse) find(criteria string) string {
	for _, e := range *d {
		if e.TypeDate == criteria {
			return e.Date
		}
	}
	return ""
}

type thumbnailResponse struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

type items struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type itemsResponse struct {
	Available int     `json:"available"`
	Returned  int     `json:"returned"`
	Items     []items `json:"items"`
}

func (i *itemsResponse) toString() string {
	data := []string{}
	for _, e := range i.Items {
		data = append(data, e.Name)
	}
	return strings.Join(data, ", ")
}

type result struct {
	ID         int               `json:"id"`
	Dates      datesResponse     `json:"dates"`
	Thumbnail  thumbnailResponse `json:"thumbnail"`
	Creators   itemsResponse     `json:"creators"`
	Characters itemsResponse     `json:"characters"`
}

type dataResponse struct {
	Total   int      `json:"total"`
	Results []result `json:"results"`
}

type response struct {
	Code int          `json:"code"`
	Data dataResponse `json:"data"`
}

func NewMarvelAPI(pubKey, priKey string) MarvelAPI {
	return MarvelAPI{
		publicKey:  pubKey,
		privateKey: priKey,
	}
}

func (m *MarvelAPI) Find(collection string, num float64, start, end int) (string, error) {
	parameters := fmt.Sprintf("%s&title=%s&issueNumber=%v&dateRange=%v-01-01,%v-12-31",
		m.getDefaultParameters(), url.QueryEscape(collection), num, start, end)
	marvelURL := fmt.Sprintf("%s/comics?%s", baseURL, parameters)
	resp, err := do(marvelURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resp.Data.Results[0].ID), nil
}

func (m *MarvelAPI) FindByID(id string) (MarvelResponse, error) {
	marvelResp := MarvelResponse{}
	parameters := m.getDefaultParameters()
	marvelURL := fmt.Sprintf("%s/comics/%s?%s", baseURL, id, parameters)
	resp, err := do(marvelURL)
	if err != nil {
		return marvelResp, err
	}
	date, err := time.Parse(marvelDateFormat, resp.Data.Results[0].Dates.find("onsaleDate"))
	if err != nil {
		return marvelResp, err
	}
	marvelResp.Date = date.Format(marvelResponseFormat)
	marvelResp.Pic = fmt.Sprintf("%s.%s", resp.Data.Results[0].Thumbnail.Path, resp.Data.Results[0].Thumbnail.Extension)
	marvelResp.Creators = resp.Data.Results[0].Creators.toString()
	marvelResp.Characters = resp.Data.Results[0].Characters.toString()
	return marvelResp, nil
}

func (m *MarvelAPI) getDefaultParameters() string {
	now := time.Now().UTC()
	ts := now.Format(time.RFC3339)
	data := []byte(fmt.Sprintf("%s%s%s", ts, m.privateKey, m.publicKey))
	hash := fmt.Sprintf("%x", md5.Sum(data))
	return fmt.Sprintf("ts=%v&apikey=%s&hash=%s", ts, m.publicKey, hash)
}

func do(url string) (response, error) {
	decodedResp := response{}
	resp, err := http.Get(url)
	if err != nil {
		return decodedResp, err
	}
	if resp.StatusCode != http.StatusOK {
		return decodedResp, fmt.Errorf("[Fail] HTTP %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return decodedResp, err
	}
	err = decode(body, &decodedResp)
	if err != nil {
		return decodedResp, err
	}
	return decodedResp, nil
}

func decode(b []byte, resp *response) error {
	err := json.Unmarshal(b, resp)
	if err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return fmt.Errorf("[Fail] HTTP %v", resp.Code)
	}
	if resp.Data.Total != 1 {
		return fmt.Errorf("[Fail] Total comics found: %v", resp.Data.Total)
	}
	if resp.Data.Results[0].Creators.Available != resp.Data.Results[0].Creators.Returned {
		return fmt.Errorf("[Fail] Total creators is %v but found %v", resp.Data.Results[0].Creators.Available, resp.Data.Results[0].Creators.Returned)
	}
	if resp.Data.Results[0].Characters.Available != resp.Data.Results[0].Characters.Returned {
		return fmt.Errorf("[Fail] Total characters is %v but found %v", resp.Data.Results[0].Characters.Available, resp.Data.Results[0].Characters.Returned)
	}
	return nil
}
