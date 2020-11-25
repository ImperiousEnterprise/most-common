package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	q "new-relic/lib/queue"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type FileContainer struct {
	index    int
	filename string
	size     int64
}

const mb = 1024 * 1024
const gb = 1024 * mb

func main() {
	start := time.Now()
	// A waitgroup to wait for all go-routines to finish.
	wg := sync.WaitGroup{}

	// This channel is used to send every read word in various go-routines.
	channel := make(chan string)

	// A dictionary which stores the count of unique words.
	dict := make(map[string]int64)

	// Done is a channel to signal the main thread that all the words have been
	// entered in the dictionary.
	done := make(chan bool, 1)

	// Read all incoming words from the channel and add them to the dictionary.
	go func() {
		for s := range channel {
			dict[s]++
		}

		// Signal the main thread that all the words have entered the dictionary.
		done <- true
	}()

	stat, _ := os.Stdin.Stat()
	//Detect if text is coming from piping
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := findMostCommon(os.Stdin, 0, 4*mb, channel)
		if err != nil {
			panic(err)
		}
	} else {
		argsWithoutProg := os.Args[1:]
		for index, x := range argsWithoutProg {
			info, err := os.Stat(x)
			if os.IsNotExist(err) {
				fmt.Printf("%s does not exist \n", x)
			} else {
				container := FileContainer{index: index + 1, filename: x, size: info.Size()}
				if info.Size() >= 1*gb {
					parseData(&container, 4*mb, &wg, channel)
				} else if info.Size() >= 100*mb {
					parseData(&container, 4*mb, &wg, channel)
				} else {
					parseData(&container, container.size, &wg, channel)
				}
			}
		}
	}

	// Wait for all go routines to complete.
	wg.Wait()
	close(channel)

	// Wait for dictionary to process all the words.
	<-done
	close(done)
	//Generate list of top 100
	ans := generatePriorityQueue(dict)
	for i, c := range ans {
		fmt.Printf("%d.) %d - %s  \n", i+1, dict[c], c)
	}

	duration := time.Since(start)
	fmt.Println(duration)
}
func parseData(container *FileContainer, chunkSize int64, wg *sync.WaitGroup, channel chan string) {
	// current signifies the chunk size of file to be processed by every thread.
	var current int64
	file, _ := os.Open(container.filename)
	partition := sync.WaitGroup{}
	wg.Add(1)
	for i := 0; current < container.size; i++ {
		partition.Add(1)
		start := current
		go func() {
			read(file, start, chunkSize, channel)
			partition.Done()
		}()
		// increment current to get to the next byte
		current += chunkSize + 1
	}
	// close file after its been scanned
	go func() {
		partition.Wait()
		fmt.Printf("Worker %d finished scanning %s \n", container.index, container.filename)
		file.Close()
		wg.Done()
	}()
}
func generatePriorityQueue(dict map[string]int64) []string {
	an_length := 100
	if len(dict) < 100 {
		an_length = len(dict)
	}

	pq := &q.PQ{}
	heap.Init(pq)
	for w, _ := range dict {
		heap.Push(pq, &q.WordCnt{w, dict[w]})
		if pq.Len() > an_length {
			heap.Pop(pq)
		}
	}

	ans := make([]string, an_length)
	for i := an_length - 1; i >= 0; i-- {
		wc := heap.Pop(pq).(*q.WordCnt)
		ans[i] = wc.Word
	}
	return ans
}

func read(file *os.File, offset int64, limit int64, channel chan (string)) {
	// Move the pointer of the file to the start of designated chunk.
	_, err := file.Seek(offset, 0)
	if err != nil {
		panic(err)
	}
	err = findMostCommon(file, offset, limit, channel)
	if err != nil {
		panic(err)
	}
}

func findMostCommon(file io.Reader, offset int64, limit int64, channel chan string) error {
	if file == nil {
		return errors.New("blank file provided")
	}
	reader := bufio.NewReader(file)

	// This block of code ensures that the start of chunk is a new word. If
	// a character is encountered at the given position it moves a few bytes till
	// the end of the word.
	if offset != 0 {
		_, err := reader.ReadBytes('\n')
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}

	var readSize int64

	for {
		// Break if read size has exceed the chunk size.
		if readSize > limit {
			break
		}

		b, err := reader.ReadBytes('\n')

		// Break if end of file is encountered.
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		readSize += int64(len(b))
		m1 := regexp.MustCompile(`(\\n|\\r|\\t)|[^0-9a-zA-Z']`)
		s := strings.ToLower(m1.ReplaceAllString(
			string(b), " "))
		p := strings.Fields(s)
		for i := 0; i < len(p); i++ {
			//ensuring I'm only grabbing 3 words at a time
			if i+3 > len(p) {
				break
			}
			channel <- strings.Join(p[i:i+3], " ")
		}

	}
	return nil
}
