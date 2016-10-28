package web

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"strings"
)

// Index
func getIndexPage(phases *service.PhaseList, events *service.EventList) (string, error) {
	return getTemplate(c["intro"], phases, events, 0), nil
}

// Issues
func getIssuesPage(phases *service.PhaseList, events *service.EventList, issues *service.ComicList) (string, error) {
	if issues.IsEmpty() {
		return getNotFoundPage(phases, events), nil
	}

	issuesContent := ""
	for _, e := range *issues {
		name := fmt.Sprintf("%s vol. %v #%v", e.Collection, e.Vol, e.Num)
		essential := "NO"
		if e.Essential {
			essential = "YES"
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
			e.Event,
			essential,
			strings.Join(e.Characters, ", "),
			strings.Join(e.Creators, ", "),
			displayComments,
			commentList,
		)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["content-issues"], (*issues)[0].PhaseID, (*issues)[0].SortID, (*issues)[0].Title, issuesContent)
	return getTemplate(content, phases, events, -1), nil
}

// Phase
func getPhasePage(phases *service.PhaseList, events *service.EventList, fissues *service.Fissues) (string, error) {
	if fissues.IsEmpty() {
		return getNotFoundPage(phases, events), nil
	}

	issues := (*fissues).List
	phaseID := fissues.Phase.ID
	issuesContent := ""
	for _, i := range issues {
		year := ""
		if i.Date != "" {
			year = i.Date[:4]
		}
		conIssue := fmt.Sprintf(c["content-fissue"], phaseID, i.SortID, i.Pic, i.Title, year, "Protagonist", phaseID, i.SortID, i.Title)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, conIssue)
	}
	content := fmt.Sprintf(c["content"], fissues.Phase.Name, issuesContent)
	return getTemplate(content, phases, events, 1), nil
}

// About
func getAboutPage(phases *service.PhaseList, events *service.EventList) string {
	return getTemplate(c["about"], phases, events, 3)
}

// Not found
func getNotFoundPage(phases *service.PhaseList, events *service.EventList) string {
	return getTemplate(c["not-found"], phases, events, -1)
}

// Utils
// Menu
func getPhasesMenuList(phases []service.Phase, n int) []string {
	result := make([]string, n)
	m := len(phases) / n
	r := len(phases) % n
	start := 0
	for i := 0; i < n; i++ {
		list := ""
		end := start + m
		if r != 0 {
			end++
			r--
		}
		for j := start; j < end; j++ {
			link := fmt.Sprintf("/phases/%s", phases[j].ID)
			title := fmt.Sprintf("%v - %s", j+1, phases[j].Name)
			li := fmt.Sprintf(c["list-link"], link, title)
			list = fmt.Sprintf("%s%s", list, li)
		}
		result[i] = list
		start = end
	}
	return result
}

func getEventsMenuList(events []service.Event, n int) []string {
	result := make([]string, n)
	m := len(events) / n
	r := len(events) % n
	start := 0
	for i := 0; i < n; i++ {
		list := ""
		end := start + m
		if r != 0 {
			end++
			r--
		}
		for j := start; j < end; j++ {
			link := fmt.Sprintf("/events/%s", events[j].ID)
			li := fmt.Sprintf(c["list-link"], link, events[j].Name)
			list = fmt.Sprintf("%s%s", list, li)
		}
		result[i] = list
		start = end
	}
	return result
}

func getTemplate(content string, phases *service.PhaseList, events *service.EventList, n int) string {
	phasesMenu := getPhasesMenuList(*phases, 3)
	eventsMenu := getEventsMenuList(*events, 3)
	active := [4]string{"", "", "", ""}
	if n >= 0 && n <= 3 {
		active[n] = "active"
	}
	return fmt.Sprintf(c["template"],
		active[0],
		active[1],
		phasesMenu[0],
		phasesMenu[1],
		phasesMenu[2],
		active[2],
		eventsMenu[0],
		eventsMenu[1],
		eventsMenu[2],
		active[3],
		content,
	)
}
