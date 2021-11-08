# Try Goroutines

`findIn(haystack []string, needles []string) map[string]int` is a function that looks for each needle from `needles` in the `hackstack`, and records the index at which it was found.

`haystack` and `needles` are lists of English words, loaded from the .txt files provided.


To run the serial implementation, use `go run main.go`

It takes about 68s to run on a 2020 M1 Mac Mini.

Without changing the `findIn` function, try to make the program concurrent so it can run at a specified concurrency level. 

Hint: goroutines


