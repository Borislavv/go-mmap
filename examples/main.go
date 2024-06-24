package main

import (
	"fmt"
	"github.com/Borislavv/go-mmap/pkg/mmap"
	"sync"
)

const (
	// on the 1.000.000.000 rows it's working approximately for 0.025ms.
	// on my macbook M2 Max for:
	//	- ./main  0,01s user 0,01s system 96% cpu 0,012 total
	//  - ./main  0,01s user 0,01s system 99% cpu 0,016 total
	//  - ./main  0,01s user 0,01s system 98% cpu 0,021 total
	path        = "1-000-000-000.txt"
	workerCount = 24
	chunkSize   = 4096 * 4096
)

func main() {
	chunksCh, closeFunc, err := mmap.Read(path, workerCount, chunkSize)
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
			}
		}()
	}
	wg.Wait()

	fmt.Println("Reading completed")
}
