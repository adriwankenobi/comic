package web

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"sort"
	"strings"
)

// Index
func getIndexPage(menu service.Menu) (string, error) {
	return getTemplate(c["intro"], menu, 0), nil
}

// Issues
func getIssuesPage(menu service.Menu, issues *service.ComicList) (string, error) {
	if issues.IsEmpty() {
		return getNotFoundPage(menu), nil
	}

	issuesContent := ""
	for _, e := range *issues {
		name := fmt.Sprintf("%s vol. %v #%v", e.Collection, e.Vol, e.Num)
		essential := "NO"
		if e.Essential {
			essential = "YES"
		}
		characters := []string{}
		for _, ch := range e.Characters {
			char := fmt.Sprintf(c["a-link"], "characters", ch.ID, ch.Name)
			characters = append(characters, char)
		}
		creators := []string{}
		for _, cr := range e.Creators {
			creat := fmt.Sprintf(c["a-link"], "creators", cr.ID, cr.Name)
			creators = append(creators, creat)
		}
		displayEvent := "block"
		if e.Event == "" {
			displayEvent = "none"
		}
		displayComments := "block"
		if len(e.Comments) <= 0 {
			displayComments = "none"
		}
		commentList := ""
		for _, e := range e.Comments {
			comment := fmt.Sprintf(c["list"], strings.Trim(e, " "))
			commentList = fmt.Sprintf("%s%s", commentList, comment)
		}

		con := fmt.Sprintf(c["content-issue"], name, e.PhaseID, e.SortID, e.Pic, name,
			e.Collection,
			e.Vol,
			e.Num,
			e.Date,
			e.Universe,
			e.PhaseID,
			e.PhaseName,
			displayEvent,
			e.EventID,
			e.Event,
			essential,
			strings.Join(characters, ", "),
			strings.Join(creators, ", "),
			displayComments,
			commentList,
		)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["content-issues"], (*issues)[0].PhaseID, (*issues)[0].SortID, (*issues)[0].Title, issuesContent)
	return getTemplate(content, menu, -1), nil
}

// Fissues
func getCharactersFissuesPage(menu service.Menu, fissues *service.Fissues) (string, error) {
	return getFissuesPage(menu, fissues, 1)
}

func getPhasesFissuesPage(menu service.Menu, fissues *service.Fissues) (string, error) {
	return getFissuesPage(menu, fissues, 2)
}

func getEventsFissuesPage(menu service.Menu, fissues *service.Fissues) (string, error) {
	return getFissuesPage(menu, fissues, 3)
}

func getCreatorsFissuesPage(menu service.Menu, fissues *service.Fissues) (string, error) {
	return getFissuesPage(menu, fissues, 4)
}

func getFissuesPage(menu service.Menu, fissues *service.Fissues, activeTab int) (string, error) {
	if fissues.IsEmpty() {
		return getNotFoundPage(menu), nil
	}

	issues := (*fissues).List
	phaseID := fissues.Namable.ID
	issuesContent := ""
	for _, i := range issues {
		comicList := ""
		l := []string{}
		m := map[string][]int{}
		for _, e := range i.ComicList {
			name := fmt.Sprintf("%s vol. %v", e.Collection, e.Vol)
			_, exists := m[name]
			if !exists {
				l = append(l, name)
				m[name] = []int{}
			}
			found := false
			for _, n := range m[name] {
				if e.Num == n {
					found = true
					break
				}
			}
			if !found {
				m[name] = append(m[name], e.Num)
			}
		}
		for _, k := range l {
			v := m[k]
			sort.Ints(v)
			name := fmt.Sprintf("%s #%v", k, v[0])
			if len(v) > 1 {
				for i := 1; i < len(v); i++ {
					if v[i] != v[i-1]+1 {
						name = fmt.Sprintf("%s - #%v", name, v[i-1])
						h6 := fmt.Sprintf(c["h6"], name)
						comicList = fmt.Sprintf("%s%s", comicList, h6)
						name = fmt.Sprintf("%s #%v", k, v[i])
					}
				}
				name = fmt.Sprintf("%s - #%v", name, v[len(v)-1])
			}
			h6 := fmt.Sprintf(c["h6"], name)
			comicList = fmt.Sprintf("%s%s", comicList, h6)
		}

		conIssue := fmt.Sprintf(c["content-fissue"], i.PhaseID, i.SortID, i.Pic, i.Title, i.Date[:4],
			i.Characters[0].ID, i.Characters[0].Name, phaseID, i.SortID, i.Title, comicList)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, conIssue)
	}
	content := fmt.Sprintf(c["content"], fissues.Namable.Name, issuesContent)
	return getTemplate(content, menu, activeTab), nil
}

// Creators
func getCreatorsPage(menu service.Menu, creators *service.NamableList) string {
	sort.Sort(service.ByName(*creators))
	columns := ""
	n := 8
	for i := 0; i < n; i++ {
		k := len(*creators) / n
		first := i * k
		last := first + k - 1
		list := ""
		for j := first; j <= last; j++ {
			e := (*creators)[j]
			a := fmt.Sprintf(c["a-link"], "creators", e.ID, e.Name)
			li := fmt.Sprintf(c["list"], a)
			list = fmt.Sprintf("%s%s", list, li)
		}
		ul := fmt.Sprintf(c["ul"], list)
		column := fmt.Sprintf(c["div-left"], ul)
		columns = fmt.Sprintf("%s%s", columns, column)
	}
	content := fmt.Sprintf(c["content"], "Creators", columns)
	return getTemplate(content, menu, 4)
}

// About
func getAboutPage(menu service.Menu) string {
	return getTemplate(c["about"], menu, 5)
}

// Not found
func getNotFoundPage(menu service.Menu) string {
	return getTemplate(c["not-found"], menu, -1)
}

// Utils
// Menu
func getMenuList(namables service.NamableList, n int, link string, showID bool) []string {
	result := make([]string, n)
	m := len(namables) / n
	r := len(namables) % n
	start := 0
	for i := 0; i < n; i++ {
		list := ""
		end := start + m
		if r != 0 {
			end++
			r--
		}
		for j := start; j < end; j++ {
			title := namables[j].Name
			if showID {
				title = fmt.Sprintf("%v - %s", j+1, title)
			}
			alink := fmt.Sprintf(c["a-link"], link, namables[j].ID, title)
			li := fmt.Sprintf(c["list"], alink)
			list = fmt.Sprintf("%s%s", list, li)
		}
		result[i] = list
		start = end
	}
	return result
}

func getTemplate(content string, menu service.Menu, activeTab int) string {
	phasesMenu := getMenuList(*menu.Phases, 3, "phases", true)
	eventsMenu := getMenuList(*menu.Events, 3, "events", false)
	charactersMenu := getMenuList(*menu.Characters, 8, "characters", false)
	active := [6]string{"", "", "", "", "", ""}
	if activeTab >= 0 && activeTab <= 5 {
		active[activeTab] = "active"
	}
	return fmt.Sprintf(c["template"],
		active[0],
		active[1],
		charactersMenu[0],
		charactersMenu[1],
		charactersMenu[2],
		charactersMenu[3],
		charactersMenu[4],
		charactersMenu[5],
		charactersMenu[6],
		charactersMenu[7],
		active[2],
		phasesMenu[0],
		phasesMenu[1],
		phasesMenu[2],
		active[3],
		eventsMenu[0],
		eventsMenu[1],
		eventsMenu[2],
		active[4],
		active[5],
		content,
	)
}
