package queue

import "sync"

type FileQueue struct {
	mutex   sync.Mutex
	queue   []string
	count   int
	current string
}
