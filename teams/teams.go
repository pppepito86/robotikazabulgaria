package teams

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"robotikazabulgaria/ws"
	"strings"
	"time"
)

type Team struct {
	Name   string
	Pass   string
	City   string
	School string
	Id     string
	Time   time.Time
}

func RegisterTeam(team, pass1, pass2, city, school, id string) error {
	if err := checkTeamNameIsValid(team); err != nil {
		return err
	} else if err := checkPassIsValid(pass1, pass2); err != nil {
		return err
	} else if err := checkCityIsValid(city); err != nil {
		return err
	} else if err := checkSchoolIsValid(school); err != nil {
		return err
	} else if err := checkIdIsValid(id); err != nil {
		return err
	} else if err := checkTeamNameIsUnique(team); err != nil {
		return err
	} else if err := checkIdIsOK(id); err != nil {
		return err
	}
	t := Team{team, pass1, city, school, id, time.Now()}
	AddTeam(t)
	return nil
}

func checkTeamNameIsValid(team string) error {
	if len(team) == 0 {
		return errors.New("Не е въведено име на отбора")
	} else if len(team) < 3 {
		return errors.New("Името на отбора трябва да е поне 3 символа")
	} else if len(team) > 30 {
		return errors.New("Името на отбора не може да надвишава 30 символа")
	}

	for _, l := range team {
		if (l >= 'a' && l <= 'z') || (l >= 'A' && l <= 'Z') || (l >= 'а' && l <= 'я') || (l >= 'А' && l <= 'Я') {
			return nil
		}
	}

	return errors.New("Нямате букви в името на отбора")
}

func checkPassIsValid(pass1, pass2 string) error {
	if pass1 != pass2 {
		return errors.New("Паролите не съвпадат")
	} else if len(pass1) == 0 {
		return errors.New("Не е въведена парола")
	} else if len(pass1) < 8 {
		return errors.New("Паролата трябва да е поне 8 символа")
	} else if len(pass1) > 50 {
		return errors.New("Паролата не може да надвишава 50 символа")
	} else {
		return nil
	}
}

func checkCityIsValid(city string) error {
	if len(city) == 0 {
		return errors.New("Не е въведен град")
	} else if len(city) < 3 {
		return errors.New("Името на града трябва да е поне 3 символа")
	} else if len(city) > 30 {
		return errors.New("Името на града не може да надвишава 30 символа")
	} else {
		return nil
	}
}

func checkSchoolIsValid(school string) error {
	if len(school) == 0 {
		return errors.New("Не е въведено училище")
	} else if len(school) < 3 {
		return errors.New("Името на училището трябва да е поне 3 символа")
	} else if len(school) > 50 {
		return errors.New("Името на училището не трябва да надвишава 50 символа")
	} else {
		return nil
	}
}

func checkIdIsValid(id string) error {
	if len(id) == 0 {
		return errors.New("Не е въведен идентификационен номер")
	} else {
		return nil
	}
}

func checkTeamNameIsUnique(team string) error {
	team = strings.ToLower(team)
	teams := GetTeams()
	for _, t := range teams {
		if strings.ToLower(t.Name) == team {
			return errors.New("Вече съществува отбор с това име")
		}
	}
	if team == "pesho" || team == "marin" || team == "monica" {
		return errors.New("Вече съществува отбор с това име")
	}
	return nil
}

func checkIdIsOK(id string) error {
	teams := GetTeams()
	for _, t := range teams {
		if t.Id == id {
			return errors.New("Идентификационият номер не е валиден")
		}
	}
	teamIds := GetTeamIds()
	for _, t := range teamIds {
		if t.Id == id {
			return nil
		}
	}
	return errors.New("Id is not valid")
}

func GetTeams() []Team {
	var teams []Team
	file := ws.ReadFile("teams.json")
	err := json.Unmarshal(file, &teams)
	if err != nil {
		teams = make([]Team, 0)
	}
	fmt.Println(teams)
	return teams
}

func AddTeam(team Team) {
	teams := GetTeams()
	teams = append(teams, team)
	writeTeams(teams)
}

func writeTeams(teams []Team) {
	file := ws.GetFilePath("teams.json")
	json, _ := json.Marshal(teams)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}

func Authenticate(username, password string) bool {
	teams := GetTeams()
	for _, t := range teams {
		if (t.Name == username || t.Id == username) && t.Pass == password {
			return true
		}
	}
	return false
}

func GetRegisteredIds() map[string]bool {
	m := make(map[string]bool)
	tt := GetTeams()
	for _, t := range tt {
		m[t.Id] = true
	}
	return m
}

func GetTeamId(loginname string) string {
	teams := GetTeams()
	for _, t := range teams {
		if t.Name == loginname || t.Id == loginname {
			return t.Id
		}
	}
	return ""
}

func GetTeamName(id string) string {
	teams := GetTeams()
	for _, t := range teams {
		if t.Id == id {
			return t.Name
		}
	}
	return ""
}
