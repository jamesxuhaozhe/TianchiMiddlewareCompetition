package engine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils/ds"
	"io"
	"net/http"
	"strings"
)

const (
	// we initialize a ring buffer with 15 slots to hold trace spans
	batchCount = 15
)

var (
	batchTraceList = make([]*ds.ConcurMap, 0, batchCount)
	initDone       = make(chan struct{}, 1)
)

// Init populates the data structure we need for further processing.
func Init() {
	go func() {
		for i := 0; i < batchCount; i++ {
			batchTraceList = append(batchTraceList, ds.NewConcurMap(constants.BatchSize))
		}
		fmt.Printf("batch trace list is populated. len is %d, cap is %d\n",
			len(batchTraceList), cap(batchTraceList))
		initDone <- struct{}{}
	}()
}

// processData executes the core logic of the client process, polling data from the data source and
func ProcessData() error {
	// wait until Init is done
	<-initDone

	// start polling the data from the data source
	url := getUrl()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("Start download from url: %s\n", url)

	buf := bufio.NewReader(resp.Body)
	count := 0
	pos := 0
	badTraceIdSet := ds.NewStrSet()
	traceBatchMap := batchTraceList[pos]
	for {
		count++
		line, err := buf.ReadString('\n')
		cols := strings.Split(line, "|")
		if cols != nil && len(cols) > 1 {
			traceId := cols[0]
			var spanList []string
			_, ok := traceBatchMap.Get(traceId)
			if !ok {
				spanList = make([]string, 1000)
				traceBatchMap.Put(traceId, &spanList)
			}
			spanList = append(spanList, line)
			if len(cols) > 8 {
				tag := cols[8]
				if isBadSpan(&tag) {
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
			badTraceIdSetBatchPos := count/constants.BatchSize - 1
			fmt.Printf("batch size: %d, badTraceSet size: %d, count: %d, badTraceIdSetBatchPos: %d\n",
				traceBatchMap.Size(), badTraceIdSet.Size(), count, badTraceIdSetBatchPos)
			sendBadTraceIds(badTraceIdSet.GetStrSlice(), badTraceIdSetBatchPos)
			badTraceIdSet.Clear()
		}
		if err != nil {
			if err == io.EOF {
				// done reading all lines
				break
			}
			return err
		}
	}
	sendBadTraceIds(badTraceIdSet.GetStrSlice(), count/constants.BatchSize - 1)
	markFinish()
	fmt.Printf("Total span count: %d\n", count)
	return nil
}

func sendBadTraceIds(traceIds []string, batchPos int) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["ids"] = traceIds
	data["batchPos"] = batchPos
	bytesData, _ := json.Marshal(data)
	fmt.Printf("request body: %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://"+constants.CommonUrlPrefix+constants.BackendProcessPort1+
		"/setBadTraceIds", bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}

// isBadSpan checks if the given tag belongs to a bad span.
// a bad span is defined as whose tags contains 'error=1' or 'http.status_code!=200'
func isBadSpan(tag *string) bool {
	if strings.Contains(*tag, "error=1") {
		return true
	} else if strings.Contains(*tag, "http.status_code=") && !strings.Contains(*tag, "http.status_code=200") {
		return true
	}
	return false
}

// markFinish informs the backend process server that the client has finished its job.
func markFinish() bool {
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
		return "http://localhost:" + conf.GetDatasourcePort() + "/trace1.data"
	} else if svrPort == constants.ClientProcessPort2 {
		return "http://localhost:" + conf.GetDatasourcePort() + "/trace2.data"
	}
	return ""
}
