// Запустить сервер ./web-mongo
// Выполнить go test .

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestWebMongo(t *testing.T) {

	// URL тестового сорвера локально. Для облака указать внешний IP ВМ.
	apiUrl := "http://localhost:8017/mgo"

	// Формирование json параметров запроса. JSON params of request
	payload, _ := json.Marshal(struct {
		user_uuid string `json:"user_uuid"`
		text      string `json:"text"`
		datetime  string `json:"datetime"`
	}{
		user_uuid: "3e266244-0e23-4f2e-8cb5-b4d118054777",
		text:      "Hello",
		datetime:  "timestamp",
	})

	// Форматирование запроса. Formatting of the request
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
	if err != nil {
		t.Log(err)
	}
	// Формирование заголовков запроса. Headers of request
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Формирование метаданных структуры запроса. Struct of request
	client := &http.Client{
		Transport: &http.Transport{},
	}

	resp, err := client.Do(req) // Выполнение запроса. Send of request
	if err != nil {
		t.Logf("Error on response.%v\n[ERROR] -", err)
	}

	t.Logf("Status = %v ", "Ok") // Статус ответа сервера. Status of response

	// Чтение данных сервера, обработка ошибок. Reads data from server, check errors
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("Error while reading the response bytes:", err)
	}
	t.Logf("\nResponse of server: %v\n", string([]byte(body)))
}

func BenchmarkPost(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 10; i++ {

		apiUrl := "http://localhost:8017/mgo"
		payload, _ := json.Marshal(struct {
			user_uuid string `json:"user_uuid"`
			text      string `json:"text"`
			datetime  string `json:"datetime"`
		}{
			user_uuid: "3e266244-0e23-4f2e-8cb5-b4d118054777",
			text:      "Hello",
			datetime:  "timestamp",
		})

		client := &http.Client{
			Transport: &http.Transport{},
		}

		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
		if err != nil {
			b.Log(err)
		}
		req.Header.Set("Content-type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Error of request: %v", err)
		}

		defer resp.Body.Close()

		b.Logf("Status = %v\n", resp.Status)
	}
}
