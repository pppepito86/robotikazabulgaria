package dashboard

import (
	"robotikazabulgaria/hw"
	"robotikazabulgaria/teams"
	"robotikazabulgaria/admin"
	"time"
	"fmt"
)

type Dashboard struct {
	Active bool
	Name    string
	Challenge   admin.Challenge
	Homeworks map[string][]hw.Homework
	Marks map[string]admin.Mark
	Teams []teams.Team
}

func GetDashboard(user string) Dashboard {
	homeworks := hw.ReadHomeworks(user)
	tms := admin.GetTeamMarks("pesho")
	t := tms[user]
	ch := admin.GetActiveChallenge()
	act := time.Now().UTC().Before(ch.EndTime.UTC())
	dashboard := Dashboard{
		Active: act,
		Name:    teams.GetTeamName(user),
		Challenge: ch,
		Homeworks: make(map[string][]hw.Homework),
		Marks: t.Marks,
	}
	for _, homework := range homeworks {
		if dashboard.Homeworks[homework.Task] == nil {
			dashboard.Homeworks[homework.Task] = make([]hw.Homework, 0)
		}
		dashboard.Homeworks[homework.Task] = append(dashboard.Homeworks[homework.Task], homework)
	}
	return dashboard
}

func GetHistoryDashboard(user string, team string) Dashboard {
	if len(team) > 0 {
		user = team
	}
	homeworks := hw.ReadHomeworks(user)
	tms := admin.GetTeamMarks("pesho")
	t := tms[user]
	//m := t.Marks[taskname]
	
	fmt.Println("team marks", t)
	chs := admin.GetChallenges()
	ch := admin.Challenge{}
	for _, ccc := range chs.Challenges {
		if ccc.State == "finished" {
			ch = ccc
			break
		}
	}
	if ch.Id == "" {
		return Dashboard{}
	}
	act := time.Now().UTC().Before(ch.EndTime.UTC())
	dashboard := Dashboard{
		Active: act,
		Name:    teams.GetTeamName(user),
		Challenge: ch,
		Homeworks: make(map[string][]hw.Homework),
		Marks: t.Marks,
		Teams: teams.GetTeams(),
	}
	for _, homework := range homeworks {
		if dashboard.Homeworks[homework.Task] == nil {
			dashboard.Homeworks[homework.Task] = make([]hw.Homework, 0)
		}
		dashboard.Homeworks[homework.Task] = append(dashboard.Homeworks[homework.Task], homework)
	}

	return dashboard
}

