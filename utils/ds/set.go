package ds

// StrSet is a customized data structure to hold a set of string.
type StrSet struct {
	set map[string]bool
	cap int
}

// NewStrSet creates an instance of StrSet.
func NewStrSet(cap int) *StrSet {
	return &StrSet{
		set: make(map[string]bool, cap),
		cap: cap,
	}
}

// Add adds an element to the set, return false if the element already exists in the set,
// true otherwise
func (s *StrSet) Add(element string) bool {
	if !s.set[element] {
		s.set[element] = true
		return true
	}
	return false
}

// Contains checks if the given element exists in the set.
func (s *StrSet) Contains(element string) bool {
	return s.set[element]
}

// Clear clears out the set.
func (s *StrSet) Clear() {
	s.set = make(map[string]bool, s.cap)
}

// Size returns number of the elements in the set.
func (s *StrSet) Size() int {
	return len(s.set)
}
