# Reading via Mmap obviously uses syscall.Mmap to improve the efficiency of multi-threaded reading of a file in chunks.

<img width="626" alt="image" src="https://github.com/Borislavv/go-mmap/assets/50691459/b7c81b57-008a-4050-ad1b-647698098027">


# How is it fast:
  - ./main  0,01s user 0,01s system 96% cpu 0,012 total
  - ./main  0,01s user 0,01s system 99% cpu 0,016 total
  - ./main  0,01s user 0,01s system 98% cpu 0,021 total

Approximately it's working for 0,017ms.
