package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/lxn/win"
)

type CounterResult struct {
	Name  string
	Value string
}

// APIV1OSSpecificHandler creates a http response for the API /info path
func APIV1OSSpecificHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reponseWriter.WriteHeader(http.StatusOK)

	counter := "\\Memory\\Available MBytes"

	getValue := getResult(counter)

	resultJSON, _ := json.Marshal(getValue)
	reponseWriter.Write(resultJSON)
}

func getResult(counter string) CounterResult {
	counterValue, _ := ReadPerformanceCounter(counter)

	return counterValue
}

func ReadPerformanceCounter(counter string) (CounterResult, error) {

	returnvalue := CounterResult{}

	var queryHandle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER

	ret := win.PdhOpenQuery(0, 0, &queryHandle)
	if ret != win.ERROR_SUCCESS {
		return returnvalue, errors.New("unable to open query through dll call")
	}

	// test path
	ret = win.PdhValidatePath(counter)
	if ret == win.PDH_CSTATUS_BAD_COUNTERNAME {
		return returnvalue, errors.New("unable to fetch counter (this is unexpected)")
	}

	ret = win.PdhAddEnglishCounter(queryHandle, counter, 0, &counterHandle)
	if ret != win.ERROR_SUCCESS {
		return returnvalue, errors.New(fmt.Sprintf("unable to add process counter. error code is %x\n", ret))
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret != win.ERROR_SUCCESS {
		return returnvalue, errors.New(fmt.Sprintf("got an error: 0x%x\n", ret))
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret == win.ERROR_SUCCESS {

		var bufSize uint32
		var bufCount uint32
		var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
		var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.

		ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &emptyBuf[0])
		if ret == win.PDH_MORE_DATA {
			filledBuf := make([]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
			ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &filledBuf[0])
			if ret == win.ERROR_SUCCESS {
				for i := 0; i < int(bufCount); i++ {
					c := filledBuf[i]
					s := win.UTF16PtrToString(c.SzName)

					metricName := counter
					if len(s) > 0 {
						metricName = fmt.Sprintf("%s.%s", counter, s)
					}

					returnvalue.Name = metricName
					returnvalue.Value = fmt.Sprintf("%f", c.FmtValue.DoubleValue)
				}
			}
		}
	}

	return returnvalue, nil

}
