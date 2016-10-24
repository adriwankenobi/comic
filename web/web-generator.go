package web

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"strings"
)

func getIndexPage(issues *service.FissuesList) (string, error) {
	introID := "intro"
	issuesLi := fmt.Sprintf(c["tab-li"], "active", introID, introID, introID, "Intro")
	issuesContent := fmt.Sprintf(c["tab-content"], "active", introID, introID, introID, c["tab-content-intro"])
	phases := []service.Phase{}
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
			pic := "http://i.annihil.us/u/prod/marvel/i/mg/b/40/image_not_available.jpg"
			if i.Pic != "" {
				pic = i.Pic
			}
			conIssue := fmt.Sprintf(c["tab-content-phase"], e.Phase.ID, i.SortID, pic, i.Title, year, "Protagonist", e.Phase.ID, i.SortID, i.Title)
			conPhase = fmt.Sprintf("%s%s", conPhase, conIssue)
		}
		con := fmt.Sprintf(c["tab-content"], "", phaseID, phaseID, phaseID, conPhase)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["tabs"], issuesLi, issuesContent)
	phasesMenu := getPhasesMenuList(phases, 3)
	content = fmt.Sprintf(c["template"], phasesMenu[0], phasesMenu[1], phasesMenu[2], content)
	return content, nil
}

func getIssuesPage(phases *service.PhaseList, issues *service.ComicList) (string, error) {
	phasesMenu := getPhasesMenuList(*phases, 3)
	if issues.IsEmpty() {
		return getNotFoundPage(phasesMenu), nil
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
		pic := "http://i.annihil.us/u/prod/marvel/i/mg/b/40/image_not_available.jpg"
		if e.Pic != "" {
			pic = e.Pic
		}
		// TODO: For each comment
		commentsList := fmt.Sprintf(c["list"], e.Comments)

		con := fmt.Sprintf(c["issue-content"], name, e.PhaseID, e.SortID, pic, name,
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
	content = fmt.Sprintf(c["template"], phasesMenu[0], phasesMenu[1], phasesMenu[2], content)
	return content, nil
}

func getPhasePage(phases *service.PhaseList, fissues *service.Fissues) (string, error) {
	phasesMenu := getPhasesMenuList(*phases, 3)
	if fissues.IsEmpty() {
		return getNotFoundPage(phasesMenu), nil
	}
	issues := (*fissues).List
	issuesContent := ""
	for _, e := range issues {
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
		pic := "http://i.annihil.us/u/prod/marvel/i/mg/b/40/image_not_available.jpg"
		if e.Pic != "" {
			pic = e.Pic
		}
		// TODO: For each comment
		commentsList := fmt.Sprintf(c["list"], e.Comments)

		con := fmt.Sprintf(c["issue-content"], name, e.PhaseID, e.SortID, pic, name,
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
	content := fmt.Sprintf(c["issues"], issues[0].PhaseID, issues[0].SortID, issues[0].Title, issuesContent)
	content = fmt.Sprintf(c["template"], phasesMenu[0], phasesMenu[1], phasesMenu[2], content)
	return content, nil
}

func getAboutPage(phases *service.PhaseList) string {
	phasesMenu := getPhasesMenuList(*phases, 3)
	return fmt.Sprintf(c["template"], phasesMenu[0], phasesMenu[1], phasesMenu[2], c["about"])
}

func getNotFoundPage(phasesMenu []string) string {
	return fmt.Sprintf(c["template"], phasesMenu[0], phasesMenu[1], phasesMenu[2], c["not-found"])
}

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
		    li := fmt.Sprintf(c["list-link"], link, phases[j].Name)
		    list = fmt.Sprintf("%s%s", list, li)
		}
		result[i] = list
		start = end
	}
	return result
}