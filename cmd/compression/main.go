package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func main() {
	dat, err := os.ReadFile("./pageLoad.dat")
	if err == nil {
		fmt.Println(string(dat))
		fmt.Println("length of dat:", len(dat))
		sdata, _ := base64.RawStdEncoding.DecodeString(string(dat)) //RawStdEncoding.DecodeString(string(dat))
		fmt.Println(sdata)
		b := bytes.NewReader(sdata)
		// fmt.Println(b)
		r, err := zlib.NewReader(b)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		io.Copy(buf, r)

		var pg PageLoadData
		byteBuf := buf.Bytes()
		json.Unmarshal(byteBuf, &pg)

		strBuff := string(byteBuf)
		fmt.Println(len(strBuff), strBuff)
		fmt.Println()
		fmt.Println()

		for k, v := range pg.PageLoad.Resource {
			fmt.Printf("%d %+v\n", k, v)
		}

		r.Close()
	} else {
		fmt.Println(err)
	}
}

// package main

// import (
// 	"bytes"
// 	"compress/zlib"
// 	"fmt"
// 	"io"
// 	"os"
// )

// func main() {
// 	buff := []byte{120, 156, 237, 148, 77, 111, 212, 48, 16, 134, 255, 139, 15, 156, 66, 119, 243, 229, 108, 34, 33, 36, 4, 130, 222, 138, 232, 13, 113, 112, 237, 201, 218, 217, 196, 19, 236, 113, 233, 182, 218, 255, 142, 67, 90, 186, 160, 106, 187, 82, 123, 169, 212, 67, 148, 120, 230, 245, 19, 59, 121, 228, 239, 55, 204, 147, 112, 116, 110, 6, 96, 77, 202, 121, 145, 220, 23, 190, 145, 24, 198, 185, 90, 151, 197, 138, 243, 154, 167, 9, 131, 75, 176, 116, 250, 145, 53, 236, 170, 236, 209, 117, 174, 144, 158, 180, 180, 44, 97, 180, 29, 35, 133, 181, 64, 82, 199, 97, 112, 125, 28, 105, 162, 209, 55, 139, 69, 231, 209, 142, 189, 144, 160, 177, 87, 224, 78, 98, 216, 72, 84, 112, 34, 113, 88, 140, 232, 201, 191, 55, 234, 93, 154, 229, 111, 198, 120, 207, 178, 25, 240, 37, 54, 34, 228, 208, 228, 57, 120, 38, 72, 199, 224, 76, 154, 75, 95, 3, 184, 109, 172, 61, 128, 61, 115, 72, 40, 241, 239, 250, 166, 197, 155, 193, 216, 53, 107, 110, 152, 3, 101, 28, 72, 154, 158, 85, 112, 130, 12, 90, 214, 188, 173, 235, 250, 246, 227, 204, 131, 93, 194, 164, 144, 26, 142, 136, 161, 181, 71, 241, 148, 245, 143, 135, 188, 239, 143, 32, 225, 47, 219, 163, 80, 143, 39, 91, 227, 60, 125, 216, 210, 17, 251, 184, 239, 230, 177, 101, 174, 227, 148, 101, 44, 59, 240, 24, 156, 132, 83, 219, 226, 4, 25, 128, 52, 198, 55, 179, 207, 159, 206, 217, 31, 6, 133, 184, 175, 108, 57, 133, 201, 137, 189, 228, 173, 72, 149, 175, 49, 87, 245, 74, 23, 22, 220, 244, 47, 174, 230, 122, 177, 233, 64, 144, 74, 235, 213, 102, 89, 176, 221, 46, 249, 95, 215, 242, 176, 174, 217, 190, 174, 215, 89, 232, 32, 112, 197, 125, 238, 215, 171, 39, 233, 26, 175, 33, 98, 159, 199, 216, 59, 216, 171, 180, 47, 75, 218, 212, 87, 88, 138, 11, 7, 40, 186, 106, 79, 218, 188, 202, 138, 203, 160, 59, 240, 235, 49, 125, 72, 90, 126, 88, 218, 252, 159, 51, 182, 42, 127, 146, 86, 94, 154, 50, 136, 240, 36, 105, 69, 127, 17, 134, 231, 81, 118, 70, 189, 10, 251, 178, 132, 173, 76, 91, 90, 114, 213, 178, 215, 84, 239, 9, 91, 173, 55, 188, 204, 7, 23, 84, 11, 109, 20, 246, 199, 111, 16, 235, 156, 235}
// 	fmt.Println(buff)

// 	b := bytes.NewReader(buff)
// 	r, err := zlib.NewReader(b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	io.Copy(os.Stdout, r)

// 	r.Close()
// }
