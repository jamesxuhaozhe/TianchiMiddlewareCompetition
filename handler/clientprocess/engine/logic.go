package engine

import (
	"bufio"
	"fmt"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// processData executes the core logic of the client process, polling data from the data source and
func ProcessData() error {
	// start polling the data from the data source
	url := getUrl()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf := bufio.NewReader(resp.Body)
	count := 0

	fmt.Printf("Start download from url: %s\n", url)
	for {
		count++
		line, err := buf.ReadString('\n')
		//line = strings.TrimSpace(line)
		//fmt.Printf("Line %s, %s", strconv.Itoa(count), line)
		cols := strings.Split(line, "|")
		fmt.Printf("Line number: %s, Col length: %s \n", strconv.Itoa(count), strconv.Itoa(len(cols)))
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func getUrl() string {
	svrPort := conf.GetServerPort()
	if svrPort == constants.ClientProcessPort1 {
		return "http://localhost:" + conf.GetDatasourcePort() + "/trace1.data"
	} else if svrPort == constants.ClientProcessPort2 {
		return "http://localhost:" + conf.GetDatasourcePort() + "/trace2.data"
	}
	return ""
}