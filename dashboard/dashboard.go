package dashboard

import (
	"robotikazabulgaria/hw"
	"robotikazabulgaria/teams"
	"robotikazabulgaria/admin"
)

type Dashboard struct {
	Name    string
	Challenge   admin.Challenge
	Homeworks map[string][]hw.Homework
}

func GetDashboard(user string) Dashboard {
	homeworks := hw.ReadHomeworks(user)
	dashboard := Dashboard{
		Name:    teams.GetTeamName(user),
		Challenge: admin.GetActiveChallenge(),
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
