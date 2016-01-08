package hw

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"robotikazabulgaria/ws"
	"strconv"
	"time"
)

type Homework struct {
	Filename    string
	Link        string
	Description string
	Task        string
	Time        time.Time
}

func ReadHomeworks(user string) []Homework {
	var homeworks []Homework
	file := ws.ReadFile(user, "homework.json")
	json.Unmarshal(file, &homeworks)
	return homeworks
}

func WriteHomeworks(user string, homework []Homework) {
	file := ws.GetFilePath(user, "homework.json")
	json, _ := json.Marshal(homework)
	os.Create(file)
	ioutil.WriteFile(file, json, 0700)
}

func AddHomework(user string, homework Homework) {
	homeworks := ReadHomeworks(user)
	homeworks = append(homeworks, homework)
	WriteHomeworks(user, homeworks)
}

func DeleteHomework(user string, id string) {
	homeworks := ReadHomeworks(user)
	for i, homework := range homeworks {
		if id == strconv.FormatInt(homework.Time.UnixNano(), 10) {
			homeworks = append(homeworks[:i], homeworks[i+1:]...)
		}
	}
	WriteHomeworks(user, homeworks)
}
