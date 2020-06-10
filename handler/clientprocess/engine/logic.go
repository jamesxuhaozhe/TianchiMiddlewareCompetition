package engine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/log"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils/ds"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// we initialize a ring buffer with 15 slots to hold trace spans
	batchCount = 15
)

var (
	batchTraceList = make([]*ds.ConcurMap, 0, batchCount)
	initDone       = make(chan struct{}, 1)
	start          time.Time
	after          time.Time
	lineChan       = make(chan string, 300)
)

// Start populates the data structure we need for further processing.
func Start() {
	go func() {
		for i := 0; i < batchCount; i++ {
			batchTraceList = append(batchTraceList, ds.NewConcurMap(constants.BatchSize))
		}
		log.Infof("batch trace list is populated. len is %d, cap is %d",
			len(batchTraceList), cap(batchTraceList))
		initDone <- struct{}{}
	}()
}

// processData executes the core logic of the client process, polling data from the data source and
func ProcessData() error {
	// wait until Start is done
	<-initDone

	// spin up one goroutine to process available data
	go func() {
		start = time.Now()
		count := 0
		pos := 0
		badTraceIdSet := ds.NewStrSet()
		traceBatchMap := batchTraceList[pos]
		for line := range lineChan {
			count++
			cols := strings.Split(line, "|")
			if cols != nil && len(cols) > 1 {
				traceId := cols[0]
				var spanList []string
				existSpans, ok := traceBatchMap.Get(traceId)
				if !ok {
					spanList = make([]string, 0, 50)
					spanList = append(spanList, line)
					traceBatchMap.Put(traceId, &spanList)
				} else {
					*existSpans = append(*existSpans, line)
				}
				if len(cols) > 8 {
					tag := cols[8]
					if isBadSpan(tag) {
						badTraceIdSet.Add(traceId)
					}
				}
			}
			if count%constants.BatchSize == 0 {
				pos++
				if pos >= batchCount {
					pos = 0
				}
				traceBatchMap = batchTraceList[pos]

				if traceBatchMap.Size() > 0 {
					for {
						time.Sleep(5 * time.Millisecond)
						if traceBatchMap.Size() == 0 {
							break
						}
					}
				}
				badTraceIdSetBatchPos := count/constants.BatchSize - 1
				sendBadTraceIds(badTraceIdSet.GetStrSlice(), badTraceIdSetBatchPos)
				badTraceIdSet.Clear()
			}
		}
		sendBadTraceIds(badTraceIdSet.GetStrSlice(), count/constants.BatchSize-1)
		markFinish()
		after = time.Now()
		log.Infof("Duration Delta: %v", after.Sub(start))
	}()

	// start polling the data from the data source
	url := getUrl()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := bufio.NewReader(resp.Body)
	for {
		//count++
		line, err := buf.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		lineChan <- line

		if err != nil {
			if err == io.EOF {
				// done reading all lines
				close(lineChan)
				log.Info("Done reading all lines")
				break
			}
			close(lineChan)
			log.Error("Unknown error occurred when reading lines.")
			return err
		}
	}
	return nil
}

func GetSpansForBadTraceId(badIds []string, batchPos int) (map[string]*[]string, error) {
	pos := batchPos % batchCount
	previous := pos - 1
	if previous == -1 {
		previous = batchCount - 1
	}
	next := pos + 1
	if next >= batchCount {
		next = 0
	}
	resultMap := make(map[string]*[]string, len(badIds))
	getSpansForBadTraceIds(previous, badIds, &resultMap)
	getSpansForBadTraceIds(pos, badIds, &resultMap)
	getSpansForBadTraceIds(next, badIds, &resultMap)
	batchTraceList[previous].Clear()
	return resultMap, nil
}

func getSpansForBadTraceIds(batchPos int, badIds []string, resultMap *map[string]*[]string) {
	batchMap := batchTraceList[batchPos]
	for _, badId := range badIds {
		spansList, _ := batchMap.Get(badId)
		if spansList != nil {
			var existSpanList []string
			resultMapValue := *resultMap
			if existSpans, ok := resultMapValue[badId]; ok {
				*existSpans = append(*existSpans, *spansList...)
			} else {
				existSpanList = append(existSpanList, *spansList...)
				resultMapValue[badId] = &existSpanList
			}
		}
	}
}

// sendBadTraceIds sends the info to client for answers.
func sendBadTraceIds(traceIds []string, batchPos int) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["ids"] = traceIds
	data["batchPos"] = batchPos
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "http://"+constants.CommonUrlPrefix+constants.BackendProcessPort1+
		"/setBadTraceIds", bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("sendBadTraceIds error, batchPos: %d", batchPos)
		return
	}
	defer resp.Body.Close()
}

// isBadSpan checks if the given tag belongs to a bad span.
// a bad span is defined as whose tags contains 'error=1' or 'http.status_code!=200'
func isBadSpan(tag string) bool {
	if strings.Contains(tag, "error=1") {
		return true
	} else if strings.Contains(tag, "http.status_code=") &&
		!strings.Contains(tag, "http.status_code=200") {
		return true
	}
	return false
}

// markFinish informs the backend process server that the client has finished its job.
func markFinish() bool {
	log.Info("markFinish gets called")
	resp, err := http.Get("http://" + constants.CommonUrlPrefix + constants.BackendProcessPort1 + "/markFinish")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

// getUrl get the download url for the given server instance.
func getUrl() string {
	svrPort := conf.GetServerPort()
	if svrPort == constants.ClientProcessPort1 {
		return "http://localhost:" + conf.GetLocalTestPort() + "/trace1.data"
	} else if svrPort == constants.ClientProcessPort2 {
		return "http://localhost:" + conf.GetLocalTestPort() + "/trace2.data"
	}
	return ""
}
