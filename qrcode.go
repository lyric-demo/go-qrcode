package main

import (
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/LyricTian/logger"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

var (
	text    string
	logo    string
	percent int
	size    int
	out     string
)

func init() {
	flag.StringVar(&text, "t", "", "二维码内容")
	flag.StringVar(&logo, "l", "", "二维码Logo(png)")
	flag.IntVar(&percent, "p", 15, "二维码Logo的显示比例(默认15%)")
	flag.IntVar(&size, "s", 256, "二维码的大小(默认256)")
	flag.StringVar(&out, "o", "", "输出文件")
}

func main() {
	flag.Parse()

	if text == "" {
		logger.Fatalf("请指定二维码的生成内容")
	}

	if out == "" {
		logger.Fatalf("请指定输出文件")
	}

	if exists, err := checkFile(out); err != nil {
		logger.Fatalf("检查输出文件发生错误：%s", err.Error())
	} else if exists {
		logger.Fatalf("输出文件已经存在，请重新指定")
	}

	code, err := qrcode.New(text, qrcode.Highest)
	if err != nil {
		logger.Fatalf("创建二维码发生错误：%s", err.Error())
	}

	srcImage := code.Image(size)
	if logo != "" {
		logoSize := float64(size) * float64(percent) / 100

		logger.Infof("Logo:%v", logoSize)

		srcImage, err = addLogo(srcImage, logo, int(logoSize))
		if err != nil {
			logger.Fatalf("增加Logo发生错误：%s", err.Error())
		}
	}

	outAbs, err := filepath.Abs(out)
	if err != nil {
		logger.Fatalf("获取输出文件绝对路径发生错误：%s", err.Error())
	}

	os.MkdirAll(filepath.Dir(outAbs), 0777)
	outFile, err := os.Create(outAbs)
	if err != nil {
		logger.Fatalf("创建输出文件发生错误：%s", err.Error())
	}
	defer outFile.Close()

	png.Encode(outFile, srcImage)

	logger.Infof("二维码生成成功，文件路径：%s", outAbs)
}

func checkFile(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resizeLogo(logo string, size uint) (image.Image, error) {
	file, err := os.Open(logo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	img = resize.Resize(size, size, img, resize.Lanczos3)
	return img, nil
}

func addLogo(srcImage image.Image, logo string, size int) (image.Image, error) {
	logoImage, err := resizeLogo(logo, uint(size))
	if err != nil {
		return nil, err
	}

	offset := srcImage.Bounds().Max.X/2 - size/2
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			srcImage.(*image.Paletted).Set(offset+x, offset+y, logoImage.At(x, y))
		}
	}

	return srcImage, nil
}
