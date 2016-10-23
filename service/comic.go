package service

import (
	"encoding/json"
	"fmt"
	"github.com/adriwankenobi/comic/marvel"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

// TODO: Definitions class

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

var Datastore = DatastoreType{
	"comics":  nil,
	"phases":  nil,
	"fissues": nil,
}

// Comics
type Comic struct {
	ID         string   `json:"id,omitempty"`         // From Marvel API
	Collection string   `json:"collection,omitempty"` // From XLSX
	Title      string   `json:"title,omitempty"`      // From XLSX
	Vol        int      `json:"vol,omitempty"`        // From XLSX
	Num        int      `json:"num,omitempty"`        // From XLSX
	Date       string   `json:"date,omitempty"`       // From Marvel API
	Event      string   `json:"event,omitempty"`      // From XLSX
	Characters []string `json:"characters,omitempty"` // From Marvel API
	Creators   []string `json:"creators,omitempty"`   // From Marvel API
	Pic        string   `json:"pic,omitempty"`        // From Marvel API
	Universe   string   `json:"universe,omitempty"`   // From XLSX
	Essential  bool     `json:"essential,omitempty"`  // From XLSX
	Comments   string   `json:"comments,omitempty"`   // From XLSX
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
			c.Comments = e.(string)
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

// Constructor from XLSX
func NewComicListFromXLSX(path, out string) error {
	// New comic list
	comics := ComicList{}

	// New phase list
	phases := PhaseList{}

	// New issues phase list
	fissues := FissuesList{}

	// Open file
	xls, err := xlsx.OpenFile(path)
	if err != nil {
		return err
	}

	// Loop through file sheets
	for sheet_i, sheet := range xls.Sheets {
		p := Phase{}
		p.ID, err = getCode(sheet_i + 1)
		if err != nil {
			return err
		}
		p.Name = sheet.Name
		phases = append(phases, p)

		i := Fissues{}
		i.Phase = p
		i.List = ComicList{}

		lastTitle := ""
		sortID := 0
		for _, row := range sheet.Rows[1:] {
			id, err := row.Cells[id_col].String()
			if err != nil {
				return err
			}
			collection, err := row.Cells[collection_col].String()
			if err != nil {
				return err
			}
			vol, err := row.Cells[vol_col].Int()
			if err != nil {
				return err
			}
			num, err := row.Cells[num_col].Int()
			if err != nil {
				return err
			}
			title, err := row.Cells[title_col].String()
			if err != nil {
				return err
			}
			date, err := row.Cells[date_col].String()
			if err != nil {
				return err
			}
			event, err := row.Cells[event_col].String()
			if err != nil {
				return err
			}
			characters, err := row.Cells[characters_col].String()
			if err != nil {
				return err
			}
			creators, err := row.Cells[creators_col].String()
			if err != nil {
				return err
			}
			pic, err := row.Cells[pic_col].String()
			if err != nil {
				return err
			}
			universe, err := row.Cells[universe_col].String()
			if err != nil {
				return err
			}
			essential, err := row.Cells[essential_col].String()
			if err != nil {
				return err
			}
			var comments string
			if len(row.Cells) > mandatory_cols {
				comments, err = row.Cells[comments_col].String()
				if err != nil {
					return err
				}
			}
			c := Comic{}
			c.ID = id
			c.Collection = collection
			c.Vol = vol
			c.Num = num
			c.Title = title
			c.Date = date
			c.Event = event
			c.Characters = strings.Split(characters, ", ")
			c.Creators = strings.Split(creators, ", ")
			c.Pic = pic
			c.Universe = universe
			c.Essential = essential == "YES"
			c.Comments = comments
			c.PhaseID = p.ID
			c.PhaseName = p.Name
			if title != lastTitle {
				sortID++
				lastTitle = title
				sID, err := getCode(sortID)
				if err != nil {
					return err
				}
				i.List = append(i.List, Comic{
					Pic:    pic,
					Title:  title,
					Date:   date,
					SortID: sID,
				})
			}
			c.SortID, err = getCode(sortID)
			if err != nil {
				return err
			}
			comics = append(comics, c)
		}

		fissues = append(fissues, i)
	}
	Datastore["comics"] = &comics
	Datastore["phases"] = &phases
	Datastore["fissues"] = &fissues
	return nil
}

// Update XLSX from MARVEL API
func UpdateXLSX(path string, start, end int, mPubKey, mPriKey string) error {
	// Open file
	xls, err := xlsx.OpenFile(path)
	if err != nil {
		return err
	}

	// New Marvel API
	m := marvel.NewMarvelAPI(mPubKey, mPriKey)

	// Gracefull shutdown if user presses Ctrl+C
	stop := false
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		stop = true
	}()

	// Loop through file sheets
	for _, sheet := range xls.Sheets {
		if !stop {
			for _, row := range sheet.Rows[1:] {
				if !stop {
					id, err := row.Cells[id_col].String()
					if err != nil {
						return err
					}
					fill := row.Cells[id_col].GetStyle().Fill
					if id == "" && fill.PatternType != "solid" && fill.FgColor != "FFFF0000" {
						collection, err := row.Cells[collection_col].String()
						if err != nil {
							return err
						}
						num, err := row.Cells[num_col].Int()
						if err != nil {
							return err
						}
						fmt.Printf("[Finding] %s %v\n", collection, num)
						id, err = m.Find(collection, num, start, end)
						if err != nil || id == "" {
							fmt.Printf("%s\n", err.Error())
						} else {
							row.Cells[id_col].SetString(id)
						}
					}
					if id != "" {
						date, err := row.Cells[date_col].String()
						if err != nil {
							return err
						}
						characters, err := row.Cells[characters_col].String()
						if err != nil {
							return err
						}
						creators, err := row.Cells[creators_col].String()
						if err != nil {
							return err
						}
						pic, err := row.Cells[pic_col].String()
						if err != nil {
							return err
						}
						if date == "" && characters == "" && creators == "" && pic == "" {
							fmt.Printf("[Finding] %s\n", id)
							data, err := m.FindByID(id)
							if err != nil {
								fmt.Printf("%s\n", err.Error())
							} else {
								row.Cells[date_col].SetString(data.Date)
								row.Cells[characters_col].SetString(data.Characters)
								row.Cells[creators_col].SetString(data.Creators)
								row.Cells[pic_col].SetString(data.Pic)
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("Done!")

	// Save file
	fmt.Println("Saving file")
	xls.Save(path)
	return nil
}

// Create folders structure based on XLSX
func CreateFolders(f, path string) error {
	// Open file
	xls, err := xlsx.OpenFile(f)
	if err != nil {
		return err
	}

	// Read all files in path
	phaseFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// Loop through file sheets
	for sheet_i, sheet := range xls.Sheets {
		// Get starter code
		starter, err := getCode(sheet_i + 1)
		if err != nil {
			return err
		}
		// Find folder
		phaseFolderName := fmt.Sprintf("%v - %s", starter, sheet.Name)
		phaseFolderNameFull := fmt.Sprintf("%s/%s", path, phaseFolderName)
		found := false
		for _, file := range phaseFiles {
			if file.IsDir() && file.Name() == phaseFolderName {
				found = true
			}
		}
		// Create folder if it doesn't exist
		if !found {
			fmt.Printf("[Creating folder] %s\n", phaseFolderNameFull)
			os.Mkdir(phaseFolderNameFull, os.ModeDir)
		}

		// Read all files in this phase
		files, err := ioutil.ReadDir(phaseFolderNameFull)
		if err != nil {
			return err
		}

		// Loop throw rows to check second level folders
		sortID := 0
		lastTitle := ""
		for _, row := range sheet.Rows[1:] {
			filesModified := false
			id, err := row.Cells[id_col].String()
			if err != nil {
				return err
			}
			fill := row.Cells[id_col].GetStyle().Fill
			collection, err := row.Cells[collection_col].String()
			if err != nil {
				return err
			}
			title, err := row.Cells[title_col].String()
			if err != nil {
				return err
			}
			vol, err := row.Cells[vol_col].Int()
			if err != nil {
				return err
			}
			num, err := row.Cells[num_col].Int()
			if err != nil {
				return err
			}
			isNew := false
			if id == "" && fill.PatternType != "solid" && fill.FgColor != "FFFF0000" {
				isNew = true
				fmt.Printf("[New comic found] %s %v %v\n", collection, vol, num)
			}
			needNewFolder := false
			if title != lastTitle {
				sortID++
				lastTitle = title
				needNewFolder = isNew
			}
			// Get code
			code, err := getCode(sortID)
			if err != nil {
				return err
			}
			// Find folder
			folderName := fmt.Sprintf("%s/%s", phaseFolderNameFull, code)
			exists := ""
			for _, file := range files {
				if file.IsDir() && strings.HasPrefix(file.Name(), code) {
					exists = file.Name()
				}
			}
			if needNewFolder && exists != "" {
				// Rename folders from 'sortID' until last
				for i := len(files); i >= sortID; i-- {
					oldCode, err := getCode(i)
					if err != nil {
						return err
					}
					newCode, err := getCode(i + 1)
					if err != nil {
						return err
					}
					found := ""
					for _, file := range files {
						if file.IsDir() && strings.HasPrefix(file.Name(), oldCode) {
							found = file.Name()
						}
					}
					if found == "" {
						return fmt.Errorf("[Error] Cannot find folder to rename")
					}

					oldPath := fmt.Sprintf("%s/%s", phaseFolderNameFull, found)
					newPath := fmt.Sprintf("%s/%s", phaseFolderNameFull, newCode)
					fmt.Printf("[Renaming folders] From %v to %v\n", oldPath, newPath)
					os.Rename(oldPath, newPath)
					filesModified = true
				}
			}
			// Create folder if it doesn't exist
			if (needNewFolder && exists != "") || exists == "" {
				fmt.Printf("[Creating folder] %s (%s %v %v)\n", folderName, collection, vol, num)
				os.Mkdir(folderName, os.ModeDir)
				filesModified = true
			}
			// Refresh files
			if filesModified {
				files, err = ioutil.ReadDir(phaseFolderNameFull)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Util
func getCode(i int) (string, error) {
	if i > 999 {
		return "", fmt.Errorf("[Error] Cannot get code higer than 999")
	}
	if i < 10 {
		return fmt.Sprintf("00%v", i), nil
	}
	if i < 100 {
		return fmt.Sprintf("0%v", i), nil
	}
	return fmt.Sprintf("%v", i), nil
}
