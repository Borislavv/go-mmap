package main

import (
	"fmt"
	"github.com/Borislavv/go-mmap/pkg/mmap"
	"sync"
)

const (
	// on the 1.000.000.000 rows it's working approximately for 0.025ms.
	// on my macbook M2 Max for:
	//	- ./main  0,02s user 0,02s system 170% cpu 0,024 total
	//  - ./main  0,02s user 0,02s system 169% cpu 0,024 total
	//  - ./main  0,02s user 0,02s system 163% cpu 0,025 total
	path          = "1-000-000-000-rows-file.txt"
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
			}
		}()
	}
	wg.Wait()

	fmt.Println("Reading completed")
}
