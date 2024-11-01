package main

import (
	"context"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInserter(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if ctx == nil && cancel != nil {
		t.Errorf("Check func WithTimeout() ctx = nil %v, want cancel = nil %v", ctx, cancel)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsnMongo))
	if client == nil && err != nil {
		t.Errorf("Check func mongo.Connect() client = nil %v, want err = nil %v", client, err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			switch p := recover(); p {
			case err != nil:
				panic(err)
				panic(p)
			default:
				panic(p)
			}
		} else {
			t.Log("Disconnect not errors")
		}
	}()

	saved := inserter
	defer func() { inserter = saved }()

	var gtests = []struct {
		gid1   string
		gid2   string
		id1    string `json:"user_uuid"`
		id2    string `json:"user_uuid"`
		intop2 []string
		insub1 []string
		insub2 []string
		key    int64
		want   bool
	}{
		{"gid1", "gid2", "id1", "id2",
			[]string{"newtopic"},
			[]string{"3e266244-0e23-4f2e-8cb5-b4d118054222"}, []string{"3e266244-0e23-4f2e-8cb5-b4d118054222"}, 8765432123456, true},
	}

	var prev_gid1 string
	for _, test := range gtests {
		if test.gid1 != prev_gid1 {
			t.Logf("%s", test.gid1)
			prev_gid1 = test.gid1
		}

		var prev_gid2 string
		if test.gid2 != prev_gid2 {
			t.Logf("%s", test.gid2)
			prev_gid2 = test.gid2
		}

		var prev_id1 string
		if test.id1 != prev_id1 {
			t.Logf("%s", test.id1)
			prev_id1 = test.id1
		}

		var prev_id2 string
		if test.id2 != prev_id2 {
			t.Logf("%v", test.id2)
			prev_id2 = test.id2
		}

		var prev_intop2 []string
		if test.intop2 != nil && prev_intop2 == nil {
			t.Logf("%s", test.intop2)
			prev_intop2 = test.intop2
		}

		var prev_insub1 []string
		if test.insub1 != nil && prev_insub1 == nil {
			t.Logf("%s", test.insub1)
			prev_insub1 = test.insub1
		}

		var prev_insub2 []string
		if test.insub2 != nil && prev_insub2 == nil {
			t.Logf("%s\n", test.insub2)
			prev_insub2 = test.insub2
		}

		var prev_key int64
		if test.key != prev_key {
			t.Logf("%v", test.key)
			prev_key = test.key
		}

		var pkey int64
		if !reflect.DeepEqual(test.key, pkey) {
			t.Logf("%v", test.key)
			pkey = test.key
		}
		keygot := time.Now().UTC().UnixNano()
		if keygot == 0 {
			t.Errorf("Check UnixNano:(%v) = %v", pkey, test.want)
		}

		var prevKey int64
		if test.key != prevKey {
			t.Logf("%v", test.key)
			test.key = prevKey
		}
		rngot := rand.New(rand.NewSource(keygot))
		if rngot == nil {
			t.Errorf("Check UnixNano:(%v) = %v", pkey, test.want)
		}

		var prng *rand.Rand
		if !reflect.DeepEqual(rngot, prng) {
			t.Logf("%v", rngot)
			prng = rngot
		}

		if got := genTopicName(prng); got == "" {
			t.Errorf("Check genTopicName:(%v) = %v", prng, test.want)
		}

		cn := client.Database("gotest").Collection("topics")
		if cn == nil {
			t.Errorf("Check client.Database:(%v) = %v", cn, test.want)
		}
		var pcn *mongo.Collection
		if cn != pcn {
			t.Logf("%v", cn)
			pcn = cn
		}

		var prevOpts *options.InsertManyOptions
		opts := options.InsertMany().SetOrdered(false)
		if opts == nil {
			t.Errorf("Check InsertMany:(%v) = %v", opts, test.want)
		}
		if !reflect.DeepEqual(opts, prevOpts) {
			t.Logf("%v", opts)
			prevOpts = opts
		}

		var prevResult *mongo.InsertManyResult
		result, _ := cn.InsertMany(context.TODO(), dtopic, prevOpts)
		if result == nil {
			t.Errorf("Check InsertManyRun:(%v) = %v", result, test.want)
		}
		if !reflect.DeepEqual(result, prevResult) {
			t.Logf("%v", result)
			prevResult = result
		}

		// Срез строковых значений с размером результата. Slice len of cursor
		testop := make([]string, 0, len(prevResult.InsertedIDs))

		// Формирование строкового слайса. Gets []string slice
		for _, v := range prevResult.InsertedIDs {
			if v != nil {
				testop = append(testop, v.(string))
			}
			t.Logf("%v", testop)
		}

	}
}
