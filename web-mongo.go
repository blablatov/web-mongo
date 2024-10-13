// Web-server-gateway for data exchange with mongodb
// Веб-сервер-шлюз для обмена данными с mongodb

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Структура демаршалинга данных запроса.
// Structs for unmarshalling data
type mgoChat struct {
	User_1 struct {
		Uuid     string `json:"user_uuid"`
		Text     string `json:"text"`
		Datetime string `json:"datetime"`
	} `json:"user_1"`
	User_2 struct {
		Uuid     string `json:"user_uuid"`
		Text     string `json:"text"`
		Datetime string `json:"datetime"`
	} `json:"user_2"`
}

type result struct {
	inres1 []string
	inres2 []string
	finres string
	inerr  error
	setweb string
}

var dsnMongo string

// Инит dsn mongodb из конфига. Init dsn from config
func init() {
	dsnMongo = readMongoConf()
	log.Println("DsnMongo: ", dsnMongo)
}

func main() {

	log.SetPrefix("Server event: ")
	log.SetFlags(log.Lshortfile)

	LogInfo("web-server listening on :8017")

	// Мультиплексор запросов. Router of http-requests.
	mux := http.NewServeMux()
	mux.HandleFunc("/mgo", mgo)
	log.Fatal(http.ListenAndServe(":8017", mux))
}

// Хэндл декодера для создания диалога с внесением сообщений - возвращает чат айди
// Handle for create topic
func mgo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request of server:")

	// Параметры http-запроса. Parameters of headers
	fmt.Fprintf(w, "Method = %s\nURL = %s\nProto = %s\n", r.Method, r.URL, r.Proto)
	fmt.Printf("Method = %s\nURL = %s\nProto = %s\n", r.Method, r.URL, r.Proto)

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Printf("Host = %q\n", r.Host)

	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	fmt.Printf("RemoteAddr = %q\n", r.RemoteAddr)

	// Получение байтового среза запроса. Get byte slice of body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error get body: ", err)
	}
	//log.Println("reqBody: ", string(reqBody))

	var wg sync.WaitGroup
	//var mu sync.Mutex

	// Чек параметра id - структура SubData. Check `id` param.
	res := (string)([]byte(reqBody))

	//chdsn := make(chan string, 1)
	in := new(result)

	// // Получение dsn mongodb из конфига. Gets dsn from config
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	chdsn <- readMongoConf()
	// }()
	// mu.Lock()
	// dsnMongo := <-chdsn
	// mu.Unlock()
	// go func() {
	// 	wg.Wait()
	// 	close(chdsn)
	// }()

	// log.Println("DsnMongo", dsnMongo)

	// Чек параметра user_uuid - структура MgoChat. Check `user_uuid` param.
	if strings.Contains(res, "user_") {
		fmt.Println("user_- ОК!")

		chs := make(chan string, 7)
		chr := make(chan any, 1)

		wg.Add(1)
		go func() {
			defer wg.Done()

			///////////////////////
			uc := new(mgoChat)
			// Вызов метода MgoChat структуры демаршалинга. Calls unmarshal
			stu, err := uc.userChat(reqBody)
			if err != nil {
				log.Fatalf("Error unmarshal: %v", err)
			}
			log.Printf("Return userChat: %v", stu)

			///////////////////////
			// Вызов метода поиска sid в коллекции users mongodb.
			// Call of find method
			fd := new(findMongo)

			in.finres = fd.finder(dsnMongo, stu)
			if in.finres == "user_1 == user_2 - Second not found" ||
				in.finres == "user_1||user_2||both - Not found" {
				in.setweb = in.finres
			}
			if in.finres == "user_1 && user_2 - ОК!" {
				log.Println("Response find - OK!")

				///////////////////////
				// Вызов метода записи коллекций в mongodb.
				// Call of insert method
				ins := new(insertMongo)
				sin1, sin2, err := ins.inserter(dsnMongo, stu)
				if err != nil {
					log.Printf("Error insert data: %v\n", err)
					in.inerr = err
				}

				log.Printf("Result inserted: %v\t%v", sin1, sin2)
				in.inres1 = sin1
				in.inres2 = sin2
			}

			chs <- stu.User_1.Uuid
			chs <- stu.User_1.Text
			chs <- stu.User_1.Datetime

			chs <- stu.User_2.Uuid
			chs <- stu.User_2.Text
			chs <- stu.User_2.Datetime

		}()
		fmt.Fprintf(w, "Responses serv:\nUuid user_1: %v\nText user_1: %v\nDatetime user_1: %v\nUuid user_2: %v\nText user_2: %v\nDatetime user_2: %v\nResponse mongodb:\nError finder: %v\nError insert: %v\nCreated topic_1: %v\nCreated topic_2: %v\n",
			<-chs, <-chs, <-chs, <-chs, <-chs, <-chs, in.setweb, in.inerr, in.inres1, in.inres2)
		go func() {
			for range <-chs {
			}
			for range chr {
			}
			wg.Wait()
			close(chs)
			close(chr)
		}()
	}
}

// Демаршалинг тела запроса - структура MgoChat. Unmarshalling post
func (*mgoChat) userChat(reqBody []byte) (mgoChat, error) {

	var out mgoChat

	err := json.Unmarshal([]byte(reqBody), &out)
	if err != nil {
		log.Fatalf("Error unmarshal: %v", err)
	}

	// Объект json user_1
	log.Println(out.User_1.Uuid)
	log.Println(out.User_1.Text)
	log.Println(out.User_1.Datetime)

	// Объект json user_2
	log.Println(out.User_2.Uuid)
	log.Println(out.User_2.Text)
	log.Println(out.User_2.Datetime)

	return out, nil
}

// Logger
var logger = log.Default()

func LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}

// Func reads DSN from file the ./mongo.conf
func readMongoConf() string {
	var dsn string
	rf, err := os.Open("mongo.conf")
	if err != nil {
		log.Fatalf("Error open a config file mongo: %v", err)
	}
	defer rf.Close()
	input := bufio.NewScanner(rf)
	for input.Scan() {
		dsn = input.Text()
	}
	return dsn
}
