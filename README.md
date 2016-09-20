# golevel

[![Build Status](https://travis-ci.org/brg-liuwei/golevel.svg?branch=master)](https://travis-ci.org/brg-liuwei/golevel)

golang lib of leveldb  

### Install leveldb into path `/usr/local/leveldb`:

    git clone https://github.com/google/leveldb.git /usr/local/leveldb
    git checkout -b v1.8 v1.8
    make

### Install:

    export PKG_CONFIG_PATH=${PKG_CONFIG_PATH}:${GOPATH}/src/github.com/brg-liuwei/golevel
    go get -u github.com/brg-liuwei/golevel

### Example

    package main
    
    import (
        "log"

    	level "github.com/brg-liuwei/golevel"
    )
    
    func main() {
	    level.Init(1) // how many tables prepareing to open
	    defer level.Cleanup()

	    table := "./test.db"

	    level.Open(&table)
	    defer level.Close(&table)

	    key := "Key"
	    val := "Value"

	    if err := level.Put(&table, &key, &val); err != nil {
            log.Fatal("level.Put err: ", err)
	    }

	    if v, err := level.Get(&table, &key); err != nil {
	    	log.Fatal("level.Get error: ", err)
	    } else if v != val {
	    	log.Fatal("level.Get(", key, ") = ", v, ", ", val, " expected")
	    }
    }

