package main

import (
	"context"
	"math/rand"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestInserter(t *testing.T) {

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
	}{
		{"gid1", "gid2", "id1", "id2",
			[]string{"newtopic"},
			[]string{"3e266244-0e23-4f2e-8cb5-b4d118054222"}, []string{"3e266244-0e23-4f2e-8cb5-b4d118054222"}},
	}

	var prev_gid1 string
	for _, test := range gtests {
		if test.gid1 != prev_gid1 {
			t.Logf("%s", test.gid1)
			prev_gid1 = test.gid1
		}
	}

	var prev_gid2 string
	for _, test := range gtests {
		if test.gid2 != prev_gid2 {
			t.Logf("%s", test.gid2)
			prev_gid2 = test.gid2
		}
	}

	var prev_id1 string
	for _, test := range gtests {
		if test.id1 != prev_id1 {
			t.Logf("%s", test.id1)
			prev_id1 = test.id1
		}
	}

	var prev_id2 string
	for _, test := range gtests {
		if test.id2 != prev_id2 {
			t.Logf("%v", test.id2)
			prev_id2 = test.id2
		}
	}

	var prev_intop2 []string
	for _, test := range gtests {
		if test.intop2 != nil && prev_intop2 == nil {
			t.Logf("%s", test.intop2)
			prev_intop2 = test.intop2
		}
	}

	var prev_insub1 []string
	for _, test := range gtests {
		if test.insub1 != nil && prev_insub1 == nil {
			t.Logf("%s", test.insub1)
			prev_insub1 = test.insub1
		}
	}

	var prev_insub2 []string
	for _, test := range gtests {
		if test.insub2 != nil && prev_insub2 == nil {
			t.Logf("%s\n", test.insub2)
			prev_insub2 = test.insub2
		}
	}

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

	var ktest = []struct {
		key int64
	}{
		{111},
		{222222},
		{333333333333},
		{5555555555555555},
	}

	var prev_key int64
	for _, test := range ktest {
		if test.key != prev_key {
			t.Logf("%v", test.key)
			prev_key = test.key
		}

		var pkey int64
		key := time.Now().UTC().UnixNano()
		if !reflect.DeepEqual(key, pkey) {
			t.Logf("%v", key)
			pkey = key
		}

		var prng *rand.Rand
		rng := rand.New(rand.NewSource(pkey))
		if !reflect.DeepEqual(rng, prng) {
			t.Logf("%v", rng)
			prng = rng
		}

		prev_gid1 = genTopicName(prng)
		if prev_gid1 == "" {
			t.Log("Error of generate string")
		}
		t.Log("Random string: ", prev_gid1)
	}

	cn := client.Database("gotest").Collection("topics")
	var pcn *mongo.Collection
	if cn != pcn {
		t.Logf("%v", cn)
		pcn = cn
	}
	if !reflect.DeepEqual(pcn, cn) {
		t.Errorf("Check client.Database() cn = nil %v, want = dtopic %v", pcn, dtopic)
	}

	var prev_opts *options.InsertManyOptions
	opts := options.InsertMany().SetOrdered(false)
	if !reflect.DeepEqual(opts, prev_opts) {
		t.Logf("%v", opts)
		prev_opts = opts
	}

	var prev_res *mongo.InsertManyResult
	res, _ := cn.InsertMany(context.TODO(), dtopic, prev_opts)
	if !reflect.DeepEqual(res, prev_res) {
		t.Logf("%v", res)
		prev_res = res
	}
}
