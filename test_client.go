package main

import (
	"YandexAlice/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("requests.json")
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	var requests []models.AliceRequest
	err = json.NewDecoder(file).Decode(&requests)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	for i, req := range requests {
		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("Ошибка сериализации запроса %d: %v", i, err)
			continue
		}

		resp, err := http.Post("http://localhost:8080/post", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Ошибка отправки запроса %d: %v", i, err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("\n--- Запрос %d ---\n", i+1)
		fmt.Println("➡ Отправлено:")
		fmt.Println(string(data))
		fmt.Println("⬅ Ответ:")
		fmt.Println(string(body))
	}
}
