package main

import (
	CusWaitGroup "app/WaitGroups"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

func countFileVowels(filePath string, t *atomic.Int32, wg *CusWaitGroup.WGroup) {
	defer wg.Done()
	const vowels = "aeiou"
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	for _, v := range file {
		c := strings.ToLower(string(v))
		if strings.ContainsAny(vowels, c) {
			t.Add(1)
		}
	}
}

func main() {
	var total atomic.Int32
	wg := CusWaitGroup.NewWGroup()
	wg.Add(4)
	for i := 1; i <= 4; i++ {
		filename := "data/" + fmt.Sprintf("file_%d.txt", i)
		go countFileVowels(filename, &total, wg)
	}
	wg.Wait()
	fmt.Printf("Total vowels in all files are: %d \n", total.Load())
}
