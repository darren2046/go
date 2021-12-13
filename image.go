package golanglibs

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
)

func resizeImg(srcPath string, dstPath string, width int, height ...int) {
	file, err := os.Open(srcPath)
	Panicerr(err)

	var img image.Image
	ftype := fileType(srcPath)
	if ftype == "jpg" {
		img, err = jpeg.Decode(file)
		Panicerr(err)
	} else if ftype == "png" {
		img, err = png.Decode(file)
		Panicerr(err)
	} else if ftype == "webp" {
		img, err = webp.Decode(file)
		Panicerr(err)
	} else {
		err = errors.New("只支持读取jpg、png或者webp格式图片")
		Panicerr(err)
	}
	file.Close()

	var h int
	if len(height) != 0 {
		h = height[0]
	} else {
		h = 0
	}

	resizeWidth := 0
	resizeHeight := 0
	if img.Bounds().Max.X > width {
		resizeWidth = width
	} else {
		resizeWidth = img.Bounds().Max.X
	}

	if h == 0 || img.Bounds().Max.Y > h {
		resizeHeight = h
	} else {
		resizeHeight = img.Bounds().Max.Y
	}

	m := resize.Resize(Uint(resizeWidth), Uint(resizeHeight), img, resize.Lanczos3)

	out, err := os.Create(dstPath)
	Panicerr(err)
	defer out.Close()

	if String(dstPath).EndsWith(".jpg") || String(dstPath).EndsWith(".jpeg") {
		err = jpeg.Encode(out, m, nil)
	} else if String(dstPath).EndsWith(".png") {
		err = png.Encode(out, m)
		// } else if strEndsWith(dstPath, ".webp") {
		// 	err = webp.Encode(out, m)
	} else {
		// err = errors.New("只支持导出jpg、png或者webp，以文件扩展名来判定。")
		err = errors.New("只支持导出jpg或者png，以文件扩展名来判定。")
		Panicerr(err)
	}

	Panicerr(err)
}
