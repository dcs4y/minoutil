package fileutil

import (
	"fmt"
	"testing"
)

func TestQREncode(t *testing.T) {
	content := "这次介绍的是压缩与解压文件，它的基本原理是创建初始的zip文件，然后循环遍历每个文件，并使用zip编写器将其添加到存档中，确保指定deflate方法以获得更好的压缩。"
	fmt.Println(QREncodeToFile(content, 400, "D:\\temp\\temp\\qr.png"))
}

func TestQRDecode(t *testing.T) {
	fmt.Println(QRDecodeFromFile("D:\\temp\\temp\\qr.png"))
}
