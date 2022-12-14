package main

import (
	"fmt"
	"regexp"
	"strings"
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

type Human interface {
	Learn()
}

type Teacher struct {
	Name string
}

func (m *Teacher) Learn() {
	fmt.Println("Teacher can learn")
}

func (m *Teacher) Teach() {
	fmt.Println("Teacher can teach")
}

type Student struct {
	Name string
}

func (m *Student) Learn() {
	fmt.Println("Student can learn")
}

func Study(h Human) {
	if h != nil {
		h.Learn()

		var s *Teacher = h.(*Teacher)
		fmt.Printf("Name: %v\n", s.Name)
		s.Teach() // ERROR
	}
}

func main() {
	Study(&Teacher{Name: "John"})

	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("I'm a bool.")
		case int:
			fmt.Println("I'm an int.")
		case string:
			fmt.Println("I'm a string.")
		default:
			fmt.Printf("Don't know type %T.\n", t)
		}
	}
	whatAmI(true)
	whatAmI("a")
	whatAmI(345)
	var rs ResourceData
	whatAmI(rs)

	var a, b, c, d bool = false, false, false, false

	switch {
	case a:
		fmt.Println("a")
	case b:
		fmt.Println("b")
	case c:
		fmt.Println("c")
	case d:
		fmt.Println("d")
	default:
		fmt.Println("Unknown")
	}

	url := "http://hello.com/128/kr/"
	temp := "http://///13/ljljl/sss/d/39a/a11/122pppppppppppppp/12"
	match, _ := regexp.MatchString("/([0-9]+)/", url)
	fmt.Println("regexp match:", match)

	r := regexp.MustCompile("([0-9]+)")
	fmt.Println(r.ReplaceAllString(url, "?"))
	fmt.Println(r.FindAllString(temp, -1))

	fmt.Println("Replace numbers to ? in url:", ReplaceNumbers(url))
	fmt.Println("Replace numbers to ? in url:", ReplaceNumbers(temp))

	// now := time.Now()
	// for i := 0; i < 1000; i++ {
	// 	ReplaceNumbers(temp)
	// }
	// then := time.Now()
	// diff := then.Sub(now)
	// fmt.Println("ReplaceNumbers 1000 times:", diff)

	fmt.Println(temp, ReplaceNumbers(temp))

	fmt.Println(IsNumber(""))

	split_strings := strings.Split(temp, "/")
	fmt.Println(strings.Join(split_strings, "/"))

	fmt.Println("ReplaceNumbers():", ReplaceNumbers(""))

}

func ReplaceNumbers2(temp string) string {
	r := regexp.MustCompile("/([0-9]+)/")
	r1 := regexp.MustCompile("/([0-9]+)$")
	return r1.ReplaceAllString(r.ReplaceAllString(temp, "/?/"), "/?")
}

func ReplaceNumbers(temp string) string {
	if temp == "" {
		return ""
	}
	var nums []string
	for _, v := range strings.Split(temp, "/") {
		if v != "" && IsNumber(v) {
			nums = append(nums, "?")
		} else {
			nums = append(nums, v)
		}
	}
	return strings.Join(nums, "/")
}

func IsNumber(num string) bool {
	if num == "" {
		return false
	}
	for _, v := range num {
		if (v > '9') || (v < '0') {
			return false
		}
	}
	return true
}
