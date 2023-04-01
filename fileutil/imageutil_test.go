package fileutil

import (
	"fmt"
	"os"
	"testing"
)

func Test_image(t *testing.T) {
	path := "D:\\temp\\temp\\qr.png"
	// 支持的格式：png,jpeg,gif
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	topImage, formatName, err := ReadImage(f) // 打开图片
	if err != nil {
		panic(err)
	}
	fmt.Println(formatName)

	// 缩放
	//topImage = ResizeImage(topImage, uint(topImage.Bounds().Dx()*2), uint(topImage.Bounds().Dy()*4))

	// 裁剪
	//topImage = SubImage(topImage, 10, 50, topImage.Bounds().Max.X, topImage.Bounds().Max.Y)

	// 新建背景图片
	backgroundImage := CreateBackgroundImage(400, 450, 0xffff)

	textImage, err := CreateTextImage(400, 25, "这两篇讲解了gg库中关于文字绘制相关的内容", "C:\\Windows\\Fonts\\simsun.ttc", 18)
	if err != nil {
		fmt.Println(err)
	}
	textImage2, err := CreateTextImage(400, 25, "这两篇讲解了gg库中关于文字绘制相关的内容", "C:\\Windows\\Fonts\\simsun.ttc", 18)

	// 合成图片
	CompositeImage(backgroundImage, topImage, textImage, textImage2)

	outFile, err := os.Create("D:\\temp\\temp\\qr_1.png")
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()
	// 保存图片
	err = WriteImage(formatName, outFile, backgroundImage)
	if err != nil {
		fmt.Println(err)
	}
}
