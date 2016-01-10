package dashboard

import "robotikazabulgaria/hw"

type Dashboard struct {
	Team    []hw.Homework
	Project []hw.Homework
	Robot   []hw.Homework
}

func GetDashboard(user string) Dashboard {
	homeworks := hw.ReadHomeworks(user)
	dashboard := Dashboard{
		Team:    make([]hw.Homework, 0),
		Project: make([]hw.Homework, 0),
		Robot:   make([]hw.Homework, 0)}
	for _, homework := range homeworks {
		if homework.Task == "team" {
			dashboard.Team = append(dashboard.Team, homework)
		} else if homework.Task == "project" {
			dashboard.Project = append(dashboard.Project, homework)
		} else if homework.Task == "robot" {
			dashboard.Robot = append(dashboard.Robot, homework)
		}
	}
	return dashboard
}
