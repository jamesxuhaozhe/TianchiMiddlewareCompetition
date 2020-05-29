package ds

import "sync"

// ConcurMap is a customized map that supports concurrent read and write.
type ConcurMap struct {
	cap   int
	rwMap map[string]*[]string
	mu    *sync.RWMutex
	sem   chan struct{}
}

// NewConcurMap returns a pointer of type ConcurMap.
func NewConcurMap(cap int) *ConcurMap {
	return &ConcurMap{
		cap:   cap,
		rwMap: make(map[string]*[]string, cap),
		mu:    &sync.RWMutex{},
		sem:   make(chan struct{}, 1),
	}
}

// Size returns the size of the map.
func (m *ConcurMap) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.rwMap)
}

// Put puts a key value pair to the map.
func (m *ConcurMap) Put(key string, value *[]string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.rwMap[key]
	m.rwMap[key] = value
	return ok
}

// Get gets a value for the given key, returns true if the key exists
// false otherwise.
func (m *ConcurMap) Get(key string) (*[]string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.rwMap[key]
	return value, ok
}

// Clear clears out the map content.
func (m *ConcurMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rwMap = make(map[string]*[]string, m.cap)
	//m.Done()
}

func (m *ConcurMap) Done() {
	m.sem <- struct{}{}
}

func (m *ConcurMap) Wait() {
	<-m.sem
}
