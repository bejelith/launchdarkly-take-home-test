package db

import "sync"

// Average keeps track of average value of groups of samples
type Average[T comparable] struct {
	samples map[T][]float64
	sums    map[T]float64
	rw      *sync.RWMutex
}

// New builds  a new Average DB
func New[T comparable]() *Average[T] {
	return &Average[T]{
		map[T][]float64{},
		map[T]float64{},
		&sync.RWMutex{},
	}
}

// Add adds a sample to a set
func (s *Average[T]) Add(set T, r float64) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.samples[set] = append(s.samples[set], r)
	s.sums[set] += r
}

// Get retrieves the average for a set of samples
func (s *Average[T]) Get(key T) ([]float64, float64) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.samples[key], s.sums[key] / float64(len(s.samples[key]))
}

// GetAll retrieves all recorded samples for a set
func (s *Average[T]) GetAll() []T {
	s.rw.RLock()
	defer s.rw.RUnlock()
	output := make([]T, 0, len(s.samples))
	for k := range s.samples {
		output = append(output, k)
	}
	return output
}
