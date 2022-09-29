package ds

import (
	"fmt"
	"sync"
	"time"
)

type ResourceData struct {
	Meta struct {
		SendEventID      string `json:"sendEventID"`
		PageLocation     string `json:"pageLocation"`
		Host             string `json:"host"`
		Path             string `json:"path"`
		Query            string `json:"query"`
		Protocol         string `json:"protocol"`
		PageTitle        string `json:"pageTitle"`
		PCode            int    `json:"pCode"`
		ProjectAccessKey string `json:"projectAccessKey"`
		Screen           struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"screen"`
		SessionID string `json:"sessionID"`
		UserAgent string `json:"userAgent"`
		UserID    string `json:"userID"`
	} `json:"meta"`
	Resource []struct {
		StartTime      int    `json:"startTime"`
		StartTimeStamp int64  `json:"startTimeStamp"`
		EventID        string `json:"eventID"`
		Type           string `json:"type"`
		URL            string `json:"url"`
		URLHost        string `json:"urlHost"`
		URLPath        string `json:"urlPath"`
		URLQuery       string `json:"urlQuery"`
		URLProtocol    string `json:"urlProtocol"`
		Timing         struct {
			Redirect struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"redirect"`
			Cache struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"cache"`
			Connect struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"connect"`
			DNS struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"dns"`
			Ssl struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"ssl"`
			Download struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"download"`
			FirstByte struct {
				Duration int `json:"duration"`
				Start    int `json:"start"`
			} `json:"firstByte"`
			Duration int `json:"duration"`
			Size     int `json:"size"`
		} `json:"timing"`
		ResourceInfo struct {
			Method string `json:"method"`
			Status int    `json:"status"`
		} `json:"resourceInfo"`
		TraceInfo struct {
			MtID string `json:"mtID"`
			TxID string `json:"txID"`
		} `json:"traceInfo"`
	} `json:"resource"`
}

/*
rum_page_load_ajax_each_page
tags: page_path	request_host	request_path
fields:
ajax_2xx_count
ajax_3xx_count
ajax_4xx_count
ajax_5xx_count
ajax_count
ajax_duration
*/

type AJAXMap struct {
	Type string // resource 'r', pageload 'p'
	sync.RWMutex
	// request_host, request_path, pCode, page_path
	items map[string]map[string]map[int64]map[string][]AJAXStatus
	stop  chan bool
}

func NewAJAXMap(Type string) *AJAXMap {
	ajaxMap := new(AJAXMap)
	ajaxMap.Type = Type
	ajaxMap.stop = make(chan bool)
	ajaxMap.SendAJAXTagCountGoFunc()
	return ajaxMap
}

type AJAXStatus struct {
	Status        int
	Ajax_duration float32
}

type AJAXCount struct {
	ajax_count     int
	ajax_2xx_count int
	ajax_3xx_count int
	ajax_4xx_count int
	ajax_5xx_count int
	ajax_duration  float32
}

type AJAXCountUrlSum struct {
	url string
	sum AJAXCount
}

func (um *AJAXMap) GetStatsType() string {
	return um.Type
}

func (um *AJAXMap) Sum(ajaxStatuses []AJAXStatus) (total AJAXCount) {
	total = AJAXCount{0, 0, 0, 0, 0, 0}
	for _, v := range ajaxStatuses {
		total.ajax_count++
		total.ajax_duration += v.Ajax_duration
		if v.Status < 200 {
			fmt.Println("AJAX StatusCode Error(<200):", v.Status)
		} else if v.Status < 300 {
			total.ajax_2xx_count++
		} else if v.Status < 400 {
			total.ajax_3xx_count++
		} else if v.Status < 500 {
			total.ajax_4xx_count++
		} else if v.Status < 600 {
			total.ajax_5xx_count++
		} else {
			fmt.Println("AJAX StatusCode Error(>500):", v.Status)
		}
	}
	if total.ajax_count != 0 {
		total.ajax_duration /= float32(total.ajax_count)
	}
	return total
}

func (um *AJAXMap) GetPcodeSums(request_host, request_path string, pCode int64) (total AJAXCount) {
	um.RLock()
	defer um.RUnlock()
	total = AJAXCount{0, 0, 0, 0, 0, 0}

	if um.items != nil {
		for _, v := range um.items[request_host][request_path][pCode] {
			for _, v1 := range v {
				total.ajax_count++
				total.ajax_duration += v1.Ajax_duration
				if v1.Status < 200 {
					fmt.Println("AJAX StatusCode Error(<200):", v1.Status)
				} else if v1.Status < 300 {
					total.ajax_2xx_count++
				} else if v1.Status < 400 {
					total.ajax_3xx_count++
				} else if v1.Status < 500 {
					total.ajax_4xx_count++
				} else if v1.Status < 600 {
					total.ajax_5xx_count++
				} else {
					fmt.Println("AJAX StatusCode Error(>500):", v1.Status)
				}
			}
		}
		if total.ajax_count != 0 {
			total.ajax_duration /= float32(total.ajax_count)
		}
	}
	return total
}

func (um *AJAXMap) GetUrlSums(request_host, request_path string, pCode int64) (ajaxUrlSums []AJAXCountUrlSum) {
	um.RLock()
	defer um.RUnlock()
	if um.items != nil {
		for k, v := range um.items[request_host][request_path][pCode] {
			ajaxUrlSums = append(ajaxUrlSums, AJAXCountUrlSum{k, um.Sum(v)})
		}
	}
	return ajaxUrlSums
}

func (um *AJAXMap) SendAjaxStatsTagCounter() {
	statsType := um.GetStatsType()
	switch statsType {
	case "r":
		um.SendResourceStats()
	case "p":
		um.SendPageLoadStats()
	default:
		fmt.Println("SendAjaxStatsTagCounter(): Unknown Stats Type:", statsType)
	}
}

func (um *AJAXMap) SendResourceStats() {
	if um.items != nil {
		// request_host, request_path, pCode, page_path
		for request_host, v := range um.items {
			for request_path, v1 := range v {
				for pCode := range v1 {
					pCodeSum := um.GetPcodeSums(request_host, request_path, pCode)
					fmt.Println("AJAX Resource Stats pCodeSum:", pCodeSum)
					urlSums := um.GetUrlSums(request_host, request_path, pCode)
					for _, urlSum := range urlSums {
						fmt.Println("AJAX Resource Stats urlSum:", urlSum)
					}
				}
			}
		}
		um.RemoveAll()
	}
}

func (um *AJAXMap) SendPageLoadStats() {
	if um.items != nil {
		// request_host, request_path, pCode, page_path
		for request_host, v := range um.items {
			for request_path, v1 := range v {
				for pCode := range v1 {
					pCodeSum := um.GetPcodeSums(request_host, request_path, pCode)
					fmt.Println("AJAX PageLoad Stats pCodeSum:", pCodeSum)
					urlSums := um.GetUrlSums(request_host, request_path, pCode)
					for _, urlSum := range urlSums {
						fmt.Println("AJAX PageLoad Stats urlSum:", urlSum)
					}
				}
			}
		}
		um.RemoveAll()
	}
}

func (um *AJAXMap) SendAJAXTagCountGoFunc() {
	fiveSecondsTicker := time.NewTicker(10 * time.Second)
	now := time.Now().UTC()

	// 9초, 4초
	then := time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second()+5-now.Second()%5+4, 0, time.UTC)
	diff := then.Sub(now)
	delay5s := time.NewTimer(diff)
	go func() {
		for {
			select {
			case <-delay5s.C:
				fiveSecondsTicker.Reset(5 * time.Second)
				delay5s.Stop()
				um.SendAjaxStatsTagCounter()
			case <-fiveSecondsTicker.C:
				um.SendAjaxStatsTagCounter()
			case <-um.stop:
				return
			}
		}
	}()
}

// request_host, request_path, pCode, page_path
func (um *AJAXMap) Add(request_host, request_path string, pCode int64, page_path string, ajaxCnt AJAXStatus) {
	um.Lock()
	defer um.Unlock()
	if um.items == nil {
		um.items = make(map[string]map[string]map[int64]map[string][]AJAXStatus, 10)
	}
	if um.items[request_host] == nil {
		um.items[request_host] = make(map[string]map[int64]map[string][]AJAXStatus, 10)
	}
	if um.items[request_host][request_path] == nil {
		um.items[request_host][request_path] = make(map[int64]map[string][]AJAXStatus, 10)
	}
	if um.items[request_host][request_path][pCode] == nil {
		um.items[request_host][request_path][pCode] = make(map[string][]AJAXStatus, 10)
	}
	if um.items[request_host][request_path][pCode][page_path] == nil {
		um.items[request_host][request_path][pCode][page_path] = make([]AJAXStatus, 0, 10)
	}
	um.items[request_host][request_path][pCode][request_path] = append(um.items[request_host][request_path][pCode][request_path], ajaxCnt)
}

func (um *AJAXMap) RemoveAll() {
	um.Lock()
	defer um.Unlock()
	um.items = make(map[string]map[string]map[int64]map[string][]AJAXStatus, 10)
}

func (um *AJAXMap) GetMapDump() (strDump string) {
	um.RLock()
	defer um.RUnlock()
	strDump = fmt.Sprintf("%+v", um.items)
	return strDump
}

func (um *AJAXMap) CloseMap() {
	um.RemoveAll()
	close(um.stop)
}

/*
rum_resource_all_page  -- 따로 만들어야...
rum_resource_each_page
tags: page_path	request_host	request_path	type
fields:
resource_connection_time
resource_count
resource_dns_time
resource_download_time
resource_duration
resource_ttfb_time
*/

type ResourceMap struct {
	sync.RWMutex
	// request_host request_path type pCode page_path
	items map[string]map[string]map[string]map[int64]map[string][]ResourceRespTime
	stop  chan bool
}

type ResourceRespTime struct {
	Resource_connection_time float32
	Resource_dns_time        float32
	Resource_download_time   float32
	Resource_duration        float32
	Resource_ttfb_time       float32
}

func (um *ResourceMap) Add(request_host, request_path, rs_type string, pCode int64, page_path string, rsRespTime ResourceRespTime) {
	um.Lock()
	defer um.Unlock()
	if um.items == nil {
		um.items = make(map[string]map[string]map[string]map[int64]map[string][]ResourceRespTime, 10)
	}
	if um.items[request_host] == nil {
		um.items[request_host] = make(map[string]map[string]map[int64]map[string][]ResourceRespTime, 10)
	}
	if um.items[request_host][request_path] == nil {
		um.items[request_host][request_path] = make(map[string]map[int64]map[string][]ResourceRespTime, 10)
	}
	if um.items[request_host][request_path][rs_type] == nil {
		um.items[request_host][request_path][rs_type] = make(map[int64]map[string][]ResourceRespTime, 10)
	}
	if um.items[request_host][request_path][rs_type][pCode] == nil {
		um.items[request_host][request_path][rs_type][pCode] = make(map[string][]ResourceRespTime, 10)
	}
	if um.items[request_host][request_path][rs_type][pCode][page_path] == nil {
		um.items[request_host][request_path][rs_type][pCode][page_path] = make([]ResourceRespTime, 0, 10)
	}
	um.items[request_host][request_path][rs_type][pCode][page_path] = append(um.items[request_host][request_path][rs_type][pCode][page_path], rsRespTime)
}

func (um *ResourceMap) RemoveAll() {
	um.Lock()
	defer um.Unlock()
	um.items = make(map[string]map[string]map[string]map[int64]map[string][]ResourceRespTime, 10)
}

func (um *ResourceMap) CloseMap() {
	um.RemoveAll()
	close(um.stop)
}
