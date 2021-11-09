package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	// for each word in "needles", find that word (and its index) in "haystack"

	haystack := loadFileLines("words_alpha_shuffled.txt")
	needles := loadFileLines("150k_needles.txt")



	now := time.Now()

	concurrency := runtime.NumCPU()
	batchSize := len(needles) / concurrency

	res := make(map[string]int, len(needles))
	resLock := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			start := i*batchSize
			upto := (i+1)*batchSize
			batchResult := findIn(haystack, needles[start:upto])
			resLock.Lock()
			for k, v := range batchResult {
				res[k] = v
			}
			resLock.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Took %v\n", time.Since(now))

	err := writeOutputToFile(res, "out.txt")
	if err != nil {
		log.Fatal(err)
	}
}


func findIn(haystack []string, needles []string) map[string]int {
	res := map[string]int{}
	for i, needle := range needles {
		if i % 10000 == 0 {
			log.Println(i, "needles found so far...")
		}
		for j, blade := range haystack {
			if blade == needle {
				res[blade] = j
			}
		}
	}
	return res
}

func loadFileLines(fileName string) []string {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}
	return strings.Split(string(fileBytes), "\r\n")
}

func writeOutputToFile(res map[string]int, fileName string) error {
	f, err := os.Create(fileName)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("error closing the file", err)
		}
	}(f)

	if err != nil {
		return err
	}

	for word, pos := range res {
		_, err := fmt.Fprintf(f, "%s %d\n", word, pos)
		if err != nil {
			return err
		}
	}
	return nil
}
