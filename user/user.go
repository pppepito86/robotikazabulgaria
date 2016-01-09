package user

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"robotikazabulgaria/ws"
)

func readUsers() map[string]string {
	var users map[string]string
	file := ws.ReadFile("users.json")
	err := json.Unmarshal(file, &users)
	if err != nil {
		users = make(map[string]string)
		users["pesho"] = "test"
	}
	fmt.Println(users)
	return users
}

func RandomString() string {
	size := 32
	rb := make([]byte, size)
	rand.Read(rb)
	rs := base64.URLEncoding.EncodeToString(rb)
	fmt.Println("Generated id:", rs)
	return rs
}

func Authenticate(username, password string) bool {
	pass, found := readUsers()[username]
	return found && pass == password
}

func ContainsUser(username string) bool {
	_, found := readUsers()[username]
	return found
}
