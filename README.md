# gtrace

A system call tracer for Linux x86-64.

DISCLAIMER: This software is experimental and not considered stable. Do
not use it in mission-critical environments.

## Installation

```golang
$ git clone git@github.com:agis/gtrace.git
$ cd gtrace
$ go generate && go build
```

## Usage

Currently only attaching to an already running process is supported. Also,
arguments are not decoded yet.

Attach to a process by specifying its pid:

```shell
$ ./gtrace -p 2602
Attached to process 2602...
futex = 0
write = 2
write = 1
futex = 0
write = 2
^C
```

## Feature work

- [ ] Decode arguments
- [ ] terminal GUI with live statistics/counters
- [ ] filter for certain syscalls
- [ ] ARM support

## Building

```shell
$ go generate
$ go build
```
