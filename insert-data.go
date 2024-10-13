// Структурные литералы топиков и подпискок
// Struct literals topics and subscriptions
// TODO - в тиноде генерируется случайный _id топика, при чтении топика выполняется его валидация, при несоответствии выдается ошибка - `meta.Get.Sub failed: error decoding key access.auth: cannot decode embedded document into an integer type`

package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type TopicData struct {
	ID        string `json:"_id"`
	Createdat struct {
		Date time.Time `json:"$date"`
	} `json:"createdat"`
	Updatedat struct {
		Date time.Time `json:"$date"`
	} `json:"updatedat"`
	State     int `json:"state"`
	Touchedat struct {
		Date time.Time `json:"$date"`
	} `json:"touchedat"`
	Usebt  bool   `json:"usebt"`
	Owner  string `json:"owner"`
	Access struct {
		Auth struct {
			NumberLong string `json:"$numberLong"`
		} `json:"auth"`
		Anon struct {
			NumberLong string `json:"$numberLong"`
		} `json:"anon"`
	} `json:"access"`
	Seqid  int `json:"seqid"`
	Delid  int `json:"delid"`
	Public struct {
		Fn   string `json:"fn"`
		Note string `json:"note"`
	} `json:"public"`
	Trusted any      `json:"trusted"`
	Tags    []string `json:"tags"`
}

var (
	topic = []interface{}{
		bson.M{"_id": "grpL6GdlEDQ3DM", // в Тиноде генерируется случайный _id
			"createdat": bson.M{
				"$date": "2024-07-11T05:57:51.202Z",
			},
			"updatedat": bson.M{
				"$date": "2024-07-11T05:57:51.202Z",
			},
			"state": 0,
			"touchedat": bson.M{
				"$date": "2024-07-11T05:57:51.202Z",
			},
			"usebt": false,
			"owner": "HhlZaX2A80Y", // подставляется _id из users
			"access": bson.M{
				"auth": bson.M{
					"$numberLong": "47",
				},
				"anon": bson.M{
					"$numberLong": "0",
				},
			},
			"seqid": 0,
			"delid": 0,
			"public": bson.M{
				"fn":   "gotov_chat",
				"note": "тестовый топик",
			},
			"trusted": nil,
			"tags": bson.A{
				"3e266244-0e23-4f2e-8cb5-b4d118054777", // подставляется полученнный sid
			},
		},
	}
)

type SubscriptionsData struct {
	ID        string `json:"_id"`
	Createdat struct {
		Date time.Time `json:"$date"`
	} `json:"createdat"`
	Updatedat struct {
		Date time.Time `json:"$date"`
	} `json:"updatedat"`
	User      string `json:"user"`
	Topic     string `json:"topic"`
	Delid     int    `json:"delid"`
	Recvseqid int    `json:"recvseqid"`
	Readseqid int    `json:"readseqid"`
	Modewant  struct {
		NumberLong string `json:"$numberLong"`
	} `json:"modewant"`
	Modegiven struct {
		NumberLong string `json:"$numberLong"`
	} `json:"modegiven"`
	Private struct {
		Comment string `json:"comment"`
	} `json:"private"`
}

var (
	subscript = []interface{}{
		bson.M{"_id": "grpL6GdlEDQ3DM:HhlZaX2A80Y", // подставляется сгенеренный _id (topic) и _id из users
			"createdat": bson.M{
				"$date": "2024-07-11T05:57:51.202Z",
			},
			"updatedat": bson.M{
				"$date": "2024-07-11T05:57:51.209Z",
			},
			"user":      "HhlZaX2A80Y",    // подставляется _id из users
			"topic":     "grpL6GdlEDQ3DM", //подставляется сгенеренный _id (topic)
			"delid":     0,
			"recvseqid": 0,
			"readseqid": 0,
			"modewant": bson.M{
				"$numberLong": "255",
			},
			"modegiven": bson.M{
				"$numberLong": "255",
			},
			"private": bson.M{
				"comment": "тест топика-темы",
			},
		},
	}
)
