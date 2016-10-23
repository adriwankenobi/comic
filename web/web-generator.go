package web

import (
	"fmt"
	"github.com/adriwankenobi/comic/service"
)

var htmlFiles = []string{
	"template",
	"tabs",
	"tab-li",
	"tab-content",
	"tab-content-intro",
	"tab-content-phase",
	"clear-fix",
	"issues",
	"issue-content",
	"not-found",
}

func getIndexPage(issues *service.FissuesList) (string, error) {
	introID := "intro"
	issuesLi := fmt.Sprintf(c["tab-li"], "active", introID, introID, introID, "Intro")
	issuesContent := fmt.Sprintf(c["tab-content"], "active", introID, introID, introID, c["tab-content-intro"])
	for _, e := range *issues {
		phaseID := fmt.Sprintf("phase%s", e.Phase.ID)
		li := fmt.Sprintf(c["tab-li"], "", phaseID, phaseID, phaseID, e.Phase.Name)
		issuesLi = fmt.Sprintf("%s%s", issuesLi, li)
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
	content = fmt.Sprintf(c["template"], content)
	return content, nil
}

func getIssuesPage(issues *service.ComicList) (string, error) {
	if issues.IsEmpty() {
		return notFound(), nil
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
			e.Characters,
			e.Creators,
			displayComments,
			e.Comments,
		)
		issuesContent = fmt.Sprintf("%s%s", issuesContent, con)
	}
	issuesContent = fmt.Sprintf("%s%s", issuesContent, c["clear-fix"])
	content := fmt.Sprintf(c["issues"], (*issues)[0].PhaseID, (*issues)[0].SortID, (*issues)[0].Title, issuesContent)
	content = fmt.Sprintf(c["template"], content)
	return content, nil
}

func notFound() string {
	return fmt.Sprintf(c["template"], c["not-found"])
}
