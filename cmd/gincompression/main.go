package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rumagent/models"
	"time"

	// "github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.POST("/pageLoad", func(c *gin.Context) {

		// contentType := c.ContentType()
		// // for _, contentType := range contentTypes {
		// fmt.Println(contentType)
		// // }

		byteBuf := Decompression(c)
		var pg models.PageLoadData
		json.Unmarshal(byteBuf, &pg)

		strBuff := string(byteBuf)
		fmt.Println(len(strBuff), strBuff)
		fmt.Println()
		fmt.Println()

		for k, v := range pg.PageLoad.Resource {
			fmt.Printf("%d %+v\n", k, v)
		}

		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8888"); err != nil {
		log.Fatal(err)
	}
}

func Decompression(c *gin.Context) []byte {
	buf := new(bytes.Buffer)
	io.Copy(buf, c.Request.Body)
	sdata, _ := base64.RawStdEncoding.DecodeString(buf.String())
	b := bytes.NewReader(sdata)

	r, err := zlib.NewReader(b)
	if err != nil {
		log.Println("Decompression zlib.NewReader(b) error:", err)
		return nil
	}
	defer r.Close()

	buff := new(bytes.Buffer)
	io.Copy(buff, r)
	byteBuf := buff.Bytes()

	return byteBuf
}
