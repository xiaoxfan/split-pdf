package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"split-pdf/util"
	"time"
)

const (
	FileBaseUri = "http://%s:%s/images/%s"
)

var (
	IP   string
	port string
)

func main() {
	flag.StringVar(&IP, "ip", "10.151.125.185", "ip")
	flag.StringVar(&port, "p", "8888", "port")
	flag.Parse()
	if err := os.MkdirAll("./images", 0666); err != nil {
		log.Fatal(err)
	}
	r := gin.New()
	r.POST("/split-pdf", SplitPDF)
	r.Static("/images", "./images")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

type Result struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(images []string, cost int) *Result {
	return &Result{
		Code: "0000",
		Msg:  "success",
		Data: struct {
			Images []string `json:"images"`
			Cost   int      `json:"cost"`
		}{Images: images, Cost: cost},
	}
}

func Fail(msg string) *Result {
	return &Result{
		Code: "1",
		Msg:  msg,
		Data: nil,
	}
}

type Param struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	DPI  float64               `form:"dpi"`
}

func SplitPDF(c *gin.Context) {
	start := time.Now()
	p := new(Param)
	if err := c.ShouldBind(p); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	file, err := p.File.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	defer file.Close()
	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	images, err := util.Pdf2Images(fBytes, p.DPI, -1)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	names := func() []string {
		uid := uuid.New().String()
		ns := make([]string, len(images))
		for i := range ns {
			ns[i] = fmt.Sprintf("%s-%d.png", uid, i+1)
		}
		return ns
	}()
	go func() {
		for i, image := range images {
			err := ioutil.WriteFile(fmt.Sprintf("./images/%s", names[i]), image, 0666)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	imageUrl := make([]string, len(names))
	for i := range imageUrl {
		imageUrl[i] = fmt.Sprintf(FileBaseUri, IP, port, names[i])
	}
	c.JSON(http.StatusOK, Success(imageUrl, int(time.Since(start).Milliseconds())))
	return
}
