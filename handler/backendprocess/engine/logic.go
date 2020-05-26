package engine

import (
	"fmt"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"sync"
)

const (
	// we use 90 batches to cache the total bad trace id from two clients
	batchSize = 65
)

var (
	mu                 = &sync.RWMutex{}
	finishProcessCount int8

	badTraceIdsList = make([]*BadTraceIdsBatch, 0, batchSize)
)

type BadTraceIdsBatch struct {
	batchPos     int
	badTraceIds  []string
	processCount int
}

func Init() {
	for i := 0; i < batchSize; i++ {
		badTraceIdsList = append(badTraceIdsList, &BadTraceIdsBatch{})
	}
}

// SetBadTraceIds maps the incoming bad trace ids into a ring buffer.
func SetBadTraceIds(badTraceIds []string, batchPos int) {
	pos := batchPos % batchSize
	batch := badTraceIdsList[pos]
	if len(badTraceIds) > 0 {
		batch.batchPos = batchPos
		batch.processCount++
		before := len(batch.badTraceIds)
		batch.badTraceIds = append(batch.badTraceIds, badTraceIds...)
		after := len(batch.badTraceIds)
		fmt.Printf("Add ids len is: %d, BatchPos: %d, pos: %d, bad ids before is: %d, after ids after is %d\n",
			len(badTraceIds), batchPos, pos, before, after)
	}
}

// BumpProcessCount bumps up the finish process count by 1.
func BumpProcessCount() {
	mu.Lock()
	defer mu.Unlock()
	finishProcessCount++
}

// IsFinished checks if the whole process is finished
func IsFinished() bool {
	mu.RLock()
	defer mu.RUnlock()
	return finishProcessCount == constants.ExpectedProcessCount
}
