// Package pqueue provides a simple persistent FIFO queue backed by a directory.
//
// Queue provides the typical queue interface Enqueue and Dequeue and may store any byte slice.
// Entries are stored as files in the backing directory and are fully managed by Queue.
//
// Example:
//   q1, _ := pqueue.New("/tmp/myqueue")
//   q1.EnqueueString("my entry")
//   q2, _ := pqueue.New("/tmp/myqueue")
//   myEntry, _ := q2.DequeueString()
//
package pqueue

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

// ErrEmpty is returned by Dequeue if the queue is empty
var ErrEmpty = errors.New("queue is empty")

// Queue is a persistent FIFO queue backed by a directory
type Queue struct {
	dir     string
	entries []int
	current int
	mu      sync.Mutex
}

// New creates a new persistent FIFO queue backed by the given directory.
//
// The directory must exist, or an error is returned. The queue is initialized using
// the backed directory, and re-reads previous keys into its internal buffer. No two
// Queue instances may use the same backing directory at the same time.
func New(dir string) (*Queue, error) {
	entries, err := readKeys(dir)
	if err != nil {
		return nil, err
	}
	var current int
	if len(entries) > 0 {
		current = entries[len(entries)-1]
	}
	return &Queue{
		dir:     dir,
		entries: entries,
		current: current,
	}, nil
}

// Enqueue writes a new byte slice to the queue and persists it as file in the backing directory
func (q *Queue) Enqueue(b []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.current++
	filename := filepath.Join(q.dir, strconv.Itoa(q.current))
	if err := os.WriteFile(filename, b, 0600); err != nil {
		return err
	}
	q.entries = append(q.entries, q.current)
	return nil
}

// EnqueueString writes a new string to the queue and persists it as file in the backing directory
func (q *Queue) EnqueueString(s string) error {
	return q.Enqueue([]byte(s))
}

// Dequeue returns the first entry in the queue as a byte slice, or returns ErrEmpty if the queue
// is empty. It also removes the entry file in the backing directory.
func (q *Queue) Dequeue() ([]byte, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.entries) == 0 {
		return nil, ErrEmpty
	}
	filename := filepath.Join(q.dir, strconv.Itoa(q.entries[0]))
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(filename); err != nil {
		return nil, err
	}
	q.entries = q.entries[1:]
	return b, nil
}

// DequeueString returns the first entry in the queue as a string, or returns ErrEmpty if the queue
// is empty. It also removes the entry file in the backing directory.
func (q *Queue) DequeueString() (string, error) {
	b, err := q.Dequeue()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readKeys(dir string) ([]int, error) {
	keys := make([]int, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		key, err := strconv.Atoi(entry.Name())
		if err == nil {
			keys = append(keys, key)
		}
	}
	sort.Ints(keys)
	return keys, nil
}
