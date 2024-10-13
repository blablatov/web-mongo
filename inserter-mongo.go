// Модуль записи коллекций в mongodb
// Inserts collections to mongodb

package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insertMongo struct {
	gid1   string
	gid2   string
	id1    string `json:"user_uuid"`
	id2    string `json:"user_uuid"`
	intop2 []string
	insub1 []string
	insub2 []string
}

var inserts = []insertMongo{}

// TODO мапа для возвращаемых данных. Map to return data
var (
	mu     sync.Mutex
	insMap = make(map[string]insertMongo)
)

func (in *insertMongo) inserter(dsnMongo string, stu mgoChat) ([]string, []string, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsnMongo))

	// Отложенное отключение, после создания клиента. Waiting disconnect
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Println("Inserter_dsnMongo: ", dsnMongo)

	var wg sync.WaitGroup
	m := new(insertMongo)

	//////////////////////////////
	// Сгенерирует случайную строку _id топика user_1 .
	// Generates random string topics name
	key := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(key))
	m.gid1 = genTopicName(rng)
	if m.gid1 == "" {
		log.Println("Error of generate string")
		return nil, nil, nil
	}
	log.Println("Random string: ", m.gid1)

	//////////////////////////////
	// Сгенерирует случайную строку _id топика user_2 .
	// Generates random string topics name
	key = time.Now().UTC().UnixNano()
	rng = rand.New(rand.NewSource(key))
	m.gid2 = genTopicName(rng)
	if m.gid2 == "" {
		log.Println("Error of generate string")
		return nil, nil, nil
	}
	log.Println("Random string: ", m.gid2)

	//////////////////////////////
	// Поиск _id для user_1 и user_2 в БД, подставляется в документ topics
	// Seeks _id for user_1 && user_2 to mongo
	fd := new(findMongo)
	in.id1, in.id2, err = fd.idFinder(dsnMongo, stu)
	if err != nil {
		log.Printf("Error find of _id: %v", err)
		return nil, nil, err
	}

	// Создает обработчик коллекций. Creates handle of collections
	cn := client.Database("gotest").Collection("topics")

	// Создает топик user_1. Creates topic user_1
	// topic := []interface{}{
	// 	bson.M{"_id": m.gid1, // генерируется случайный _id
	// 		"createdat": bson.M{
	// 			"$date": stu.User_1.Datetime,
	// 		},
	// 		"updatedat": bson.M{
	// 			"$date": stu.User_1.Datetime,
	// 		},
	// 		"state": 0,
	// 		"touchedat": bson.M{
	// 			"$date": stu.User_1.Datetime,
	// 		},
	// 		"usebt": false,
	// 		"owner": id1, // подставляется _id1 из users
	// 		"access": bson.M{
	// 			"auth": bson.M{
	// 				"$numberLong": "47",
	// 			},
	// 			"anon": bson.M{
	// 				"$numberLong": "0",
	// 			},
	// 		},
	// 		"seqid": 1,
	// 		"delid": 0,
	// 		"public": bson.M{
	// 			"fn":   "gotov_chat",
	// 			"note": "тестовый топик",
	// 		},
	// 		"trusted": nil,
	// 		"tags": bson.A{
	// 			stu.User_1.Useruuid, // подставляется полученнный sid
	// 		},
	// 	},
	// }

	topic := []interface{}{
		bson.M{"_id": m.gid1,
			"access": bson.M{
				"anon": 47,
				"auth": 0,
			},
			"delid":         0,
			"createdat":     stu.User_1.Datetime,
			"lastmessageat": stu.User_1.Datetime,
			"owner":         in.id1,
			"public": bson.M{
				"fn":   "gotov_chat",
				"note": "hello chat",
			},
			//"tags":      stu.User_1.Useruuid,
			"tags": bson.A{
				stu.User_1.Uuid, // подставляется полученнный sid
			},
			"seqid":     0,
			"state":     0,
			"stateat":   nil,
			"updatedat": "2019-10-11T12:13:14.522Z",
			"usebt":     false,
		},
	}

	var res *mongo.InsertManyResult
	// Вставка документа topics для user_1
	// Inserts documents into the collection
	// Set the Ordered option to false to allow
	// both operations to happen even if one of them errors.
	opts := options.InsertMany().SetOrdered(false)
	res, err = cn.InsertMany(context.TODO(), topic, opts)
	if err != nil {
		log.Printf("Error of insert topics: %v", err)
		return nil, nil, err
	}
	log.Printf("Inserted topic %v\n", res.InsertedIDs)

	// Срез строковых значений с размером результата. Slice len of cursor
	top1 := make([]string, 0, len(res.InsertedIDs))

	// Формирование строкового слайса. Gets []string slice
	for _, v := range res.InsertedIDs {
		if v != nil {
			top1 = append(top1, v.(string))
		}
	}

	//////////////////////////////
	// Вставка документа topics для user_2
	chtp2 := make(chan []string, 1)
	var top2 []string

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Создает топик user_2. Creates topic user_2
		// topic2 := []interface{}{
		// 	bson.M{"_id": m.gid2, //TODO генератор своей строки
		// 		"createdat": bson.M{
		// 			"$date": stu.User_2.Datetime,
		// 		},
		// 		"updatedat": bson.M{
		// 			"$date": stu.User_2.Datetime,
		// 		},
		// 		"state": 0,
		// 		"touchedat": bson.M{
		// 			"$date": stu.User_2.Datetime,
		// 		},
		// 		"usebt": false,
		// 		"owner": id2, // подставляется _id2 из users
		// 		"access": bson.M{
		// 			"auth": bson.M{
		// 				"$numberLong": "47",
		// 			},
		// 			"anon": bson.M{
		// 				"$numberLong": "0",
		// 			},
		// 		},
		// 		"seqid": 0,
		// 		"delid": 0,
		// 		"public": bson.M{
		// 			"fn":   "gotov_chat",
		// 			"note": "тестовый топик",
		// 		},
		// 		"trusted": nil,
		// 		"tags": bson.A{
		// 			stu.User_2.Useruuid,
		// 		},
		// 	},
		// }

		topic2 := []interface{}{
			bson.M{"_id": m.gid2,
				"access": bson.M{
					"anon": 47,
					"auth": 0,
				},
				"delid":         0,
				"createdat":     stu.User_2.Datetime,
				"lastmessageat": stu.User_2.Datetime,
				"owner":         in.id2,
				"public": bson.M{
					"fn":   "gotov_chat",
					"note": "hello chat",
				},
				//"tags":      stu.User_2.Useruuid,
				"tags": bson.A{
					stu.User_2.Uuid, // подставляется полученнный sid
				},
				"seqid":     0,
				"state":     0,
				"stateat":   nil,
				"updatedat": stu.User_2.Datetime,
				"usebt":     false,
			},
		}

		res := new(mongo.InsertManyResult)
		// Вставка документа topics для user_2
		// Inserts documents into the collection for user_2
		opts := options.InsertMany().SetOrdered(false)
		res, err = cn.InsertMany(context.TODO(), topic2, opts)
		if err != nil {
			log.Printf("Error of insert topics: %v", err)
		}
		log.Printf("Inserted topic %v\n", res.InsertedIDs)

		// Срез строковых значений с размером результата. Slice len of cursor
		mir := make([]string, 0, len(res.InsertedIDs))

		// Формирование строкового слайса. Gets []string slice
		for _, v := range res.InsertedIDs {
			if v != nil {
				mir = append(mir, v.(string))
			}
		}

		chtp2 <- mir
	}()
	top2 = <-chtp2
	go func() {
		wg.Wait()
		close(chtp2)
	}()

	//////////////////////////////
	// Вставка документа subscriptions для user_1
	chsb1 := make(chan []string, 1)
	//var sb1 []string

	wg.Add(1)
	go func() {
		defer wg.Done()

		//Создает обработчик коллекций. Creates handle of collections
		cn := client.Database("gotest").Collection("subscriptions")

		// subscript := []interface{}{
		// 	bson.M{"_id": m.gid1 + ":" + id1, // подставляется сгенеренный _id (topic) и _id из users
		// 		"createdat": bson.M{
		// 			"$date": stu.User_1.Datetime,
		// 		},
		// 		"updatedat": bson.M{
		// 			"$date": stu.User_1.Datetime,
		// 		},
		// 		"user":      id1,    // подставляется _id из users
		// 		"topic":     m.gid1, //подставляется сгенеренный _id (topic)
		// 		"delid":     0,
		// 		"recvseqid": 1,
		// 		"readseqid": 1,
		// 		"modewant": bson.M{
		// 			"$numberLong": "255",
		// 		},
		// 		"modegiven": bson.M{
		// 			"$numberLong": "255",
		// 		},
		// 		"private": bson.M{
		// 			"comment": "тест топика-темы",
		// 		},
		// 	},
		// }

		subscript := []interface{}{
			bson.M{
				"_id":       m.gid1 + ":" + in.id1, // подставляется сгенеренный _id (topic) и _id из users
				"createdat": stu.User_1.Datetime,
				"updatedat": stu.User_1.Datetime,
				"deletedat": nil,
				"user":      in.id1, // подставляется _id из users
				"topic":     m.gid1, //подставляется сгенеренный _id (topic)
				"recvseqid": 0,
				"readseqid": 0,
				"modewant":  255,
				"modegiven": 255,
				//"private":   "About me",
				"private": bson.M{
					"comment": "About me",
				},
				"state": 0,
			},
		}

		resu := new(mongo.InsertManyResult)
		// Создает коллекцию подписок. Creates subscriptions collection
		opts := options.InsertMany().SetOrdered(false)
		resu, err = cn.InsertMany(context.TODO(), subscript, opts)
		if err != nil {
			log.Printf("Error of insert subscriptions: %v", err)
		}
		log.Printf("Inserted subscriptions: %v\n", resu.InsertedIDs)

		// Срез строковых значений с размером результата. Slice len of result
		mis := make([]string, 0, len(resu.InsertedIDs))

		// Формирование строкового слайса. Gets []string slice
		for _, v := range resu.InsertedIDs {
			if v != nil {
				mis = append(mis, v.(string))
			}
		}

		chsb1 <- mis
	}()
	<-chsb1
	go func() {
		wg.Wait()
		close(chsb1)
	}()

	//////////////////////////////
	// Вставка документа subscriptions для user_2
	chsb2 := make(chan []string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		//Создает обработчик коллекций. Creates handle of collections
		cn := client.Database("gotest").Collection("subscriptions")

		// subscript := []interface{}{
		// 	bson.M{"_id": m.gid2 + ":" + id2, // подставляется сгенеренный _id (topic) и _id из users
		// 		"createdat": bson.M{
		// 			"$date": stu.User_2.Datetime,
		// 		},
		// 		"updatedat": bson.M{
		// 			"$date": stu.User_2.Datetime,
		// 		},
		// 		"user":      id2,    // подставляется _id из users
		// 		"topic":     m.gid2, //подставляется сгенеренный _id (topic)
		// 		"delid":     0,
		// 		"recvseqid": 1,
		// 		"readseqid": 1,
		// 		"modewant": bson.M{
		// 			"$numberLong": "255",
		// 		},
		// 		"modegiven": bson.M{
		// 			"$numberLong": "255",
		// 		},
		// 		"private": bson.M{
		// 			"comment": "тест топика-темы",
		// 		},
		// 	},
		// }

		subscript := []interface{}{
			bson.M{
				"_id":       m.gid2 + ":" + in.id2, // подставляется сгенеренный _id (topic) и _id из users
				"createdat": stu.User_2.Datetime,
				"updatedat": stu.User_2.Datetime,
				"deletedat": nil,
				"user":      in.id2, // подставляется _id из users
				"topic":     m.gid2, //подставляется сгенеренный _id (topic)
				"recvseqid": 0,
				"readseqid": 0,
				"modewant":  255,
				"modegiven": 255,
				//"private":   "About me",
				"private": bson.M{
					"comment": "About me too",
				},
				"state": 0,
			},
		}

		resu := new(mongo.InsertManyResult)
		// Создает коллекцию подписок. Creates subscriptions collection
		opts := options.InsertMany().SetOrdered(false)
		resu, err := cn.InsertMany(context.TODO(), subscript, opts)
		if err != nil {
			log.Printf("Error of insert subscriptions: %v", err)
		}
		log.Printf("Inserted subscriptions: %v\n", resu.InsertedIDs)

		// Срез строковых значений с размером результата. Slice len of result
		mis := make([]string, 0, len(resu.InsertedIDs))

		// Формирование строкового слайса. Gets []string slice
		for _, v := range resu.InsertedIDs {
			if v != nil {
				mis = append(mis, v.(string))
			}
		}

		chsb2 <- mis
	}()
	<-chsb2
	go func() {
		wg.Wait()
		close(chsb2)
	}()

	return top1, top2, nil
}

// Генератор случайных строк
func genTopicName(rng *rand.Rand) string {
	n := rng.Intn(11) // Случайная длина до 11
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000)) // Случайная руна до '\u0999'
		runes[i] = r
		runes[n-1-i] = r
	}
	log.Println("genTopicName: ", string(runes))
	return string(runes)
}
