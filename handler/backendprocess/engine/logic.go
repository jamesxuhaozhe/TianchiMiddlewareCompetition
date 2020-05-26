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
	initDone        = make(chan struct{}, 1)

	currentBatch   = 0
	availableBatch = make(chan *BadTraceIdsBatch)

	backendDone = make(chan struct{}, 1)
)

type BadTraceIdsBatch struct {
	batchPos     int
	badTraceIds  []string
	processCount int
}

func Init() {
	go func() {
		for i := 0; i < batchSize; i++ {
			badTraceIdsList = append(badTraceIdsList, &BadTraceIdsBatch{})
		}
		initDone <- struct{}{}
	}()

	go func() {
		for {
			nextBatch := currentBatch + 1
			if nextBatch >= batchSize {
				nextBatch = 0
			}

			nextBadTraceIdsBatch := badTraceIdsList[nextBatch]
			currentBadTraceIdsBatch := badTraceIdsList[currentBatch]
			if (finishProcessCount >= constants.ExpectedProcessCount && currentBadTraceIdsBatch.batchPos > 0) ||
				(currentBadTraceIdsBatch.processCount >= constants.ExpectedProcessCount &&
					nextBadTraceIdsBatch.processCount >= constants.ExpectedProcessCount) {
				badTraceIdsList[currentBatch] = &BadTraceIdsBatch{}
				currentBatch = nextBatch
				availableBatch <- currentBadTraceIdsBatch
			}

			select {
			case <-backendDone:
				break
			default:
				// do nothing
			}
		}
	}()

	go func() {
		for {
			select {
			// if we manage get one available batch from the channel
			case batch := <-availableBatch:
				process(batch)
			default:
				if IsFinished() {
					sendCheckSum()
				}
			}
		}
	}()

}

// sendCheckSum computes the desired MD5 checksum results and send it to the data source
func sendCheckSum() {

}

func process(batch *BadTraceIdsBatch) {

}

// SetBadTraceIds maps the incoming bad trace ids into a ring buffer.
func SetBadTraceIds(badTraceIds []string, batchPos int) {
	pos := batchPos % batchSize
	batch := badTraceIdsList[pos]
	if len(badTraceIds) > 0 {
		batch.batchPos = batchPos
		batch.processCount++
		//before := len(batch.badTraceIds)
		batch.badTraceIds = append(batch.badTraceIds, badTraceIds...)
		//after := len(batch.badTraceIds)
		/*		fmt.Printf("Add ids len is: %d, BatchPos: %d, pos: %d, bad ids before is: %d, after ids after is %d\n",
				len(badTraceIds), batchPos, pos, before, after)*/
	}
}

// BumpProcessCount bumps up the finish process count by 1.
func BumpProcessCount() {
	mu.Lock()
	defer mu.Unlock()
	finishProcessCount++
}

// StartCheckSumService starts the service computing the checksum.
func StartCheckSumService() {
	<-initDone
	go func() {
		for {
			if IsFinished() {
				fmt.Println("from checksum service: isFinished!")
				break
			}
		}
	}()
}

// IsFinished checks if there is really no more work for us to do before we can send the md5 info to data source
// 1. if we still have badTrace batch waiting for process, then it doesn't count as finish
// 2. if we don't have all the finish signals from the client, then it doesn't count as finish
func IsFinished() bool {
	// check if all the batch in the badTraceIdsList has been processed
	for _, v := range badTraceIdsList {
		if v.batchPos != 0 {
			return false
		}
	}
	mu.RLock()
	defer mu.RUnlock()
	// checks if we have received all the finish signal from the client
	return finishProcessCount == constants.ExpectedProcessCount
}
