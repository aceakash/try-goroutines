package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	haystack := loadFileLines("words_alpha_shuffled.txt")
	needles := loadFileLines("150k_needles.txt")

	res := map[string]int{}

	now := time.Now()
	findIn(haystack, needles, res)
	fmt.Printf("Took %v\n", time.Since(now))

	err := writeOutputToFile(res, "out.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func writeOutputToFile(res map[string]int, s string) error {
	f, err := os.Create("out.txt")
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

func findIn(haystack []string, needles []string, res map[string]int) {
	for i, needle := range needles {
		if i % 10000 == 0 {
			fmt.Println(i)
		}
		for j, blade := range haystack {
			if blade == needle {
				res[blade] = j
			}
		}
	}
}

func loadFileLines(fileName string) []string {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}
	return strings.Split(string(fileBytes), "\r\n")
}
