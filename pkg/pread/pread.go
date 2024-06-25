package pread

import (
	"context"
	"sync"
	"syscall"
)

func Read(ctx context.Context, path string, readers int, chunkSize int64) (chunksCh chan []byte, closeFunc func(), err error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}

	var stat syscall.Stat_t
	if err = syscall.Fstat(fd, &stat); err != nil {
		func() { _ = syscall.Close(fd) }()
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	closeFunc = func() {
		cancel()
		_ = syscall.Close(fd)
	}

	chunksCh = make(chan []byte, readers)
	offsetsCh := make(chan int64, readers)

	wg := &sync.WaitGroup{}
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for offset := range offsetsCh {
				buffer := make([]byte, chunkSize)

				select {
				case <-ctx.Done():
					return
				default:
					n, err := syscall.Pread(fd, buffer, offset)
					if err != nil {
						panic(err)
					} else if n == 0 {
						return
					} else if int64(n) < chunkSize {
						buffer = buffer[:n]
					}

					chunksCh <- buffer
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
