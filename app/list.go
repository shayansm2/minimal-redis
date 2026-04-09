package main

import (
	"errors"
	"slices"
	"strconv"
	"sync"
	"time"
)

var subscribeChannels map[string]*LinkedList[chan<- string]
var mu sync.Mutex

func init() {
	subscribeChannels = make(map[string]*LinkedList[chan<- string])
	bgJobs = append(bgJobs, PushedElementsPublisherJob)
}

func lPushHandler(args []string) any {
	if len(args) < 2 {
		return errors.New("ERR not enough args provided")
	}
	name := args[0]
	elements := make([]string, len(args)-1)
	copy(elements, args[1:])
	slices.Reverse(elements)

	list, _ := getListFromDB(name)
	list = append(elements, list...)

	db.set(name, list)

	return len(list)
}

func rPushHandler(args []string) any {
	if len(args) < 2 {
		return errors.New("ERR not enough args provided")
	}
	name := args[0]

	list, _ := getListFromDB(name)
	list = append(list, args[1:]...)

	db.set(name, list)

	return len(list)
}

func lRangeHandler(args []string) any {
	if len(args) < 3 {
		return errors.New("ERR not enough args provided")
	}
	list, found := getListFromDB(args[0])
	if !found {
		return []BulkStr{}
	}

	start, _ := strconv.Atoi(args[1])
	stop, _ := strconv.Atoi(args[2])

	if start < 0 {
		start = max(start+len(list), 0)
	}
	if stop < 0 {
		stop = max(stop+len(list), 0)
	}
	stop++ // make exclusive
	if start >= len(list) || start >= stop {
		return []BulkStr{}
	}
	stop = min(stop, len(list))
	result := make([]BulkStr, stop-start)
	for i, val := range list[start:stop] {
		result[i] = BulkStr(val)
	}
	return result
}

func lLenHandler(args []string) any {
	if len(args) < 1 {
		return errors.New("ERR not enough args provided")
	}
	list, found := getListFromDB(args[0])
	if !found {
		return 0
	}
	return len(list)
}

func lPopHandler(args []string) any {
	name := args[0]
	count := 1
	if len(args) > 1 {
		count, _ = strconv.Atoi(args[1])
	}
	list, found := getListFromDB(args[0])
	if !found || len(list) == 0 {
		return nil
	}
	pops := list[:count]
	list = list[count:]
	db.set(name, list)

	bulks := make([]BulkStr, len(pops))
	for i, pop := range pops {
		bulks[i] = BulkStr(pop)
	}
	if count == 1 {
		return bulks[0]
	}
	return bulks
}

func bLPopHandler(args []string) any {
	if len(args) < 2 {
		return errors.New("ERR not enough args provided")
	}
	name := args[0]
	timeout, _ := strconv.ParseFloat(args[1], 64)
	ch := make(chan string, 0)
	id := subscribeToListPush(name, ch)
	if timeout == 0 {
		pop := <-ch
		return []BulkStr{BulkStr(name), BulkStr(pop)}
	}

	select {
	case pop := <-ch:
		return []BulkStr{BulkStr(name), BulkStr(pop)}
	case <-time.After(time.Millisecond * time.Duration(1000*timeout)):
		unsubscribeToListPush(name, id)
		close(ch)
		var nullArr []string
		return nullArr
	}
}

func getListFromDB(name string) (list []string, found bool) {
	l, found := db.get(name)
	if !found {
		return
	}
	list = l.([]string)
	return
}

func subscribeToListPush(name string, ch chan<- string) int {
	mu.Lock()
	defer mu.Unlock()
	if _, found := subscribeChannels[name]; !found {
		subscribeChannels[name] = NewLinkedList[chan<- string]()
	}
	return subscribeChannels[name].push(ch)
}

func unsubscribeToListPush(name string, id int) {
	mu.Lock()
	defer mu.Unlock()
	subscribeChannels[name].del(id)
}

func PushedElementsPublisherJob() {
	for {
		mu.Lock()
		for name, chans := range subscribeChannels {
			if chans.isEmpty() {
				continue
			}
			list, found := getListFromDB(name)
			if !found || len(list) == 0 {
				continue
			}
			pop := list[0]
			list = list[1:]
			ch := chans.pop()
			ch <- pop
			close(ch)
			db.set(name, list)
			subscribeChannels[name] = chans
		}
		mu.Unlock()
	}
}
