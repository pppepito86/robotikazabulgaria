package dashboard

import (
	"robotikazabulgaria/hw"
	"robotikazabulgaria/teams"
	"robotikazabulgaria/admin"
	"time"
)

type Dashboard struct {
	Active bool
	Name    string
	Challenge   admin.Challenge
	Homeworks map[string][]hw.Homework
}

func GetDashboard(user string) Dashboard {
	homeworks := hw.ReadHomeworks(user)
	ch := admin.GetActiveChallenge()
	act := time.Now().UTC().Before(ch.EndTime.UTC())
	dashboard := Dashboard{
		Active: act,
		Name:    teams.GetTeamName(user),
		Challenge: ch,
		Homeworks: make(map[string][]hw.Homework),
	}
	for _, homework := range homeworks {
		if dashboard.Homeworks[homework.Task] == nil {
			dashboard.Homeworks[homework.Task] = make([]hw.Homework, 0)
		}
		dashboard.Homeworks[homework.Task] = append(dashboard.Homeworks[homework.Task], homework)
	}
	return dashboard
}
