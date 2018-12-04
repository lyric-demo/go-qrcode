# go-qrcode

> Golang 二维码生成示例，支持logo缩放比例

## 编译并使用

```
go build -o qrcode
./qrcode -l data/owl.png -t "收" -o data/shou.jpg
```

![qrcode](https://github.com/lyric-demo/go-qrcode/blob/master/data/shou.jpg)

```
Usage of ./qrcode:
  -l string
        二维码Logo(png)
  -o string
        输出文件
  -p int
        二维码Logo的显示比例(默认15%) (default 15)
  -s int
        二维码的大小(默认256) (default 256)
  -t string
        二维码内容
```
