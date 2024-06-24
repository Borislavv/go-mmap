package main

import (
	"fmt"
	"github.com/Borislavv/go-mmap/pkg/mmap"
	"sync"
)

const (
	// on the 1.000.000.000 rows it's working approximately for 0.016-0.024ms.
	// on my macbook M2 Max for [./main  0.01s user 0.02s system 160% cpu 0.016 total]
	path          = "1000000000rows.txt"
	workerCount   = 24
	readBlockSize = 2048 * 2048
)

func main() {
	chunksCh, closeFunc, err := mmap.Read(path, workerCount, readBlockSize)
	if err != nil {
		panic(err)
	}
	defer closeFunc()

	wg := &sync.WaitGroup{}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunksCh {
				_ = chunk
				//processData(chunk)
			}
		}()
	}

	wg.Wait()

	fmt.Println("Reading completed")
}

func processData(data []byte) {
	l := len(data)
	for i := 0; i < l; i++ {
		if data[i] == '\n' {

		}
	}
}
