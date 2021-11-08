package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	// for each word in "needles", find that word (and its index) in "haystack"

	haystack := loadFileLines("words_alpha_shuffled.txt")
	needles := loadFileLines("150k_needles.txt")

	// get maxproc to determine the optimal concurrency level, otherwise the scheduler will be fighting for CPU cores.
	// increasing this will lead to diminishing returns and make it a tad more unpredictable
	maxProcs := runtime.GOMAXPROCS(0)
	fmt.Println(fmt.Sprintf("maxprocs: %d", maxProcs))

	now := time.Now()

	var chunks [][]string
	// calculate the size of each chunk
	chunkSize := (len(needles) + maxProcs - 1) / maxProcs
	// split needles in chunks to search needles concurrently
	for i := 0; i < len(needles); i += chunkSize {
		end := i + chunkSize

		if end > len(needles) {
			end = len(needles)
		}

		chunks = append(chunks, needles[i:end])
	}

	// init map with the correct size to avoid resizes
	res := make(chan map[string]int, len(chunks))
	for _, chunk := range chunks {
		go func(c []string) {
			res <- findIn(haystack, c)
		}(chunk)
	}

	// init map with the correct size to avoid resizes
	store := make(map[string]int, len(needles))
	for i := 0; i < len(chunks); i++ {
		// maps in go are not concurrent so to add concurrency here we would probably need to lock anyway
		chunkRes := <-res
		for k, v := range chunkRes {
			store[k] = v
		}
	}

	fmt.Printf("Took %v\n", time.Since(now))

	err := writeOutputToFile(store, "out.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func findIn(haystack []string, needles []string) map[string]int {
	res := map[string]int{}
	for i, needle := range needles {
		if i%10000 == 0 {
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

func writeOutputToFile(res map[string]int, s string) error {
	f, err := os.Create(s)
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
