# Reading via Mmap obviously uses syscall.Mmap to improve the efficiency of multi-threaded reading of a file in chunks.

<img width="663" alt="image" src="https://github.com/Borislavv/go-mmap/assets/50691459/855cc218-1e5a-45ad-86ff-8f04ba746b68">


# How is it fast:
  - ./main  0,01s user 0,01s system 96% cpu 0,012 total
  - ./main  0,01s user 0,01s system 99% cpu 0,016 total
  - ./main  0,01s user 0,01s system 98% cpu 0,021 total

Approximately it's working for 0,017ms.

# P.S. 
But this method is ineffective if you have heavy work on the received data (that is, if you use data either from the end or from the beginning of the slice). This is due to the way the file is mapped into memory, thus making the processor cache less efficient.
In that case, I would recommend considering using syscall.Pread with offsets.
