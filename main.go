package level

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init(2)
	table := "db"
	Open(&table)

	key := "Key"
	val := "Value"
	var wg sync.WaitGroup
	for i := 0; i != 10001; i++ {
		wg.Add(1)
		go func(i int) {
			k := fmt.Sprint(key, i)
			v := fmt.Sprint(val, i)
			BatchWrite(&table, &k, &v)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
	k := fmt.Sprint(key, 100)
	v, _ := Get(&table, &k)
	fmt.Println(k, ":", v)

	Close(&table)
	Cleanup()
}
