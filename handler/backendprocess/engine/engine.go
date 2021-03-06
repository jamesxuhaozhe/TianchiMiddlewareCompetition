package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/log"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils/ds"
	"net/http"
	"net/url"
	"sync"
)

const (
	// we use 65 batches to cache the total bad trace id from two clients
	batchSize = 65
)

var (
	mu                 = &sync.RWMutex{}
	finishProcessCount int8

	badTraceIdsList = make([]*BadTraceIdsBatch, 0, batchSize)
	initDone        = make(chan struct{})

	currentBatch = 0

	csMu        = &sync.Mutex{}
	checkSumMap = make(map[string]string, 10000)

	ports = []string{constants.ClientProcessPort1, constants.ClientProcessPort2}
)

type BadTraceIdsBatch struct {
	batchPos     int
	badTraceIds  []string
	processCount int
}

type response struct {
	Map map[string]*[]string `json:"map"`
}

// Start spins up the core backend processing logic
func Start() {
	go func() {
		for i := 0; i < batchSize; i++ {
			badTraceIdsList = append(badTraceIdsList, &BadTraceIdsBatch{})
		}
		initDone <- struct{}{}
	}()

	go func() {
		<-initDone
		count := 0
		for {
			count++
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
				process(currentBadTraceIdsBatch)
			}
			if IsFinished() {
				break
			}
		}

		csMu.Lock()
		checkSumJson, _ := json.Marshal(checkSumMap)
		csMu.Unlock()
		checkSumString := string(checkSumJson)
		resp, err := http.PostForm("http://localhost:"+conf.GetDatasourcePort()+"/api/finished",
			url.Values{"result": {checkSumString}})
		if err != nil {
			log.Error("send checkSum fail")
		} else {
			log.Info("Suc to send checksum!")
			defer resp.Body.Close()
		}

		log.Info("exiting the second goroutine!")
	}()
}

// process takes in the given batch id sets and query two remote servers for spans we are interested in.
func process(batch *BadTraceIdsBatch) {
	traceMap := make(map[string]*ds.StrSet)
	for _, port := range ports {
		tempTraceMap, err := getTraceMapFromRemote(batch.badTraceIds, batch.batchPos, port)
		if err != nil {
			log.Errorf("getTraceMapFromRemote error: batchPos: %d, port: %d", batch.batchPos, port)
			continue
		}
		for traceId, spanList := range tempTraceMap {
			if spanSet, ok := traceMap[traceId]; ok {
				spanSet.AddAll(*spanList)
			} else {
				spanSet = ds.NewStrSetWithCap(50)
				spanSet.AddAll(*spanList)
				traceMap[traceId] = spanSet
			}
		}
	}
	for traceId, spans := range traceMap {
		spanStr := spans.SortedStr() + "\n"
		md5Hash := utils.MD5(spanStr)
		csMu.Lock()
		checkSumMap[traceId] = md5Hash
		csMu.Unlock()
	}
}

// getTraceMapFromRemote queries the remote client for answers
func getTraceMapFromRemote(badTraceIds []string, batchPos int, port string) (map[string]*[]string, error) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["ids"] = badTraceIds
	data["batchPos"] = batchPos
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "http://"+constants.CommonUrlPrefix+port+
		"/getSpansForBadTraceIds", bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("net work error")
	}

	var response = response{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Map, nil
}

// SetBadTraceIds maps the incoming bad trace ids into a ring buffer.
func SetBadTraceIds(badTraceIds []string, batchPos int) {
	pos := batchPos % batchSize
	batch := badTraceIdsList[pos]

	batch.batchPos = batchPos
	batch.processCount++
	batch.badTraceIds = append(batch.badTraceIds, badTraceIds...)
}

// BumpProcessCount bumps up the finish process count by 1.
func BumpProcessCount() {
	mu.Lock()
	defer mu.Unlock()
	finishProcessCount++
}

// IsFinished checks if there is really no more work for us to do before we can send the md5 info to data source
// 1. if we still have badTrace batch waiting for process, then it doesn't count as finish
// 2. if we don't have all the finish signals from the client, then it doesn't count as finish
func IsFinished() bool {
	// check if all the batch in the badTraceIdsList has been processed
	for i := 0; i < batchSize; i++ {
		if badTraceIdsList[i].batchPos != 0 {
			return false
		}
	}
	mu.RLock()
	defer mu.RUnlock()
	// checks if we have received all the finish signal from the client
	return finishProcessCount == constants.ExpectedProcessCount
}
