# pqueue
[![Go Reference](https://pkg.go.dev/badge/heckel.io/pqueue.svg)](https://pkg.go.dev/heckel.io/pqueue)
[![Tests](https://github.com/binwiederhier/pqueue/workflows/test/badge.svg)](https://github.com/binwiederhier/pqueue/actions)

`pqueue` is a simple persistent directory-backed FIFO queue.

It provides the typical queue interface `Enqueue` and `Dequeue` and may store any byte slice or string. 
Entries are stored as files in the backing directory and are fully managed by `Queue`.

## Usage
Fetch the library:
```
go get -u heckel.io/pqueue
```
Then import it like this: 

```go
import "heckel.io/pqueue"
```
The package is imported under the `pqueue` namespace.

## Example

```go
import "heckel.io/pqueue"

q, _ := pqueue.New("/tmp/myqueue")
q.EnqueueString("some entry")
first, _ := q.DequeueString() // some entry

q2, _ := pqueue.New("/tmp/myqueue") // queue persists
second, _ := q2.DequeueString() // another entry
```

## Contributing
I welcome any and all contributions. Just create a PR or an issue.

## License
Made with ❤️ by [Philipp C. Heckel](https://heckel.io), distributed under the [Apache License 2.0](LICENSE).
