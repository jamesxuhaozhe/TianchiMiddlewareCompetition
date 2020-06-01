package ds

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// StrSet is a customized data structure to hold a set of string.
type StrSet struct {
	set map[string]bool
	cap int
}

// NewStrSetWithCap creates an instance of StrSet.
func NewStrSetWithCap(cap int) *StrSet {
	return &StrSet{
		set: make(map[string]bool, cap),
		cap: cap,
	}
}

func NewStrSet() *StrSet {
	return NewStrSetWithCap(1)
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

// AddAll adds bulky elements.
func (s *StrSet) AddAll(elements []string) {
	for _, element := range elements {
		if !s.Contains(element) {
			s.Add(element)
		}
	}
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

func (s *StrSet) ToJSON() string {
	tempKeys := make([]string, 0, len(s.set))
	for k, _ := range s.set {
		tempKeys = append(tempKeys, k)
	}
	b, err := json.MarshalIndent(tempKeys, "", "  ")
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}


type SpanSlice []string

func (s SpanSlice) Len() int {
	return len(s)
}

func (s SpanSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SpanSlice) Less(i, j int) bool {
	timeStampI, timeStampJ := toInt64(s[i]), toInt64(s[j])
	return timeStampI < timeStampJ
}

func toInt64(span string) int64 {
	cols := strings.Split(span, "|")
	if len(cols) > 8 {
		timeStamp, err := strconv.ParseInt(cols[1], 10, 64)
		if err != nil {
			return -1
		}
		return timeStamp
	}
	return -1
}

func (s *StrSet) SortedStr() string {
	if len(s.set) == 0 {
		return ""
	}
	spanSlice := s.GetStrSlice()
	sort.Stable(SpanSlice(spanSlice))
	return strings.Join(spanSlice, "\n")
}

func (s *StrSet) GetStrSlice() []string {
	tempKeys := make([]string, 0, len(s.set))
	for k := range s.set {
		tempKeys = append(tempKeys, k)
	}
	return tempKeys
}
