package dashboard

import (
	"robotikazabulgaria/admin"
	"robotikazabulgaria/hw"
	"robotikazabulgaria/teams"
	"time"
)

type Dashboard struct {
	Active    bool
	Name      string
	Challenge admin.Challenge
	Homeworks map[string][]hw.Homework
	Marks     map[string]admin.Mark
	Teams     []teams.Team
}

func GetDashboard(user string) Dashboard {
	homeworks := hw.ReadHomeworks(user)
	tms := admin.GetTeamMarks("pesho")
	t := tms[user]
	ch := admin.GetActiveChallenge()
	act := time.Now().UTC().Before(ch.EndTime.UTC())
	dashboard := Dashboard{
		Active:    act,
		Name:      teams.GetTeamName(user),
		Challenge: ch,
		Homeworks: make(map[string][]hw.Homework),
		Marks:     t.Marks,
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
	origUser := user
	if len(team) > 0 {
		user = team
	}
	homeworks := hw.ReadHomeworks(user)
	tms := admin.GetTeamMarks("pesho")
	t := tms[user]
	//m := t.Marks[taskname]

	chs := admin.GetChallenges()
	ch := admin.Challenge{}
	for iii := len(chs.Challenges) - 1; iii >= 0; iii-- {
		ccc := chs.Challenges[iii]
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
		Active:    act,
		Name:      teams.GetTeamName(user),
		Challenge: ch,
		Homeworks: make(map[string][]hw.Homework),
		Teams:     teams.GetTeams("0"),
	}
	if origUser == user {
		dashboard.Marks = t.Marks
	}
	for _, homework := range homeworks {
		if dashboard.Homeworks[homework.Task] == nil {
			dashboard.Homeworks[homework.Task] = make([]hw.Homework, 0)
		}
		dashboard.Homeworks[homework.Task] = append(dashboard.Homeworks[homework.Task], homework)
	}

	return dashboard
}
