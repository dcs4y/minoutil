package fileutil

import (
	"errors"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// ResizeImage 图片缩放
func ResizeImage(source image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, source, resize.Lanczos3)
}

// SubImage 图片裁剪
func SubImage(source image.Image, x0, y0, x1, y1 int) image.Image {
	switch temp := source.(type) {
	case *image.Alpha:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.Alpha16:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.CMYK:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.Gray:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.Gray16:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.NRGBA:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.NRGBA64:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.NYCbCrA:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.Paletted:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.RGBA:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.RGBA64:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	case *image.YCbCr:
		return temp.SubImage(image.Rect(x0, y0, x1, y1))
	}
	return nil
}

// CompositeImage 图片合成
func CompositeImage(backgroundImage draw.Image, flowImages ...image.Image) {
	var maxY int
	for _, currentImage := range flowImages {
		draw.Draw(backgroundImage, backgroundImage.Bounds(), currentImage, image.Pt(0, -maxY), draw.Over)
		maxY += currentImage.Bounds().Max.Y
	}
}

// CreateBackgroundImage 创建空白的背景图片
func CreateBackgroundImage(width, height int, colorHex uint16) draw.Image {
	backgroundImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < backgroundImage.Bounds().Dx(); x++ { // 添加背景图颜色
		for y := 0; y < backgroundImage.Bounds().Dy(); y++ {
			backgroundImage.Set(x, y, color.Gray16{Y: colorHex})
		}
	}
	return backgroundImage
}

// ReadImage 打开图片。需要引入对应图片类型的包。返回图片的格式名称：png、jpeg、gif。
func ReadImage(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

// WriteImage 保存图片到输出流
func WriteImage(format string, w io.Writer, m image.Image) error {
	switch format {
	case "png":
		return png.Encode(w, m)
	case "jpeg":
		return jpeg.Encode(w, m, nil)
	case "gif":
		return gif.Encode(w, m, nil)
	}
	return errors.New("不支持的图片类型！")
}

// CreateTextImage 创建纯白背景的文字图片
func CreateTextImage(width, height int, content string, fontPath string, points float64) (image.Image, error) {
	dc := gg.NewContext(width, height)
	//dc.SetRGB(0, 1, 1)
	//dc.Clear() // 将纯色图片覆盖到原画布的方式来实现纯色背景的效果
	dc.SetRGB(0, 0, 0)
	// 设置字体
	if err := dc.LoadFontFace(fontPath, points); err != nil {
		return nil, err
	}
	// 计算文本高宽
	textWidth, textHeight := dc.MeasureString(content)
	dc.DrawString(content, (float64(width)-textWidth)/2, (float64(height)+textHeight)/2)
	return dc.Image(), nil
}
