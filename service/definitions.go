package service

import (
	"encoding/json"
	"fmt"
)

// XLSX columns
const (
	id_col         = 0
	collection_col = 1
	vol_col        = 2
	num_col        = 3
	title_col      = 4
	date_col       = 5
	event_col      = 6
	characters_col = 7
	creators_col   = 8
	pic_col        = 9
	universe_col   = 10
	essential_col  = 11
	comments_col   = 12

	mandatory_cols = 12
)

type JsonAble interface {
	ToJson() ([]byte, error)
	IsEmpty() bool
}

// Datastore
type DatastoreType map[string]JsonAble

var Datastore = DatastoreType{}

// Comics
type Comic struct {
	ID         string   `json:"id,omitempty"`         // From Marvel API
	Collection string   `json:"collection,omitempty"` // From XLSX
	Title      string   `json:"title,omitempty"`      // From XLSX
	Vol        int      `json:"vol,omitempty"`        // From XLSX
	Num        int      `json:"num,omitempty"`        // From XLSX
	Date       string   `json:"date,omitempty"`       // From Marvel API
	Event      string   `json:"event,omitempty"`      // From XLSX
	EventID    string   `json:"eventid,omitempty"`    // From XLSX
	Characters []string `json:"characters,omitempty"` // From Marvel API
	Creators   []string `json:"creators,omitempty"`   // From Marvel API
	Pic        string   `json:"pic,omitempty"`        // From Marvel API
	Universe   string   `json:"universe,omitempty"`   // From XLSX
	Essential  bool     `json:"essential,omitempty"`  // From XLSX
	Comments   []string `json:"comments,omitempty"`   // From XLSX
	PhaseID    string   `json:"phaseid,omitempty"`    // From XLSX: Generated based on sheet position
	PhaseName  string   `json:"phasename,omitempty"`  // From XLSX: Generated based on sheet name
	SortID     string   `json:"sortid,omitempty"`     // From XLSX: Generated based on row position
}
type ComicList []Comic

func (c *Comic) ToJson() ([]byte, error) {
	return json.MarshalIndent(c, "", "	")
}

func (c *Comic) IsEmpty() bool {
	return c.ID == "" && c.Collection == ""
}

func (c *ComicList) ToJson() ([]byte, error) {
	return json.MarshalIndent(c, "", "	")
}

func (c *ComicList) IsEmpty() bool {
	return len(*c) <= 0
}

// Phases
type Phase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type PhaseList []Phase

func (p *Phase) ToJson() ([]byte, error) {
	return json.MarshalIndent(p, "", "	")
}

func (p *Phase) IsEmpty() bool {
	return p.ID == "" && p.Name == ""
}

func (p *PhaseList) ToJson() ([]byte, error) {
	return json.MarshalIndent(p, "", "	")
}

func (p *PhaseList) IsEmpty() bool {
	return len(*p) <= 0
}

// Events
type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type EventList []Event

func (e *Event) ToJson() ([]byte, error) {
	return json.MarshalIndent(e, "", "	")
}

func (e *Event) IsEmpty() bool {
	return e.ID == "" && e.Name == ""
}

func (e *EventList) ToJson() ([]byte, error) {
	return json.MarshalIndent(e, "", "	")
}

func (e *EventList) IsEmpty() bool {
	return len(*e) <= 0
}

// First issues
type Fissues struct {
	Phase Phase     `json:"phase"`
	List  ComicList `json:"list,omitempty"`
}
type FissuesList []Fissues

func (f *Fissues) ToJson() ([]byte, error) {
	return json.MarshalIndent(f, "", "	")
}

func (f *Fissues) IsEmpty() bool {
	return f.Phase.IsEmpty() && f.List.IsEmpty()
}

func (f *FissuesList) ToJson() ([]byte, error) {
	return json.MarshalIndent(f, "", "	")
}

func (f *FissuesList) IsEmpty() bool {
	return len(*f) <= 0
}

// Constructors from jsonql
func NewComic(in interface{}) (Comic, error) {
	m := in.(map[string]interface{})
	c := Comic{}
	for i, e := range m {
		switch i {
		case "id":
			c.ID = e.(string)
			break
		case "collection":
			c.Collection = e.(string)
			break
		case "title":
			c.Title = e.(string)
			break
		case "vol":
			c.Vol = int(e.(float64))
			break
		case "num":
			c.Num = int(e.(float64))
			break
		case "date":
			c.Date = e.(string)
			break
		case "event":
			c.Event = e.(string)
			break
		case "eventid":
			c.EventID = e.(string)
			break
		case "characters":
			c.Characters = NewStringList(e)
			break
		case "creators":
			c.Creators = NewStringList(e)
			break
		case "pic":
			c.Pic = e.(string)
			break
		case "universe":
			c.Universe = e.(string)
			break
		case "essential":
			c.Essential = e.(bool)
			break
		case "comments":
			c.Comments = NewStringList(e)
			break
		case "phaseid":
			c.PhaseID = e.(string)
			break
		case "phasename":
			c.PhaseName = e.(string)
			break
		case "sortid":
			c.SortID = e.(string)
			break
		default:
			return c, fmt.Errorf("Unknown field: %v", i)
		}
	}
	return c, nil
}

func NewComicList(in interface{}) (*ComicList, error) {
	all := in.([]interface{})
	comics := make(ComicList, len(all))
	for i, e := range all {
		m := e.(map[string]interface{})
		c, err := NewComic(m)
		if err != nil {
			return &comics, err
		}
		comics[i] = c
	}
	return &comics, nil
}

func NewStringList(in interface{}) []string {
	all := in.([]interface{})
	ss := make([]string, len(all))
	for i, e := range all {
		c := e.(string)
		ss[i] = c
	}
	return ss
}

func NewPhase(in interface{}) (Phase, error) {
	m := in.(map[string]interface{})
	p := Phase{}
	for i, e := range m {
		switch i {
		case "id":
			p.ID = e.(string)
			break
		case "name":
			p.Name = e.(string)
			break
		default:
			return p, fmt.Errorf("Unknown field: %v", i)
		}
	}
	if p.ID == "" {
		return p, fmt.Errorf("Phase doesn't contain 'id' field: %v", p)
	}
	return p, nil
}

func NewPhaseList(in interface{}) (*PhaseList, error) {
	all := in.([]interface{})
	phases := make(PhaseList, len(all))
	for i, e := range all {
		m := e.(map[string]interface{})
		p, err := NewPhase(m)
		if err != nil {
			return &phases, err
		}
		phases[i] = p
	}
	return &phases, nil
}

func NewEvent(in interface{}) (Event, error) {
	m := in.(map[string]interface{})
	ev := Event{}
	for i, e := range m {
		switch i {
		case "id":
			ev.ID = e.(string)
			break
		case "name":
			ev.Name = e.(string)
			break
		default:
			return ev, fmt.Errorf("Unknown field: %v", i)
		}
	}
	if ev.ID == "" {
		return ev, fmt.Errorf("Event doesn't contain 'id' field: %v", ev)
	}
	return ev, nil
}

func NewEventList(in interface{}) (*EventList, error) {
	all := in.([]interface{})
	events := make(EventList, len(all))
	for i, e := range all {
		m := e.(map[string]interface{})
		ev, err := NewEvent(m)
		if err != nil {
			return &events, err
		}
		events[i] = ev
	}
	return &events, nil
}

func NewFissues(in interface{}) (Fissues, error) {
	m := in.(map[string]interface{})
	is := Fissues{}
	for i, e := range m {
		switch i {
		case "phase":
			phase, err := NewPhase(e)
			if err != nil {
				return is, err
			}
			is.Phase = phase
			break
		case "list":
			list, err := NewComicList(e)
			if err != nil {
				return is, err
			}
			is.List = *list
			break
		default:
			return is, fmt.Errorf("Unknown field: %v", i)
		}
	}
	return is, nil
}

func NewFissuesList(in interface{}) (*FissuesList, error) {
	all := in.([]interface{})
	fissues := make(FissuesList, len(all))
	for i, e := range all {
		m := e.(map[string]interface{})
		is, err := NewFissues(m)
		if err != nil {
			return &fissues, err
		}
		fissues[i] = is
	}
	return &fissues, nil
}
