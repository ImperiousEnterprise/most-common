package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestFindMostCommon(t *testing.T) {
	dict := make(map[string]int64)
	channel := make(chan (string))
	done := make(chan (bool), 1)
	go func() {
		for {
			j, more := <-channel
			if more {
				dict[j]++
			} else {
				close(channel)
				done <- true
				return
			}
		}
	}()
	s := strings.NewReader("This is a test string\nThis is a test.\n")
	if err := findMostCommon(s, 0, int64(s.Len()), channel); err != nil {
		t.Error(err)
	}
	close(done)

	ans := generatePriorityQueue(dict)
	for i, c := range ans {
		fmt.Printf("%d.) %s - %d \n", i+1, c, dict[c])
	}
}

func TestFindMostCommonFailure(t *testing.T) {
	dict := make(map[string]int64)
	channel := make(chan (string))
	done := make(chan (bool), 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for s := range channel {
			dict[s]++
		}
		done <- true
		wg.Done()
	}()
	if err := findMostCommon(nil, 0, int64(0), channel); err == nil {
		t.FailNow()
	}
	close(done)
}

func TestParseData(t *testing.T) {
	const mb = 1024 * 1024
	channel := make(chan string)
	wg := sync.WaitGroup{}
	info, err := os.Stat("test.txt")
	if err != nil {
		t.FailNow()
	}
	container := FileContainer{index: 0 + 1, filename: "test.txt", size: info.Size()}
	parseData(&container, 1*mb, &wg, channel)
}
func TestPriorityQueue(t *testing.T) {
	dict := make(map[string]int64)
	dict["This is a"] = 3
	dict["This is test"] = 2
	dict["This is apple"] = 1
	dict["This is astroid"] = 4

	ans := generatePriorityQueue(dict)
	if ans[0] != "This is astroid" {
		t.FailNow()
	}
}

func TestPriorityQueueFail(t *testing.T) {
	dict := make(map[string]int64)
	dict["This is a"] = 3
	dict["This is test"] = 2
	dict["This is apple"] = 1
	dict["This is astroid"] = 1

	ans := generatePriorityQueue(dict)
	if ans[0] == "This is astroid" {
		t.FailNow()
	}
}
