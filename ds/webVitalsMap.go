package ds

import (
	"fmt"
	"sync"
	"time"
	// "golib/lang/pack"
)

type WebVitalsData struct {
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
	WebVitals struct {
		Cls float64 `json:"CLS"`
		Fid int     `json:"FID"`
		Lcp int     `json:"LCP"`
	} `json:"webVitals"`
}

type WebVitalsMap struct {
	sync.RWMutex
	items map[int64]map[string][]WebVitalsRespTime
	stop  chan bool
}

func (um *WebVitalsMap) SendWebVitalsTagCountGoFunc() {
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
				um.SendWebVitalsStats()
			case <-fiveSecondsTicker.C:
				// fmt.Println(time.Now().UTC())
				um.SendWebVitalsStats()
			case <-um.stop:
				fmt.Println("stop: <-um.Stop", <-um.stop)
				return
			}
		}
	}()

}

func (um *WebVitalsMap) SendWebVitalsStats() {
	fmt.Println("SendWebVitalsStats-------------------------------------------")
	if um.items != nil {
		// request_host, request_path, ,rs_type, pCode, page_path
		for pCode := range um.items {
			pCodeAvg := um.GetPcodeAvg(pCode)
			fmt.Printf("SendWebVitalsStats pCodeAvg: %d %+v\n", pCode, pCodeAvg)
			urlAvgs := um.GetUrlAvgs(pCode)
			for _, urlavg := range urlAvgs {
				fmt.Printf("SendWebVitalsStats urlavg: %d %+v\n", pCode, urlavg)
			}
		}
		fmt.Println("SendResourceStats Dump:", um.GetMapDump())
		um.RemoveAll()
		fmt.Println("SendResourceStats Dump after RemoveAll():", um.GetMapDump())
	}
}

func (um *WebVitalsMap) CloseMap() {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		for k, v := range um.items {
			for k1 := range v {
				v[k1] = nil
				delete(v, k1)
			}
			delete(um.items, k)
		}
	}
	close(um.stop)
}

func NewWebVitalsMap() *WebVitalsMap {
	webVitalsMap := new(WebVitalsMap)
	webVitalsMap.stop = make(chan bool)
	webVitalsMap.SendWebVitalsTagCountGoFunc()
	return webVitalsMap
}

type WebVitalsMapAvg struct {
	url string
	avg WebVitalsTimeAvg
}

type WebVitalsRespTime struct {
	cls float32
	fid float32
	lcp float32
}

type WebVitalsTimeAvg struct {
	cls     float32
	cls_cnt int
	fid     float32
	fid_cnt int
	lcp     float32
	lcp_cnt int
}

func (um *WebVitalsMap) MakeWebVitalsRespTimes(cls float64, fid, lcp int) (aaa WebVitalsRespTime) {
	aaa.cls = float32(cls)
	aaa.fid = float32(fid)
	aaa.lcp = float32(lcp)
	return aaa
}

func (um *WebVitalsMap) Add(pCode int64, url string, wVRespTime WebVitalsRespTime) {
	um.Lock()
	defer um.Unlock()

	if um.items == nil {
		// fmt.Println("um.items == nil")
		um.items = make(map[int64]map[string][]WebVitalsRespTime, 10)
	}
	if um.items[pCode] == nil {
		// fmt.Println("um.items[pCode] == nil")
		um.items[pCode] = make(map[string][]WebVitalsRespTime, 10)
	}
	if um.items[pCode][url] == nil {
		// fmt.Println("um.items[pCode][url] == nil")
		um.items[pCode][url] = make([]WebVitalsRespTime, 0, 10)
	}
	um.items[pCode][url] = append(um.items[pCode][url], wVRespTime)
}

func (um *WebVitalsMap) GetPcodeAvg(pCode int64) WebVitalsTimeAvg {
	um.RLock()
	defer um.RUnlock()
	avg := WebVitalsTimeAvg{0, 0, 0, 0, 0, 0}

	if um.items != nil {
		for _, wvTimes := range um.items[pCode] {
			for _, wvTime := range wvTimes {
				if wvTime.cls != -1 {
					avg.cls_cnt++
					avg.cls += wvTime.cls
				}
				if wvTime.fid != -1 {
					avg.fid_cnt++
					avg.fid += wvTime.fid
				}
				if wvTime.lcp != -1 {
					avg.lcp_cnt++
					avg.lcp += wvTime.lcp
				}
			}
		}
		if avg.cls != 0 {
			avg.cls /= float32(avg.cls_cnt)
		}
		if avg.fid != 0 {
			avg.fid /= float32(avg.fid_cnt)
		}
		if avg.lcp != 0 {
			avg.lcp /= float32(avg.lcp_cnt)
		}
	}

	return avg
}

func (um *WebVitalsMap) GetUrlAvgs(pCode int64) (avgs []WebVitalsMapAvg) {
	um.RLock()
	defer um.RUnlock()

	if um.items != nil {
		for url, wvTimes := range um.items[pCode] {
			avg := WebVitalsTimeAvg{0, 0, 0, 0, 0, 0}
			for _, wvTime := range wvTimes {
				if wvTime.cls != -1 {
					avg.cls_cnt++
					avg.cls += wvTime.cls
				}
				if wvTime.fid != -1 {
					avg.fid_cnt++
					avg.fid += wvTime.fid
				}
				if wvTime.lcp != -1 {
					avg.lcp_cnt++
					avg.lcp += wvTime.lcp
				}
			}
			if avg.cls != 0 {
				avg.cls /= float32(avg.cls_cnt)
			}
			if avg.fid != 0 {
				avg.fid /= float32(avg.fid_cnt)
			}
			if avg.lcp != 0 {
				avg.lcp /= float32(avg.lcp_cnt)
			}
			avgs = append(avgs, WebVitalsMapAvg{url, avg})
		}
	}
	return avgs
}

func (um *WebVitalsMap) Remove(pCode int64) {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		v, pexists := um.items[pCode]
		if pexists {
			for k := range v {
				v[k] = nil
				delete(v, k)
			}
		}
		delete(um.items, pCode)
	}
}

func (um *WebVitalsMap) RemoveAll() {
	um.Lock()
	defer um.Unlock()
	um.items = make(map[int64]map[string][]WebVitalsRespTime, 10)
}

func (um *WebVitalsMap) GetMapDump() (strDump string) {
	um.RLock()
	defer um.RUnlock()
	strDump = fmt.Sprintf("%+v", um.items)
	return strDump
}
