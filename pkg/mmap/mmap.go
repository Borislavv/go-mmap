package mmap

import (
	"context"
	"sync"
	"syscall"
)

// Read is take a path to file, number of readers and size of chunk.
func Read(ctx context.Context, path string, readers int, chunkSize int64) (chunksCh chan []byte, closeFunc func(), err error) {
	pageSize := int64(syscall.Getpagesize())
	if chunkSize%pageSize != 0 {
		chunkSize = ((chunkSize / pageSize) + 1) * pageSize
	}

	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}

	var stat syscall.Stat_t
	if err = syscall.Fstat(fd, &stat); err != nil {
		func() { _ = syscall.Close(fd) }()
		return nil, nil, err
	}

	data, err := syscall.Mmap(fd, 0, int(stat.Size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		func() { _ = syscall.Close(fd) }()
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	closeFunc = func() {
		cancel()
		_ = syscall.Close(fd)
		_ = syscall.Munmap(data)
	}

	chunksCh = make(chan []byte, readers)
	offsetsCh := make(chan int64, readers)

	wg := sync.WaitGroup{}
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for offset := range offsetsCh {
				select {
				case <-ctx.Done():
					return
				default:
					end := offset + chunkSize
					if end > stat.Size {
						end = stat.Size
					}
					chunksCh <- data[offset:end]
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(chunksCh)
	}()

	go func() {
		for offset := int64(0); offset < stat.Size; offset += chunkSize {
			select {
			case <-ctx.Done():
				break
			default:
				offsetsCh <- offset
			}
		}
		close(offsetsCh)
	}()

	return chunksCh, closeFunc, nil
}
