package main

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFinder(t *testing.T) {
	saved := idFinder
	defer func() { idFinder = saved }()

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

	var tests = []struct {
		uuid   string
		filter interface{}
		opts   interface{}
		values interface{}
		want   bool
	}{
		{"3e266244-0e23-4f2e-8cb5-b4d118054222", nil, nil, nil, true},
		{"3e266244-0e23-4f2e-8cb5-b4d118054111", nil, nil, nil, true},
		{"3............00,,,,,,;;;;;;;;;;;;;00", nil, nil, nil, true},
		{"    ", nil, nil, nil, true},
	}

	var prevUuid string
	for _, test := range tests {
		if test.uuid != prevUuid {
			t.Logf("%v", test.uuid)
			prevUuid = test.uuid
		}

		var prevFilter interface{}
		if test.filter != nil {
			t.Logf("%v", test.filter)
			prevFilter = test.filter
		}

		var prevOpts interface{}
		if test.opts != nil {
			t.Logf("%v", test.opts)
			prevOpts = test.opts
			t.Logf("%v", prevOpts)
		}

		var prevValues interface{}
		if test.values == nil {
			t.Logf("%v", test.values)
			prevValues = test.values
			t.Logf("%v", prevValues)
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

		filter := bson.D{{"tags", bson.D{{"$eq", prevUuid}}}}
		if filter == nil {
			t.Logf("Check filter: %v", filter)
		}

		opts := options.Distinct().SetMaxTime(2 * time.Second)
		if opts == nil {
			t.Logf("Check opts: %v", opts)
		}

		values, err := cn.Distinct(context.TODO(), "_id", prevFilter, opts)
		if err != nil {
			t.Logf("Error check Distinct: %v", err)
		} else {
			t.Logf("%v", values)
		}
	}
}
