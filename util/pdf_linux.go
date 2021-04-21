package util

import (
	"bytes"
	"errors"
	"github.com/xiaoxfan/go-fitz"
	"image/png"
	"io/ioutil"
)

const (
	defaultDPI float64 = 72
)

var (
	PageSizeErr = errors.New("page size exceed limit")
)

func Pdf2Images(src []byte, dpi float64, pageLimit int) ([][]byte, error) {
	if dpi <= 0 {
		dpi = defaultDPI
	}
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return nil, err
	}
	defer doc.Close()
	if pageLimit > 0 && doc.NumPage() > pageLimit {
		return nil, PageSizeErr
	}
	ret := make([][]byte, doc.NumPage())
	buf := &bytes.Buffer{}
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.ImageDPI(n, dpi)
		if err != nil {
			return nil, err
		}
		if err = png.Encode(buf, img); err != nil {
			return nil, err
		}
		ret[n], err = ioutil.ReadAll(buf)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
