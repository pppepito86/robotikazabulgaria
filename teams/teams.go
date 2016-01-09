package teams

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"robotikazabulgaria/admin"
	"robotikazabulgaria/ws"
)

type Team struct {
	Name   string
	Pass   string
	City   string
	School string
	Id     string
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
	t := Team{team, pass1, city, school, id}
	AddTeam(t)
	return nil
}

func checkTeamNameIsValid(team string) error {
	if len(team) == 0 {
		return errors.New("Team name is not set")
	} else if len(team) < 3 {
		return errors.New("Team name length should be at least 3 symbols")
	} else if len(team) > 20 {
		return errors.New("Team name length should not exceed 20 symbols")
	} else {
		return nil
	}
}

func checkPassIsValid(pass1, pass2 string) error {
	if pass1 != pass2 {
		return errors.New("Passwords do not match")
	} else if len(pass1) == 0 {
		return errors.New("Password is not set")
	} else if len(pass1) < 8 {
		return errors.New("Password length should be at least 8 symbols")
	} else if len(pass1) > 50 {
		return errors.New("Password length should not exceed 50 symbols")
	} else {
		return nil
	}
}

func checkCityIsValid(city string) error {
	if len(city) == 0 {
		return errors.New("City is not set")
	} else if len(city) < 3 {
		return errors.New("City length should be at least 3 symbols")
	} else if len(city) > 20 {
		return errors.New("City length should not exceed 20 symbols")
	} else {
		return nil
	}
}

func checkSchoolIsValid(school string) error {
	if len(school) == 0 {
		return errors.New("School is not set")
	} else if len(school) < 3 {
		return errors.New("School length should be at least 3 symbols")
	} else if len(school) > 20 {
		return errors.New("School length should not exceed 20 symbols")
	} else {
		return nil
	}
}

func checkIdIsValid(id string) error {
	if len(id) == 0 {
		return errors.New("Identification number is not set")
	} else {
		return nil
	}
}

func checkTeamNameIsUnique(team string) error {
	teams := GetTeams()
	for _, t := range teams {
		if t.Name == team {
			return errors.New("Team name already exists")
		}
	}
	return nil
}

func checkIdIsOK(id string) error {
	teams := GetTeams()
	for _, t := range teams {
		if t.Id == id {
			return errors.New("Id is not valid")
		}
	}
	teamIds := admin.GetTeamIds()
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
