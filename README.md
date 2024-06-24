# Reading via Mmap obviously uses syscall.Mmap to improve the efficiency of multi-threaded reading of a file in chunks.

<img width="626" alt="image" src="https://github.com/Borislavv/go-mmap/assets/50691459/b7c81b57-008a-4050-ad1b-647698098027">


# How is it fast:
  - ./main  0,02s user 0,02s system 170% cpu 0,024 total
  - ./main  0,02s user 0,02s system 169% cpu 0,024 total
  - ./main  0,02s user 0,02s system 163% cpu 0,025 total

Approximately it's working for 0,025ms.
