package ds

import (
	"fmt"
	"sync"
	"time"

	ua "github.com/mileusna/useragent"
)

type OnErrorData struct {
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
	OnError []struct {
		Message   string   `json:"message"`
		Stack     []string `json:"stack"`
		Status    int      `json:"status"`
		Timestamp int64    `json:"timestamp"`
		URL       string   `json:"url"`
		Type      string   `json:"type"`
	} `json:"onError"`
}

type OnErrorMap struct {
	sync.RWMutex
	items map[int64]map[string][]ErrorCount
	stop  chan bool
}

type ErrorCountUrlSum struct {
	url string
	sum ErrorCount
}

type ErrorCount struct {
	Error_count                   int
	Error_console_count           int
	Error_onerror_count           int
	Error_promise_rejection_count int
	Error_fetch_count             int
	Error_xhr_count               int
	Error_message_count           int
	Error_none_count              int
}

func (um *OnErrorMap) Sum(errCnts []ErrorCount) (total ErrorCount) {
	total = ErrorCount{0, 0, 0, 0, 0, 0, 0, 0}
	for _, v := range errCnts {
		total.Error_count += v.Error_count
		total.Error_console_count += v.Error_console_count
		total.Error_onerror_count += v.Error_onerror_count
		total.Error_promise_rejection_count += v.Error_promise_rejection_count
		total.Error_fetch_count += v.Error_fetch_count
		total.Error_xhr_count += v.Error_xhr_count
		total.Error_message_count += v.Error_message_count
		total.Error_none_count += v.Error_none_count
	}
	return total
}

func (um *OnErrorMap) SendOnErrorTagCountGoFunc() {
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
				fmt.Println("onErrorSend delay5s.C:", time.Now().UTC())
				for k := range um.items {
					fmt.Println("onError pCode:", k)
					fmt.Println("onError GetPcodeSums:", um.GetPcodeSums(k))
					fmt.Println("onError GetUrlSums:", um.GetUrlSums(k))
					um.Remove(k)
				}

			case <-fiveSecondsTicker.C:
				fmt.Println("onErrorSend", time.Now().UTC())
				for k := range um.items {
					fmt.Println("onError pCode:", k)
					fmt.Println("onError GetPcodeSums:", um.GetPcodeSums(k))
					fmt.Println("onError GetUrlSums:", um.GetUrlSums(k))
					um.Remove(k)
				}

			case <-um.stop:
				fmt.Println("stop: <-um.Stop", <-um.stop)
				return
			}
		}
	}()
}

func (um *OnErrorMap) CloseMap() {
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

func NewOnErrorMap() *OnErrorMap {
	onErrorMap := new(OnErrorMap)
	onErrorMap.stop = make(chan bool)
	onErrorMap.SendOnErrorTagCountGoFunc()
	return onErrorMap
}

func (um *OnErrorMap) Add(pCode int64, url string, errCnt ErrorCount) {
	um.Lock()
	defer um.Unlock()
	if um.items == nil {
		um.items = make(map[int64]map[string][]ErrorCount, 10)
	}
	if um.items[pCode] == nil {
		um.items[pCode] = make(map[string][]ErrorCount, 10)
	}
	if um.items[pCode][url] == nil {
		um.items[pCode][url] = make([]ErrorCount, 0, 10)
	}
	um.items[pCode][url] = append(um.items[pCode][url], errCnt)
}

func (um *OnErrorMap) GetUrlSums(pCode int64) (errCntUrlSum []ErrorCountUrlSum) {
	um.RLock()
	defer um.RUnlock()
	_, pexists := um.items[pCode]
	if pexists {
		for k, v := range um.items[pCode] {
			errCntUrlSum = append(errCntUrlSum, ErrorCountUrlSum{k, um.Sum(v)})
		}
	}
	return errCntUrlSum
}

func (um *OnErrorMap) GetPcodeSums(pCode int64) (sum ErrorCount) {
	um.RLock()
	defer um.RUnlock()
	sum = ErrorCount{0, 0, 0, 0, 0, 0, 0, 0}
	var sums []ErrorCount
	_, pexists := um.items[pCode]
	if pexists {
		for _, v := range um.items[pCode] {
			sums = append(sums, um.Sum(v))
		}
		sum = um.Sum(sums)
	}
	return sum
}

func (um *OnErrorMap) Remove(pCode int64) {
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

type PlusOnErrorMap struct {
	Type string // device, os, browser
	sync.RWMutex
	items map[string]map[int64]map[string][]ErrorCount
	stop  chan bool
}

func ParseUserAgent(userAgent string) (string, string, string) {
	uaparse := ua.Parse(userAgent)
	// fmt.Printf("%+v\n", uaparse)
	var strDevice string
	switch {
	case uaparse.Desktop:
		strDevice = "Desktop"
	case uaparse.Bot:
		strDevice = "Bot"
	case uaparse.Mobile:
		strDevice = "Mobile"
	case uaparse.Tablet:
		strDevice = "Tablet"
	default:
		strDevice = "Unknown"
	}
	return uaparse.Name, uaparse.OS, strDevice
}

func (um *PlusOnErrorMap) GetStatsType() string {
	return um.Type
}

func NewPlusOnErrorMap(Type string) *PlusOnErrorMap {
	plusOnErrorMap := new(PlusOnErrorMap)
	plusOnErrorMap.Type = Type
	plusOnErrorMap.stop = make(chan bool)
	// plusOnErrorMap.SendPlusOnErrorTagCountGoFunc()
	return plusOnErrorMap
}

func (um *PlusOnErrorMap) Add(plus string, pCode int64, url string, errCnt ErrorCount) {
	um.Lock()
	defer um.Unlock()
	if um.items == nil {
		um.items = make(map[string]map[int64]map[string][]ErrorCount, 10)
	}
	if um.items[plus] == nil {
		um.items[plus] = make(map[int64]map[string][]ErrorCount, 10)
	}
	if um.items[plus][pCode] == nil {
		um.items[plus][pCode] = make(map[string][]ErrorCount, 10)
	}
	if um.items[plus][pCode][url] == nil {
		um.items[plus][pCode][url] = make([]ErrorCount, 0, 10)
	}
	um.items[plus][pCode][url] = append(um.items[plus][pCode][url], errCnt)
}

func (um *PlusOnErrorMap) GetUrlSums(plus string, pCode int64) (errCntUrlSum []ErrorCountUrlSum) {
	um.RLock()
	defer um.RUnlock()
	// urlavgs := []UrlAvg{}
	if um.items != nil {
		for k, v := range um.items[plus][pCode] {
			errCntUrlSum = append(errCntUrlSum, ErrorCountUrlSum{k, um.Sum(v)})
		}
	}
	return errCntUrlSum
}

func (um *PlusOnErrorMap) GetPcodeSums(plus string, pCode int64) (sum ErrorCount) {
	um.RLock()
	defer um.RUnlock()
	sum = ErrorCount{0, 0, 0, 0, 0, 0, 0, 0}
	var sums []ErrorCount

	if um.items != nil {
		for _, v := range um.items[plus][pCode] {
			sums = append(sums, um.Sum(v))
		}
		sum = um.Sum(sums)
	}
	return sum
}

func (um *PlusOnErrorMap) Sum(errCnts []ErrorCount) (total ErrorCount) {
	total = ErrorCount{0, 0, 0, 0, 0, 0, 0, 0}
	for _, v := range errCnts {
		total.Error_count += v.Error_count
		total.Error_console_count += v.Error_console_count
		total.Error_onerror_count += v.Error_onerror_count
		total.Error_promise_rejection_count += v.Error_promise_rejection_count
		total.Error_fetch_count += v.Error_fetch_count
		total.Error_xhr_count += v.Error_xhr_count
		total.Error_message_count += v.Error_message_count
		total.Error_none_count += v.Error_none_count
	}
	return total
}

func (um *PlusOnErrorMap) Remove(plus string, pCode int64) {
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

func (um *PlusOnErrorMap) SendStatsTagCounter() {
	statsType := um.GetStatsType()
	switch statsType {
	case "device":
		// um.SendDeviceStats()
	case "browser":
		// um.SendBrowserStats()
	case "os":
		// um.SendOsStats()
	default:
		fmt.Println("SendStatsTagCounter(): Unknown Stats Type:", statsType)
	}
}

func (um *PlusOnErrorMap) RemoveAll() {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		for k, v := range um.items {
			for k1, v1 := range v {
				for k2 := range v1 {
					v1[k2] = nil
					delete(v1, k2)
				}
				delete(v, k1)
			}
			delete(um.items, k)
		}
	}
}

func (um *PlusOnErrorMap) CloseMap() {
	um.RemoveAll()
	close(um.stop)
}
