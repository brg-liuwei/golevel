package golevel_test

import (
	"os"
	"testing"
	"time"

	level "github.com/brg-liuwei/golevel"
)

func init() {
	level.Init(1)
}

func TestPutGet(t *testing.T) {
	table := "test.db"
	level.Open(&table)
	defer func() {
		if err := os.RemoveAll(table); err != nil {
			t.Error("remove table(", table, ") error: ", err)
		}
	}()
	defer level.Close(&table)

	key := "Key"
	val := "Value"

	if err := level.Put(&table, &key, &val); err != nil {
		t.Error("level.Put error: ", err)
	}

	if v, err := level.Get(&table, &key); err != nil {
		t.Error("level.Get error: ", err)
	} else if v != val {
		t.Error("level.Get(", key, ") = ", v, ", ", val, " expected")
	}
}

func TestIterator(t *testing.T) {
	table := "test.db"
	level.Open(&table)
	defer func() {
		if err := os.RemoveAll(table); err != nil {
			t.Error("remove table(", table, ") error: ", err)
		}
	}()
	defer level.Close(&table)

	type kvSt struct{ k, v string }
	kvs := []kvSt{
		{"k1", "v1"},
		{"k2", "v2"},
		{"k3", "v3"},
		{"k4", "v4"},
		{"k5", "v5"},
	}

	for _, kv := range kvs {
		level.BatchWrite(&table, &kv.k, &kv.v)
	}

	// wait for flush
	time.Sleep(2 * time.Second)

	itr, err := level.NewIterator(&table)
	if err != nil {
		t.Error("new iterator error: ", err)
	}
	defer itr.Destroy()

	i := 0
	for itr.SeekToFirst(); itr.Valid(); itr.Next() {
		if itr.Key() != kvs[i].k {
			t.Error("itr.Key(): ", itr.Key(), "; kvs[", i, "].k: ", kvs[i].k)
		}
		if itr.Value() != kvs[i].v {
			t.Error("itr.Value(): ", itr.Value(), "; kvs[", i, "].k: ", kvs[i].v)
		}
		i++
	}

	if i != len(kvs) {
		t.Error("i = ", i, "; expect: ", len(kvs))
	}
}

func TestEnd(t *testing.T) {
	level.Cleanup()
}
