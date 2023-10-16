package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type User struct {
	id       string
	InitTime time.Time
}

type Data struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	PIN   string `json:"PIN"`
}

type ResponseData struct {
	Data Data `json:"data"`
}

type Response struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
	Index   string       `json:"index"`
}

var userData = make(map[string]*Data)

func initEndPoint(user User, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	user.InitTime = time.Now()
	initURL := "http://localhost:4000/user/" + user.id
	req, err := http.NewRequest("GET", initURL, nil)
	if err != nil {
		fmt.Println("Error creating request")
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var responseObject Response
	err = json.Unmarshal([]byte(string(body)), &responseObject)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	mu.Lock()
	userData[user.id] = &responseObject.Data.Data
	mu.Unlock()
	fmt.Println(responseObject.Data.Data.Name, time.Now().Sub(user.InitTime))
	if responseObject.Data.Data.PIN != "" && (time.Now().Sub(user.InitTime)) <= 5*time.Second {
		pay(responseObject.Data.Data.Name)
	} else {
		fmt.Println("Transaction done by the user ", responseObject.Data.Data.Name, "timed out")
	}
}

func pay(name string) {
	fmt.Println(name, ":Transaction done successfully")
}

func main() {
	users := []User{
		{id: "652a35962bd1257661a14f1b"},
		{id: "652a35c32bd1257661a14f1d"},
		{id: "652a3a7b42f56b3875a800dc"},
		{id: "652a4c055b21625d6be302f6"},
		{id: "652a4c2d5b21625d6be302f8"},
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, user := range users {
		wg.Add(1)
		go initEndPoint(user, &wg, &mu)
	}
	wg.Wait()
}
