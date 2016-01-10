package teams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"robotikazabulgaria/ws"
	"sort"
	"strconv"
	"sync"
)

type TeamId struct {
	Id     string
	City   string
	School string
}

type TeamIdArray []TeamId

func (slice TeamIdArray) Len() int {
	return len(slice)
}

func (slice TeamIdArray) Less(i, j int) bool {
	a, _ := strconv.ParseInt(slice[i].Id, 10, 64)
	b, _ := strconv.ParseInt(slice[j].Id, 10, 64)
	return a < b
}

func (slice TeamIdArray) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

var teamIdMutex sync.Mutex

func GetTeamIds() TeamIdArray {
	var teamIds TeamIdArray
	file := ws.ReadFile("team_ids.json")
	err := json.Unmarshal(file, &teamIds)
	if err != nil {
		teamIds = make([]TeamId, 0)
	}
	fmt.Println(teamIds)
	sort.Sort(teamIds)
	return teamIds
}

func AddTeamId(id, city, school string) {
	if len(id) == 0 || len(city) == 0 || len(school) == 0 {
		return
	}
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

type TeamIdInfo struct {
	Id         string
	City       string
	School     string
	Registered bool
}

func GetTeamsIdInfo() []TeamIdInfo {
	r := make([]TeamIdInfo, 0)
	tt := GetTeamIds()
	m := GetRegisteredIds()
	for _, t := range tt {
		teamInfo := TeamIdInfo{t.Id, t.City, t.School, m[t.Id]}
		r = append(r, teamInfo)
	}
	return r
}
