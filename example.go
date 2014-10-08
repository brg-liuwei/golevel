package main

import (
	"fmt"
	level "github.com/brg-liuwei/golevel"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	level.Init(2)
	table := "db"
	level.Open(&table)

	key := "Key"
	val := "Value"
	var wg sync.WaitGroup
	for i := 0; i != 4; i++ {
		wg.Add(1)
		go func(j int) {
			for k := 10000 * j; k < 10000*(j+1); k++ {
				key := fmt.Sprint(key, k)
				val := fmt.Sprint(val, k)
				level.BatchWrite(&table, &key, &val)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	k := fmt.Sprint(key, 100)
	v, _ := level.Get(&table, &k)
	fmt.Println(k, ":", v)

	k = fmt.Sprint(key, 30100)
	v, _ = level.Get(&table, &k)
	fmt.Println(k, ":", v)

	level.Close(&table)
	level.Cleanup()
}
