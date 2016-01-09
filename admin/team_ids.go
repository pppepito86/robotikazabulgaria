package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"robotikazabulgaria/ws"
	"sync"
)

type TeamId struct {
	Id     string
	City   string
	School string
}

var teamIdMutex sync.Mutex

func GetTeamIds() []TeamId {
	var teamIds []TeamId
	file := ws.ReadFile("team_ids.json")
	err := json.Unmarshal(file, &teamIds)
	if err != nil {
		teamIds = make([]TeamId, 0)
	}
	fmt.Println(teamIds)
	return teamIds
}

func AddTeamId(id, city, school string) {
	teamIdMutex.Lock()
	defer teamIdMutex.Unlock()
	teamIds := GetTeamIds()
	teamIds = append(teamIds, TeamId{id, city, school})
	writeTeamIds(teamIds)
}

func writeTeamIds(teamIds []TeamId) {
	file := ws.GetFilePath("team_ids.json")
	json, _ := json.Marshal(teamIds)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}
