package web

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"strings"
)

// Index
func getIndexPage(issues *service.FissuesList) (string, error) {
	introID := "intro"
	issuesLi := fmt.Sprintf(c["tab-li"], "active", introID, introID, introID, "Intro")
	issuesContent := fmt.Sprintf(c["tab-content"], "active", introID, introID, introID, c["tab-content-intro"])
	phases := service.PhaseList{}
	for _, e := range *issues {
		phaseID := fmt.Sprintf("phase%s", e.Phase.ID)
		li := fmt.Sprintf(c["tab-li"], "", phaseID, phaseID, phaseID, e.Phase.Name)
		issuesLi = fmt.Sprintf("%s%s", issuesLi, li)
		phases = append(phases, e.Phase)
		conPhase := ""
		for _, i := range e.List {
			year := ""
			if i.Date != "" {
				year = i.Date[:4]
			}
			conIssue := fmt.Sprintf(c["tab-content-phase"], e.Phase.ID, i.SortID, i.Pic, i.Title, year, "Protagonist", e.Phase.ID, i.SortID, i.Title)
			conPhase = fmt.Sprintf("%s%s", conPhase, conIssue)
		}
		con := fmt.Sprintf(c["tab-content"], "", phaseID, phaseID, phaseID, conPhase)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["tabs"], issuesLi, issuesContent)
	return getTemplate(content, &phases), nil
}

// Issues
func getIssuesPage(phases *service.PhaseList, issues *service.ComicList) (string, error) {
	if issues.IsEmpty() {
		return getNotFoundPage(phases), nil
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
		if e.Comments == "" {
			displayComments = "none"
		}
		// TODO: For each comment
		commentsList := fmt.Sprintf(c["list"], e.Comments)

		con := fmt.Sprintf(c["issue-content"], name, e.PhaseID, e.SortID, e.Pic, name,
			e.Collection,
			e.Vol,
			e.Num,
			e.Date,
			e.Universe,
			e.PhaseName,
			displayEvent,
			e.Event,
			essential,
			strings.Join(e.Characters, ", "),
			strings.Join(e.Creators, ", "),
			displayComments,
			commentsList,
		)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["issues"], (*issues)[0].PhaseID, (*issues)[0].SortID, (*issues)[0].Title, issuesContent)
	return getTemplate(content, phases), nil
}

// Phase
func getPhasePage(phases *service.PhaseList, fissues *service.Fissues) (string, error) {
	if fissues.IsEmpty() {
		return getNotFoundPage(phases), nil
	}

	issues := (*fissues).List
	phaseID := fissues.Phase.ID
	issuesContent := ""
	for _, i := range issues {
		year := ""
		if i.Date != "" {
			year = i.Date[:4]
		}
		conIssue := fmt.Sprintf(c["tab-content-phase"], phaseID, i.SortID, i.Pic, i.Title, year, "Protagonist", phaseID, i.SortID, i.Title)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, conIssue)
	}
	content := fmt.Sprintf(c["phase"], fissues.Phase.Name, issuesContent)
	return getTemplate(content, phases), nil
}

// About
func getAboutPage(phases *service.PhaseList) string {
	return getTemplate(c["about"], phases)
}

// Not found
func getNotFoundPage(phases *service.PhaseList) string {
	return getTemplate(c["not-found"], phases)
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

func getTemplate(content string, phases *service.PhaseList) string {
	phasesMenu := getPhasesMenuList(*phases, 3)
	return fmt.Sprintf(c["template"],
		phasesMenu[0],
		phasesMenu[1],
		phasesMenu[2],
		content,
	)
}
