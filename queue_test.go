package pqueue_test

import (
	"heckel.io/pqueue"
	"testing"
)

func TestQueue_Simple(t *testing.T) {
	q, _ := pqueue.New(t.TempDir())
	if err := q.EnqueueString("entry 1"); err != nil {
		t.Fatal()
	}
	if err := q.EnqueueString("entry 2"); err != nil {
		t.Fatal()
	}
	if s, _ := q.DequeueString(); s != "entry 1" {
		t.Fatal()
	}
	if s, _ := q.DequeueString(); s != "entry 2" {
		t.Fatal()
	}
	if _, err := q.DequeueString(); err != pqueue.ErrEmpty {
		t.Fatal()
	}
}

func TestQueue_PersistAndReRead(t *testing.T) {
	dir := t.TempDir()

	// Persist stuff (first use)
	q1, _ := pqueue.New(dir)
	q1.EnqueueString("entry 1")
	q1.EnqueueString("entry 2")

	// Re-read stuff (next use)
	q2, _ := pqueue.New(dir)
	if s, _ := q2.DequeueString(); s != "entry 1" {
		t.Fatal()
	}
	if s, _ := q2.DequeueString(); s != "entry 2" {
		t.Fatal()
	}
	if _, err := q2.DequeueString(); err != pqueue.ErrEmpty {
		t.Fatal()
	}
}

func TestQueue_NonExistingDir(t *testing.T) {
	if _, err := pqueue.New("/does/not/exist"); err == nil {
		t.Fatal()
	}
}
