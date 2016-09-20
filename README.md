# golevel
=======

## golang lib of leveldb  

### *Install:*
    go get github.com/brg-liuwei/golevel  

Before Using this lib, you need to do as follows:  

### *Linux platform:* 
    cp -r /your/gopath/src/github/brg-liuwei/golevel/lib/leveldb/include/leveldb /usr/include  
    cp lib/leveldb/libleveldb.so /usr/lib  

### *MacOS platform:*
    cp -r /your/gopath/src/github/brg-liuwei/golevel/lib/leveldb/include/leveldb /usr/include  
    cp lib/leveldb/libleveldb.dylib /usr/lib  

### *Windows platform:*  
    Read my code and implements a windows leveldb go lib for yourself :)


### Example

    package main
    
    import (
    	"fmt"
    	"runtime"
    	"sync"
    	"time"

    	level "github.com/brg-liuwei/golevel"
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

