package main

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

func main1() {
	fmt.Println(reflect.TypeOf(1).String())
	fmt.Println(reflect.TypeOf("1").String())
}
func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(2)
		go func(x int) {
			defer wg.Done()
			SetM(x)
		}(i)
		go func(s string) {
			defer wg.Done()
			SetM(s)
		}(strconv.Itoa(i) + "a")
	}
	wg.Wait()
}

var clientMu sync.Mutex                                 //客户取锁的时候需要加同步锁
var muM = map[string]*sync.Mutex{"test": &sync.Mutex{}} //配置客户锁

func getClientMu(client string) *sync.Mutex {
	clientMu.Lock()
	defer clientMu.Unlock()
	_, ok := muM[client]
	if !ok {
		muM[client] = &sync.Mutex{}
	}
	return muM[client]
}

func getClient(x interface{}) string {
	ts := reflect.TypeOf(x).String()
	if ts == "int" {
		return "i"
	}
	return "s"
}

var m = map[string]interface{}{}
var m2 = map[string]interface{}{}

func SetM(x interface{}) {
	client := getClient(x)
	mu := getClientMu(client)
	mu.Lock()
	defer mu.Unlock()
	if client == "i" {
		m[client] = x
		fmt.Println(m[client])
	} else if client == "s" {
		m2[client] = x
		fmt.Println(m2[client])
	} else {
		panic("hahahahahah")
	}
}
