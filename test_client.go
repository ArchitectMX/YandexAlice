package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TestRequest struct {
	Name string          `json:"name"`
	Body json.RawMessage `json:"body"`
}

func main() {
	file, err := os.Open("requests.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var requests []TestRequest
	json.Unmarshal(byteValue, &requests)

	for _, req := range requests {
		fmt.Println("==>", req.Name)
		resp, err := http.Post("http://localhost:8080/post", "application/json", bytes.NewBuffer(req.Body))
		if err != nil {
			fmt.Println("Request failed:", err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println(string(body))
		fmt.Println("------")
	}
}
