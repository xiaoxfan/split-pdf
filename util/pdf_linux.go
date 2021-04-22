package util

import (
	"errors"
	"github.com/gen2brain/go-fitz"
	"log"
	"sync"
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
	for n := 0; n < doc.NumPage(); n++ {
		ret[n], err = doc.ImagePNG(n, dpi)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func Pdf2Images1(src []byte, dpi float64, pageLimit int) ([][]byte, error) {
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
	wg := new(sync.WaitGroup)
	wg.Add(doc.NumPage())
	for n := 0; n < doc.NumPage(); n++ {
		go func(n int) {
			defer wg.Done()
			doc, err := fitz.NewFromMemory(src)
			if err != nil {
				log.Println(err)
				return
			}
			defer doc.Close()
			ret[n], err = doc.ImagePNG(n, dpi)
			if err != nil {
				log.Println(err)
				return
			}
		}(n)
	}
	wg.Wait()
	return ret, nil
}
