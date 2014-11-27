package golevel

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include "leveldb/c.h"

#cgo CFLAGS: -I lib/leveldb/include
#cgo LDFLAGS: -L lib/leveldb -lleveldb -lpthread


static leveldb_t **db;
static leveldb_writebatch_t **batch;
static size_t *batch_cnt;

static char **errmsg;
static leveldb_options_t *opt;
static leveldb_readoptions_t *ropt;
static leveldb_writeoptions_t *wopt;

static void initDbEnv(size_t max_tables)
{
    opt = leveldb_options_create();
    leveldb_options_set_create_if_missing(opt, (unsigned char)1);
    ropt = leveldb_readoptions_create();
    wopt = leveldb_writeoptions_create();

    db = (leveldb_t **)calloc(max_tables, sizeof(leveldb_t *));
    if (db == NULL) {
        printf("init db env error(calloc db), set tables size: %lu\n", max_tables);
        abort();
    }

    errmsg = (char **)calloc(max_tables, sizeof(char *));
    if (errmsg == NULL) {
        printf("init db env error(calloc errmsg), set tables size: %lu\n", max_tables);
        abort();
    }

    batch = (leveldb_writebatch_t **)calloc(max_tables, sizeof(leveldb_writebatch_t *));
    if (batch == NULL) {
        printf("init db env error(calloc writebatch), set tables size: %lu\n", max_tables);
        abort();
    }

    batch_cnt = (size_t *)calloc(max_tables, sizeof(size_t));
    if (batch_cnt == NULL) {
        printf("init db env error(calloc batch_cnt), set tables size: %lu\n", max_tables);
        abort();
    }
}

static void deleteDbEnv()
{
    if (opt) {
        leveldb_options_destroy(opt);
        opt = NULL;
    }
    if (ropt) {
        leveldb_readoptions_destroy(ropt);
        ropt = NULL;
    }
    if (wopt) {
        leveldb_writeoptions_destroy(wopt);
        wopt = NULL;
    }
    if (db) {
        free(db);
        db = NULL;
    }
    if (errmsg) {
        free(errmsg);
        errmsg = NULL;
    }
    if (batch) {
        free(batch);
        batch = NULL;
    }
    if (batch_cnt) {
        free(batch_cnt);
        batch_cnt = NULL;
    }
}

static void batch_flush(int idx, void *batch)
{
    if (batch) {
        leveldb_write(db[idx], wopt, (leveldb_writebatch_t *)batch, &errmsg[idx]);
        leveldb_writebatch_clear((leveldb_writebatch_t *)batch);
        leveldb_writebatch_destroy((leveldb_writebatch_t *)batch);
    }
}

static void *batch_get(int idx)
{
    leveldb_writebatch_t *old_batch;
    if (batch_cnt[idx] == 0) {
        return NULL;
    }
    old_batch = batch[idx];
    batch[idx] = NULL;
    batch_cnt[idx] = 0;
    return (void *)old_batch;
}

static void *batch_write(const char *k, size_t k_len, const char *v, size_t v_len, int idx)
{
    if (batch[idx] == NULL) {
        batch[idx] = leveldb_writebatch_create();
        if (batch[idx] == NULL) {
            abort();
        }
        batch_cnt[idx] = 0;
    }
    leveldb_writebatch_put(batch[idx], k, k_len, v, v_len);
    batch_cnt[idx]++;
    if (batch_cnt[idx] >= 10000) {
        return batch_get(idx);
    }
    return NULL;
}

static void *batch_delete(const char *k, size_t k_len, int idx)
{
    void *ret = NULL;
    if (batch[idx] == NULL) {
        batch[idx] = leveldb_writebatch_create();
        if (batch[idx] == NULL) {
            abort();
        }
        batch_cnt[idx] = 0;
    }
    leveldb_writebatch_delete(batch[idx], k, k_len);
    batch_cnt[idx]++;
    if (batch_cnt[idx] >= 10000) {
        ret = (void *)batch_cnt[idx];
        batch[idx] = NULL;
        batch_cnt[idx] = 0;
    }
    return ret;
}

static int put(const char *k, size_t k_len, const char *v, size_t v_len, int idx)
{
    leveldb_put(db[idx], wopt, k, k_len, v, v_len, &errmsg[idx]);

    if (errmsg[idx] != NULL) {
        printf("put k, v failed: %s\n", errmsg[idx]);
        leveldb_free(errmsg[idx]);
        errmsg[idx] = NULL;
        return -1;
    }
    return 0;
}

static char *get(const char *k, size_t k_len, size_t *v_len, int idx)
{
    char  *val;

    val = leveldb_get(db[idx], ropt, k, k_len, v_len, &errmsg[idx]);

    if (errmsg[idx] != NULL) {
        printf("get k failed: %s\n", errmsg[idx]);
        leveldb_free(errmsg[idx]);
        errmsg[idx] = NULL;
    }
    return val;
}

static int del(const char *k, size_t k_len, int idx)
{
    leveldb_delete(db[idx], wopt, k, k_len, &errmsg[idx]);

    if (errmsg[idx] != NULL) {
        printf("delete k failed: %s\n", errmsg[idx]);
        leveldb_free(errmsg[idx]);
        errmsg[idx] = NULL;
        return -1;
    }
    return 0;
}

static void openDb(const char *name, int idx)
{
    db[idx] = leveldb_open(opt, name, &errmsg[idx]);

    if (errmsg[idx] != NULL) {
        printf("open levelDb %s error: %s\n", name, errmsg[idx]);
        leveldb_free(errmsg[idx]);
        errmsg[idx] = NULL;
        abort();
    }
}

static void closeDb(int idx)
{
    if (batch[idx] != NULL && batch_cnt[idx] != 0) {
        batch_flush(idx, batch[idx]);
    }

    if (errmsg[idx] != NULL) {
        leveldb_free(errmsg[idx]);
        errmsg[idx] = NULL;
    }
    if (opt != NULL) {
        leveldb_options_destroy(opt);
        opt = NULL;
    }
    if (ropt != NULL) {
        leveldb_readoptions_destroy(ropt);
        ropt = NULL;
    }
    if (wopt != NULL) {
        leveldb_writeoptions_destroy(wopt);
        wopt = NULL;
    }
    if (db[idx] != NULL) {
        leveldb_close(db[idx]);
        db[idx] = NULL;
    }
}

static void *newIterator(int idx)
{
    return (void *)leveldb_create_iterator(db[idx], ropt);
}

static void iterDestroy(void *itr)
{
    leveldb_iter_destroy((leveldb_iterator_t *)itr);
}

static void iterSeek(void *itr, const char *k, size_t k_len)
{
    leveldb_iter_seek((leveldb_iterator_t *)itr, k, k_len);
}

static void iterSeekToFirst(void *itr)
{
    leveldb_iter_seek_to_first((leveldb_iterator_t *)itr);
}

static void iterSeekToLast(void *itr)
{
    leveldb_iter_seek_to_last((leveldb_iterator_t *)itr);
}

static void iterPrev(void *itr)
{
    leveldb_iter_prev((leveldb_iterator_t *)itr);
}

static void iterNext(void *itr)
{
    leveldb_iter_next((leveldb_iterator_t *)itr);
}

static unsigned char iterValid(void *itr)
{
    return leveldb_iter_valid((leveldb_iterator_t *)itr);
}

static const char *iterKey(void *itr, size_t *k_len)
{
    return leveldb_iter_key((leveldb_iterator_t *)itr, k_len);
}

static const char *iterValue(void *itr, size_t *v_len)
{
    return leveldb_iter_value((leveldb_iterator_t *)itr, v_len);
}
*/
import "C"

import (
	"errors"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/brg-liuwei/spinlock"
)

var TableNotOpened error = errors.New("Table has not been opened")
var PutError error = errors.New("levelDb put failed")
var BatchError error = errors.New("levelDb batch write failed")
var NotExist error = errors.New("key not exist")
var DelError error = errors.New("levelDb del failed")

var maxTables int
var curTables int
var tableIdx map[string]int

var once sync.Once

type batchInfo struct {
	idx   int
	batch unsafe.Pointer
}

var batchChan chan batchInfo

func Init(tables int) {
	once.Do(func() {
		maxTables = tables
		curTables = 0
		tableIdx = make(map[string]int)
		batchChan = make(chan batchInfo, 1024)
		C.initDbEnv(C.size_t(tables))
		go batchLoop()
	})
}

func Cleanup() {
	C.deleteDbEnv()
}

var tableLock sync.Mutex

func Open(name *string) {
	tableLock.Lock()
	defer tableLock.Unlock()
	_open(name)
}

func _open(name *string) {
	if _, ok := tableIdx[*name]; ok {
		return
	}
	if curTables >= maxTables {
		panic("Cannot open too many tables: " + strconv.Itoa(curTables))
	}
	tableIdx[*name] = curTables
	idx := curTables
	curTables++
	cname := C.CString(*name)
	defer C.free(unsafe.Pointer(cname))
	C.openDb(cname, C.int(idx))
}

func Close(name *string) {
	tableLock.Lock()
	defer tableLock.Unlock()
	_close(name)
}

func _close(name *string) {
	if v, ok := tableIdx[*name]; ok {
		C.closeDb(C.int(v))
		curTables--
		delete(tableIdx, *name)
	}
}

var spin spinlock.SpinLock

func BatchWrite(table *string, k *string, v *string) error {
	idx, ok := tableIdx[*table]
	if !ok {
		return TableNotOpened
	}
	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))
	val := C.CString(*v)
	defer C.free(unsafe.Pointer(val))

	spin.Lock()
	batch := C.batch_write(key, C.size_t(len(*k)), val, C.size_t(len(*v)), C.int(idx))
	spin.UnLock()

	if batch != unsafe.Pointer(nil) {
		var bi batchInfo
		bi.idx = idx
		bi.batch = batch
		batchChan <- bi
	}
	return nil
}

func BatchDelete(table *string, k *string) error {
	idx, ok := tableIdx[*table]
	if !ok {
		return TableNotOpened
	}
	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))

	spin.Lock()
	batch := C.batch_delete(key, C.size_t(len(*k)), C.int(idx))
	spin.UnLock()

	if batch != unsafe.Pointer(nil) {
		var bi batchInfo
		bi.idx = idx
		bi.batch = batch
		batchChan <- bi
	}
	return nil
}

var batchLock sync.Mutex

func batchSync() {
	for i := 0; i < maxTables; i++ {
		spin.Lock()
		batch := C.batch_get(C.int(i))
		spin.UnLock()
		if batch != unsafe.Pointer(nil) {
			batchLock.Lock()
			C.batch_flush(C.int(i), unsafe.Pointer(batch))
			batchLock.Unlock()
		}
	}
}

func batchLoop() {
	for {
		select {
		case bi := <-batchChan:
			batchLock.Lock()
			C.batch_flush(C.int(bi.idx), unsafe.Pointer(bi.batch))
			batchLock.Unlock()
			_ = time.Second
			//case <-time.After(time.Second):
			//	batchSync()
		}
	}
}

func Put(table *string, k *string, v *string) error {
	idx, ok := tableIdx[*table]
	if !ok {
		return TableNotOpened
	}
	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))
	val := C.CString(*v)
	defer C.free(unsafe.Pointer(val))

	if ok := C.put(key, C.size_t(len(*k)), val, C.size_t(len(*v)), C.int(idx)); ok != 0 {
		return PutError
	}
	return nil
}

func Get(table *string, k *string) (v string, err error) {
	idx, ok := tableIdx[*table]
	if !ok {
		err = TableNotOpened
		return
	}

	var v_len C.size_t

	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))

	val := C.get(key, C.size_t(len(*k)), &v_len, C.int(idx))
	if val == nil {
		err = NotExist
		return
	}
	defer C.leveldb_free(unsafe.Pointer(val))

	v = C.GoStringN(val, C.int(v_len))
	return
}

func Delete(table *string, k *string) error {
	idx, ok := tableIdx[*table]
	if !ok {
		return TableNotOpened
	}

	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))

	if ok := C.del(key, C.size_t(len(*k)), C.int(idx)); ok != 0 {
		return DelError
	}
	return nil
}

type Iter struct {
	valid bool
	itr   unsafe.Pointer
}

func NewIterator(table *string) (*Iter, error) {
	if idx, ok := tableIdx[*table]; !ok {
		return nil, TableNotOpened
	} else {
		return &Iter{false, unsafe.Pointer(C.newIterator(C.int(idx)))}, nil
	}
	return nil, nil
}

func (this *Iter) Destroy() {
	if this.valid {
		C.iterDestroy(this.itr)
		this.itr = nil
		this.valid = false
	}
}

func (this *Iter) Seek(k *string) {
	if this.itr != nil {
		this.valid = true
		key := C.CString(*k)
		defer C.free(unsafe.Pointer(key))

		C.iterSeek(this.itr, key, C.size_t(len(*k)))
	}
}

func (this *Iter) SeekToFirst() {
	if this.itr != nil {
		this.valid = true
		C.iterSeekToFirst(this.itr)
	}
}

func (this *Iter) SeekToLast() {
	if this.itr != nil {
		this.valid = true
		C.iterSeekToLast(this.itr)
	}
}

func (this *Iter) Prev() {
	if this.valid {
		C.iterPrev(this.itr)
	}
}

func (this *Iter) Next() {
	if this.valid {
		C.iterNext(this.itr)
	}
}

func (this *Iter) Valid() bool {
	if this.valid {
		if C.iterValid(this.itr) != C.uchar(0) {
			return true
		}
		this.valid = false
	}
	return false
}

func (this *Iter) Key() string {
	if !this.valid {
		return ""
	}

	var k_len C.size_t
	key := C.iterKey(this.itr, &k_len)
	if key == nil {
		return ""
	}
	return C.GoStringN(key, C.int(k_len))
}

func (this *Iter) Value() string {
	if !this.valid {
		return ""
	}
	var v_len C.size_t
	val := C.iterValue(this.itr, &v_len)
	if val == nil {
		return ""
	}
	return C.GoStringN(val, C.int(v_len))
}
