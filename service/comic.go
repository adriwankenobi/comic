package service

import (
	"fmt"
	"github.com/adriwankenobi/comic/marvel"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

// JSON generator
func JsonGenerator(path, out string) error {
	// New comic list
	comics := ComicList{}

	// New phase list
	phases := NamableList{}

	// New events list
	events := NamableList{}

	// New issues phase list
	fissues := FissuesList{}

	// Open file
	xls, err := xlsx.OpenFile(path)
	if err != nil {
		return err
	}

	eventsMap := map[string]Namable{}
	eventsComics := map[string]*ComicList{}
	eventID := 0

	// Loop through file sheets
	for sheet_i, sheet := range xls.Sheets {
		p := Namable{}
		p.ID, err = getCode(sheet_i + 1)
		if err != nil {
			return err
		}
		p.Name = sheet.Name
		phases = append(phases, p)

		i := Fissues{}
		i.Phase = p
		i.List = ComicList{}

		cp := ComicList{}

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
			if event != "" {
				c.Event = event
				e, exists := eventsMap[event]
				if !exists {
					eventID++
					eID, err := getCode(eventID)
					if err != nil {
						return err
					}
					e = Namable{ID: eID, Name: event}
					eventsMap[event] = e
					eventsComics[eID] = &ComicList{}
					events = append(events, e)
				}
				c.EventID = e.ID
			}
			c.Characters = strings.Split(characters, ", ")
			c.Creators = strings.Split(creators, ", ")
			c.Pic = pic
			c.Universe = universe
			c.Essential = essential == "YES"
			if comments != "" {
				c.Comments = strings.Split(comments, ", ")
			}
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
			cp = append(cp, c)
			if event != "" {
				*(eventsComics[c.EventID]) = append(*(eventsComics[c.EventID]), c)
			}
		}

		fissues = append(fissues, i)
		Datastore[fmt.Sprintf("comics-phase-%s", p.ID)] = &cp
	}
	Datastore["comics"] = &comics
	Datastore["phases"] = &phases
	Datastore["fissues"] = &fissues
	Datastore["events"] = &events
	for key, value := range eventsComics {
		Datastore[fmt.Sprintf("comics-event-%s", key)] = value
	}
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
