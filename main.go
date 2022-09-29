package main

import (
	"fmt"
	"net/http"
	"pageLoadMap/ds"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
	// "golib/lang/pack"
	// "golib/lang/pack"
)

type PageLoadData struct {
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
	PageLoad struct {
		NavigationTiming struct {
			StartTimeStamp int64 `json:"startTimeStamp"`
			Data           struct {
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
				DomInteractive   int `json:"domInteractive"`
				DomContentLoaded struct {
					Duration int `json:"duration"`
					Start    int `json:"start"`
				} `json:"domContentLoaded"`
				DomComplete int `json:"domComplete"`
				DomLoad     struct {
					Duration int `json:"duration"`
					Start    int `json:"start"`
				} `json:"domLoad"`
				LoadTime     int `json:"loadTime"`
				BackendTime  int `json:"backendTime"`
				FrontendTime struct {
					Duration int `json:"duration"`
					Start    int `json:"start"`
				} `json:"frontendTime"`
				RenderTime struct {
					Duration int `json:"duration"`
					Start    int `json:"start"`
				} `json:"renderTime"`
			} `json:"data"`
			EventID string `json:"eventID"`
		} `json:"navigationTiming"`
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
		TotalDuration int `json:"totalDuration"`
	} `json:"pageLoad"`
}

type PageLoadMap struct {
	sync.RWMutex
	items map[int64][]PageLoadData
	Stop  chan bool
}

func NewPageLoadMap() *PageLoadMap {
	urlMap := new(PageLoadMap)
	urlMap.Stop = make(chan bool)
	urlMap.SendPageLoadMapTagCountGoFunc()
	return urlMap
}

func (um *PageLoadMap) SendPageLoadMapTagCountGoFunc() {
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
				// fmt.Println("delay5s.C:", time.Now())
				if um.items != nil {
					for k := range um.items {
						fmt.Println(k)
						fmt.Println(um.GetPcodeAvg(k))
						um.Remove(k)
					}
				}
			case <-fiveSecondsTicker.C:
				// fmt.Println("5s", time.Now())
				if um.items != nil {
					for k := range um.items {
						fmt.Println(k)
						fmt.Println(um.GetPcodeAvg(k))
						um.Remove(k)
					}
				}

			case <-um.Stop:
				fmt.Println("stop: <-um.Stop", <-um.Stop)
				return
			}
		}
	}()
}

func (um *PageLoadMap) CloseUrlMap() {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		for k := range um.items {
			delete(um.items, k)
		}
	}
	close(um.Stop)
}

func (um *PageLoadMap) Add(pCode int64, pgload PageLoadData) {
	um.Lock()
	defer um.Unlock()
	if um.items == nil {
		um.items = make(map[int64][]PageLoadData, 10)
	}
	um.items[pCode] = append(um.items[pCode], pgload)
}

func (um *PageLoadMap) Remove(pCode int64) {
	um.Lock()
	defer um.Unlock()
	if um.items != nil {
		delete(um.items, pCode)
	}
}

func (um *PageLoadMap) GetPcodeAvg(pCode int64) (int, int, int, int, int) {
	um.RLock()
	defer um.RUnlock()

	count := 0
	pgBackendSum := 0
	pgFrontendSum := 0
	pgDurationSum := 0
	pgLoadTimeSum := 0

	pgList, exists := um.items[pCode]
	if exists {
		for _, v := range pgList {
			pgBackendSum += v.PageLoad.NavigationTiming.Data.BackendTime
			pgFrontendSum += v.PageLoad.NavigationTiming.Data.FrontendTime.Duration
			pgDurationSum += v.PageLoad.TotalDuration
			pgLoadTimeSum += v.PageLoad.NavigationTiming.Data.LoadTime
		}
		count = len(pgList)
	}

	if count != 0 {
		return count, pgBackendSum / count, pgFrontendSum / count,
			pgDurationSum / count, pgLoadTimeSum / count
	}

	return count, pgBackendSum, pgFrontendSum, pgDurationSum, pgLoadTimeSum
}

func (um *PageLoadMap) GetDeviceAvg(pCode int64) (cnt int, avg int) {
	um.RLock()
	defer um.RUnlock()

	return 0, 0
}

var pgMap *PageLoadMap
var wvMap *ds.WebVitalsMap

var oEMap *ds.OnErrorMap
var browserOeMap *ds.PlusOnErrorMap

var rChangeMap *ds.RouteChangeMap

var ajaxRMap *ds.AJAXMap
var ajaxPMap *ds.AJAXMap

func main() {

	r := gin.Default()

	pgMap = NewPageLoadMap()
	defer pgMap.CloseUrlMap()
	wvMap = ds.NewWebVitalsMap()
	defer wvMap.CloseMap()

	oEMap = ds.NewOnErrorMap()
	defer oEMap.CloseMap()
	browserOeMap = ds.NewPlusOnErrorMap("browser")
	defer browserOeMap.CloseMap()

	rChangeMap = ds.NewRouteChangeMap()
	defer rChangeMap.CloseMap()

	ajaxRMap = ds.NewAJAXMap("r")
	defer ajaxRMap.CloseMap()
	ajaxPMap = ds.NewAJAXMap("p")
	defer ajaxPMap.CloseMap()

	v1 := r.Group("/")
	{
		v1.POST("/pageLoad", pageLoad)
		v1.POST("/webVitals", webVitals)
		v1.POST("/resource", resource)
		v1.POST("/routeChange", routeChange)
		v1.POST("/onError", onError)
	}
	r.Run(":8848")
}

func pageLoad(c *gin.Context) {
	var pg PageLoadData // page performance timing data
	err := c.ShouldBindJSON(&pg)

	if err != nil {
		fmt.Println("pageLoad Error")
		c.JSON(http.StatusBadRequest, gin.H{"pageLoad Error": err.Error()})
		return
	}

	if pg.Meta.PCode == 0 {
		fmt.Println("pg.Meta.PCode == 0")
		c.JSON(http.StatusBadRequest, gin.H{"pCode == 0": "pg.Meta.PCode == 0"})
		return
	}

	for _, v := range pg.PageLoad.Resource {
		if v.Type == "fetch" || v.Type == "xhr" {
			ajaxPMap.Add(v.URLHost, v.URLPath, int64(pg.Meta.PCode), pg.Meta.Path,
				ds.AJAXStatus{Status: v.ResourceInfo.Status, Ajax_duration: float32(v.Timing.Duration)})
		}
	}

	pgMap.Add(int64(pg.Meta.PCode), pg)
	c.JSON(http.StatusOK, gin.H{
		"message": "pageLoad OK",
	})
}

func resource(c *gin.Context) {
	// time.Sleep(1000 * time.Microsecond)
	// rawData, _ := c.GetRawData()
	// strRawData := string(rawData)
	// fmt.Println(strRawData)
	var rs ds.ResourceData // resource performance timing data
	err := c.ShouldBindJSON(&rs)

	if err != nil {
		fmt.Println("resource Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"resource Error": err.Error()})
		return
	}

	if rs.Meta.PCode == 0 {
		fmt.Println("rs.Meta.PCode == 0")
		c.JSON(http.StatusBadRequest, gin.H{"pCode == 0": "rs.Meta.PCode == 0"})
		return
	}

	// var ajaxCnt ds.AJAXCount
	for _, v := range rs.Resource {
		if v.Type == "fetch" || v.Type == "xhr" {
			ajaxRMap.Add(v.URLHost, v.URLPath, int64(rs.Meta.PCode), rs.Meta.Path,
				ds.AJAXStatus{Status: v.ResourceInfo.Status, Ajax_duration: float32(v.Timing.Duration)})
		}
	}

	// fmt.Printf("resource: %+v\n", rs)
	c.String(http.StatusOK, "resource OK")
}

func routeChange(c *gin.Context) {

	// rawData, _ := c.GetRawData()
	// strRawData := string(rawData)
	// fmt.Println(strRawData)
	var rchange ds.RouteChangeData // page performance timing data
	err := c.ShouldBindJSON(&rchange)

	if err != nil {
		fmt.Println("routeChange Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"routeChange Error": err.Error()})
		return
	}

	if rchange.Meta.PCode == 0 {
		fmt.Println("rchange.Meta.PCode == 0")
		c.JSON(http.StatusBadRequest, gin.H{"pCode == 0": "rchange.Meta.PCode == 0"})
		return
	}

	// fmt.Printf("routeChange: %+v\n", rchange)

	for _, v := range rchange.RouteChange {
		rChangeMap.Add(v.RouterChangeTiming.IsComplete, int64(rchange.Meta.PCode),
			v.RouterChangeTiming.Path, float32(v.RouterChangeTiming.LoadTime))
	}
	c.String(http.StatusOK, "routeChange OK")
}

func webVitals(c *gin.Context) {

	var wv ds.WebVitalsData // page performance timing data
	err := c.ShouldBindJSON(&wv)

	if err != nil {
		fmt.Println("webVitals Error")
		c.JSON(http.StatusBadRequest, gin.H{"webVitals Error": err.Error()})
		return
	}

	if wv.Meta.PCode == 0 {
		fmt.Println("wv.Meta.PCode == 0")
		c.JSON(http.StatusBadRequest, gin.H{"pCode == 0": "wv.Meta.PCode == 0"})
		return
	}

	wvrt := wvMap.MakeWebVitalsRespTimes(wv.WebVitals.Cls, wv.WebVitals.Fid, wv.WebVitals.Lcp)
	wvMap.Add(int64(wv.Meta.PCode), wv.Meta.Path, wvrt)

	fmt.Printf("webVitals: %+v", wv)
	c.String(http.StatusOK, "webVitals OK")
}

func onError(c *gin.Context) {
	var oe ds.OnErrorData
	// var oe map[string]interface{}
	err := c.ShouldBindJSON(&oe)

	if err != nil {
		fmt.Println("onError Error")
		c.JSON(http.StatusBadRequest, gin.H{"onError Error": err.Error()})
		return
	}

	if oe.Meta.PCode == 0 {
		fmt.Println("oe.Meta.PCode == 0")
		c.JSON(http.StatusBadRequest, gin.H{"pCode == 0": "oe.Meta.PCode == 0"})
		return
	}

	errCount := ds.ErrorCount{Error_count: 0,
		Error_console_count:           0,
		Error_onerror_count:           0,
		Error_promise_rejection_count: 0,
		Error_fetch_count:             0,
		Error_xhr_count:               0,
		Error_message_count:           0}

	/*
	   ErrorType {
	     console = 'console',
	     onError = 'onError',
	     promiseRejection = 'promiseRejection',
	     fetchError = 'fetchError',
	     xhrError = 'xhrError',
	     messageError = 'messageError',
	     none = 'none',
	   }
	*/
	for _, v := range oe.OnError {
		errCount.Error_count++
		switch v.Type {
		case "console":
			errCount.Error_console_count++
		case "onError":
			errCount.Error_onerror_count++
		case "promiseRejection":
			errCount.Error_promise_rejection_count++
		case "fetchError":
			errCount.Error_fetch_count++
		case "xhrError":
			errCount.Error_xhr_count++
		case "messageError":
			errCount.Error_message_count++
		default:
			// for none & error_count
			errCount.Error_none_count++
		}
	}
	oEMap.Add(int64(oe.Meta.PCode), oe.Meta.Path, errCount)

	browser, _, _ := ds.ParseUserAgent(oe.Meta.UserAgent)
	browserOeMap.Add(browser, int64(oe.Meta.PCode), oe.Meta.Path, errCount)

	c.String(http.StatusOK, "onError OK")
}

// func mapStringTest(oe map[string]interface{}) {
// 	for k, v := range oe {
// 		switch c := v.(type) {
// 		case string:
// 			fmt.Printf("Item %q is a string, containing %q\n", k, c)
// 		case float64:
// 			fmt.Printf("Looks like item %q is a number, specifically %f\n", k, c)
// 		case map[string]interface{}:
// 			for k1, v1 := range v.(map[string]interface{}) {
// 				switch c := v1.(type) {
// 				case string:
// 					fmt.Printf("Item %q is a string, containing %q\n", k1, c)
// 				case float64:
// 					fmt.Printf("Looks like item %q is a number, specifically %f\n", k1, c)
// 				default:
// 					fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k1, c)
// 				}
// 			}
// 			fmt.Printf("Looks like item %q is a number, specifically %T\n", k, c)
// 		default:
// 			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
// 		}
// 	}
// 	fmt.Printf("onError: %+v", oe)
// }
