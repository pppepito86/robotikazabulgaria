package hw

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"robotikazabulgaria/ws"
)

type Homework struct {
	Filename    string
	Link        string
	Description string
	Task        string
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
