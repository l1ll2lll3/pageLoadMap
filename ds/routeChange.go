package ds

import (
	"fmt"
	"sync"
	"time"
)

type RouteChangeData struct {
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
	RouteChange []struct {
		RouterChangeTiming struct {
			IsComplete     bool   `json:"isComplete"`
			StartTimeStamp int64  `json:"startTimeStamp"`
			EndTimeStamp   int64  `json:"endTimeStamp"`
			LoadTime       int    `json:"loadTime"`
			PageLocation   string `json:"pageLocation"`
			Host           string `json:"host"`
			Path           string `json:"path"`
			Query          string `json:"query"`
			Protocol       string `json:"protocol"`
		} `json:"routerChangeTiming"`
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
	} `json:"routeChange"`
}

type RouteChangeMap struct {
	sync.RWMutex
	// completed_loading, pCode, page_path, router_change_times
	items map[bool]map[int64]map[string][]float32
	stop  chan bool
}

type RouteChangeUrlAvg struct {
	url   string
	count int
	avg   float32
}

func (um *RouteChangeMap) Sum(router_chagnge_times []float32) (total float32) {
	total = 0
	for _, v := range router_chagnge_times {
		total += v
	}
	return total
}

func (um *RouteChangeMap) Avg(durations []float32) (avg float32) {
	length := len(durations)
	avg = 0 // error
	if length != 0 {
		avg = um.Sum(durations) / float32(length)
	}
	return avg
}

func (um *RouteChangeMap) SendRouteChangeTagCountGoFunc() {
	fiveSecondsTicker := time.NewTicker(10 * time.Second)
	now := time.Now().UTC()

	// 9초, 4초
	then := time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second()+5-now.Second()%5+4, 0, time.UTC)
	// fmt.Println("then:", then)
	diff := then.Sub(now)
	// fmt.Println("diff:", diff)
	delay5s := time.NewTimer(diff)
	go func() {
		for {
			select {
			case <-delay5s.C:
				fiveSecondsTicker.Reset(5 * time.Second)
				delay5s.Stop()
				fmt.Println("RouteChange delay5s.C:", time.Now().UTC())
				if um.items != nil {
					for k, v := range um.items {
						for k1 := range v {
							cnt, pCodeAvg := um.GetPcodeAvg(k, k1)
							fmt.Println(k, cnt, pCodeAvg)
							uavgs := um.GetUrlAvgs(k, k1)
							for _, v1 := range uavgs {
								fmt.Println(k, v1)
							}
						}
					}
					um.RemoveAll()
				}

			case <-fiveSecondsTicker.C:
				fmt.Println("RouteChange:", time.Now().UTC())
				if um.items != nil {
					for k, v := range um.items {
						for k1 := range v {
							cnt, pCodeAvg := um.GetPcodeAvg(k, k1)
							fmt.Println(k, cnt, pCodeAvg)
							uavgs := um.GetUrlAvgs(k, k1)
							for _, v1 := range uavgs {
								fmt.Println(k, v1)
							}
						}
					}
					um.RemoveAll()
				}
			case <-um.stop:
				fmt.Println("stop: <-um.Stop", <-um.stop)
				return
			}
		}
	}()
}

func (um *RouteChangeMap) CloseMap() {
	um.RemoveAll()
	close(um.stop)
}

func NewRouteChangeMap() *RouteChangeMap {
	routeChangeMap := new(RouteChangeMap)
	routeChangeMap.stop = make(chan bool)
	routeChangeMap.SendRouteChangeTagCountGoFunc()
	return routeChangeMap
}

func (um *RouteChangeMap) Add(plus bool, pCode int64, url string, router_chagnge_time float32) {
	um.Lock()
	defer um.Unlock()
	if um.items == nil {
		um.items = make(map[bool]map[int64]map[string][]float32, 10)
	}
	if um.items[plus] == nil {
		um.items[plus] = make(map[int64]map[string][]float32, 10)
	}
	if um.items[plus][pCode] == nil {
		um.items[plus][pCode] = make(map[string][]float32, 10)
	}
	if um.items[plus][pCode][url] == nil {
		um.items[plus][pCode][url] = make([]float32, 0, 10)
	}
	um.items[plus][pCode][url] = append(um.items[plus][pCode][url], router_chagnge_time)
	fmt.Println(um.items)
}

func (um *RouteChangeMap) Remove(plus bool, pCode int64) {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		for k := range um.items[plus][pCode] {
			um.items[plus][pCode][k] = nil
			delete(um.items[plus][pCode], k)
		}
		delete(um.items[plus], pCode)
	}
}

func (um *RouteChangeMap) RemoveAll() {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		for k, v := range um.items {
			for k1, v1 := range v {
				for k2, _ := range v1 {
					v1[k2] = nil
					delete(v1, k2)
				}
				delete(v, k1)
			}
			delete(um.items, k)
		}
	}
}

func (um *RouteChangeMap) GetUrlAvgs(plus bool, pCode int64) (routeChangeUrlAvgs []RouteChangeUrlAvg) {
	um.RLock()
	defer um.RUnlock()
	if um.items != nil {
		for k, v := range um.items[plus][pCode] {
			routeChangeUrlAvgs = append(routeChangeUrlAvgs, RouteChangeUrlAvg{k, len(v), um.Avg(v)})
		}
	}
	return routeChangeUrlAvgs
}

func (um *RouteChangeMap) GetPcodeAvg(plus bool, pCode int64) (cnt int, avg float32) {
	um.RLock()
	defer um.RUnlock()
	cnt = 0 // error or no value
	avg = 0 // error or no value
	sum := float32(0)

	if um.items != nil {
		for _, v := range um.items[plus][pCode] {
			sum += um.Sum(v)
			cnt += len(v)
		}
	}
	if cnt != 0 {
		avg = sum / float32(cnt)
	}
	return cnt, avg
}
