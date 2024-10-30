// Модуль поиска по `user_uuid` пользователей в mongodb
// Checks user_uuid registered users

package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type findMongo struct{}

var finder = func(dsnMongo string, stu mgoChat) string {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsnMongo))

	// Отложенный дисконнект, после создания клиента
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Инит коллекции. Initial collection
	cn := client.Database("gotest").Collection("users")

	// Поиск документов сортировкой по параметру
	// Founds docs with sort on param
	cursor, err := cn.Distinct(context.TODO(), "tags", bson.D{})
	if err != nil {
		log.Printf("Error find of tags: %v", err)
	}

	// Срез строковых значений с размером курсора. Slice len of cursor
	mp := make([]string, 0, len(cursor))

	// Формирование строкового слайса. Gets []string slice
	for _, v := range cursor {
		if v != nil {
			mp = append(mp, v.(string))
		}
	}

	if stu.User_1.Uuid == stu.User_2.Uuid {
		log.Println("user_1 == user_2 - Second not found")
		return "user_1 == user_2 - Second not found"
	}

	if strings.Contains(strings.Join(mp, ""), stu.User_1.Uuid) &&
		strings.Contains(strings.Join(mp, ""), stu.User_2.Uuid) {
		//log.Println("user_1 && user_2 - ОК!")
		return "user_1 && user_2 - ОК!"
	} else {
		log.Println("user_1||user_2||both - Not found")
		return "user_1||user_2||both - Not found"
	}
}

var idFinder = func(dsnMongo string, stu mgoChat) (string, string, error) {

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsnMongo))

	// Отложенный дисконнект, после создания клиента
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Инит коллекции. Initial collection
	cn := client.Database("gotest").Collection("users")

	// Поиск _id по параметру tags для user_1.
	// Opts задаваемое время выполнения на сервере.
	// Find all unique values for the "_id" field for documents in which the
	// "tags" field is user_1.
	// MaxTime указывает время выполнения операции на сервере.
	// Specify the MaxTime option to limit the amount of time the operation can
	// run on the server.
	filter := bson.D{{"tags", bson.D{{"$eq", stu.User_1.Uuid}}}}
	opts := options.Distinct().SetMaxTime(2 * time.Second)
	values, err := cn.Distinct(context.TODO(), "_id", filter, opts)
	if err != nil {
		log.Printf("Error find of id1: %v", err)
	}

	// Срез строковых значений с размером курсора. Slice len of cursor
	mp := make([]string, 0, len(values))

	// Формирование строкового слайса. Gets []string slice
	for _, v := range values {
		if v != nil {
			mp = append(mp, v.(string))
		}
	}
	log.Println("Id of user_1", mp)
	id1 := strings.Join(mp, "")

	///////////////////////////
	var wg sync.WaitGroup
	chid := make(chan string, 1)
	var id2 string

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Поиск _id по параметру tags для user_2
		// Founds _id with sort for user_2
		filter := bson.D{{"tags", bson.D{{"$eq", stu.User_2.Uuid}}}}
		opts := options.Distinct().SetMaxTime(2 * time.Second)
		values, err := cn.Distinct(context.TODO(), "_id", filter, opts)
		if err != nil {
			log.Printf("Error find of id2: %v", err)
		}

		// Срез строковых значений с размером курсора. Slice len of cursor
		mp := make([]string, 0, len(values))

		// Формирование строкового слайса. Gets []string slice
		for _, v := range values {
			if v != nil {
				mp = append(mp, v.(string))
			}
		}
		log.Println("Id of user_2", mp)

		id2 = strings.Join(mp, "")
		chid <- id2

	}()

	id2 = <-chid
	go func() {
		wg.Wait()
		close(chid)
	}()

	return id1, id2, err
}
