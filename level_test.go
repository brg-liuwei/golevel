package golevel_test

import (
	"os"
	"testing"

	level "github.com/brg-liuwei/golevel"
)

func TestPutGet(t *testing.T) {
	level.Init(1)
	defer level.Cleanup()

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
