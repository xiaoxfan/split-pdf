package util

import (
	"bytes"
	"golang.org/x/image/tiff"
	"image/png"
)

func Tiff2PNG(src []byte) (data []byte, err error) {
	img, err := tiff.Decode(bytes.NewReader(src))
	if err != nil {
		return
	}
	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)
	data = buf.Bytes()
	return
}
