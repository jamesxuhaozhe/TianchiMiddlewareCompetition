package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/log"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils/ds"
	"net/http"
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
	initDone        = make(chan struct{})

	currentBatch   = 0
	availableBatch = make(chan *BadTraceIdsBatch)

	csMu = &sync.Mutex{}
	checkSumMap = make(map[string]string)
)

type BadTraceIdsBatch struct {
	batchPos     int
	badTraceIds  []string
	processCount int
}

type response struct {
	Map map[string]*[]string `json:"map"`
}

func Init() {
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
				//fmt.Printf("Send batchpos: %d\n", currentBadTraceIdsBatch.batchPos)
				availableBatch <- currentBadTraceIdsBatch
				//fmt.Printf("finish sending batchpos: %d\n", currentBadTraceIdsBatch.batchPos)
			}
			if IsFinished() {
				close(availableBatch)
				break
			}
		}
		log.Info("exiting the second goroutine!")
	}()

	go func() {
		log.Info("Entering second goroutine.")

		ports := []string{constants.ClientProcessPort1} //, constants.ClientProcessPort2}
		for batch := range availableBatch {
			process(batch, &ports)
		}
		log.Info("Exiting second goroutine.!!!!!!")
	}()

}

// sendCheckSum computes the desired MD5 checksum results and send it to the data source
func sendCheckSum() {
	log.Info("Send check sum method invoked")
}

func process(batch *BadTraceIdsBatch, ports *[]string) {
	log.Infof("process batchPos: %d", batch.batchPos)
	traceMap := make(map[string]*ds.StrSet)
	for _, port := range *ports {
		tempTraceMap, err := getTraceMapFromRemote(batch.badTraceIds, batch.batchPos, port)
		if err != nil {
			log.Errorf("getTraceMapFromRemote error: batchPos: %d, port: %d", batch.batchPos, port)
			continue
		}
		for traceId, spanList := range tempTraceMap {
			if spanSet, ok := traceMap[traceId]; ok {
				spanSet.AddAll(*spanList)
			} else {
				spanSet = ds.NewStrSet()
				spanSet.AddAll(*spanList)
				traceMap[traceId] = spanSet
			}
		}
	}
	//log.Infof("traceMap: %s", traceMap)
	//getTraceMapFromRemote(batch.badTraceIds, batch.batchPos, "")
	// update the checksum map
	for traceId, spans := range traceMap {

		spanStr := spans.SortedStr() + "\n"
		csMu.Lock()
		// TODO we need to get md5 digest of the spanstr
		checkSumMap[traceId] = spanStr
		csMu.Unlock()
	}
}

func getTraceMapFromRemote(badTraceIds []string, batchPos int, port string) (map[string]*[]string, error) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["ids"] = badTraceIds
	data["batchPos"] = batchPos
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "http://"+constants.CommonUrlPrefix+port+
		"/getSpansForBadTraceIds", bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("net work error")
	}

	var respo = response{}
	if err := json.NewDecoder(resp.Body).Decode(&respo); err != nil {
		return nil, err
	}
	//log.Infof("get back data %v", respo)
	return respo.Map, nil
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
