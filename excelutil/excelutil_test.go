package excelutil

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

// 简单表格生成示例
func TestWriteExcelSampleList(t *testing.T) {
	fileName := "D:\\usr\\local\\temp\\SampleList.xlsx"
	sheetName := "会员报表"

	colWidth := make(map[interface{}]float64)
	colWidth["B"] = 25
	colWidth["D"] = 25
	colWidth["I"] = 25
	titleList := []string{"会员姓名", "支付流水号", "收支类型", "员工卡号", "支付方式", "加油升数(升)", "金额（元）", "余额（元）", "收支时间",
		"收支时间j",
		"收支时间k",
		"收支时间l",
		"收支时间m",
		"收支时间n",
		"收支时间o",
		"收支时间p",
		"收支时间q",
		"收支时间r",
		"收支时间s",
		"收支时间t",
		"收支时间u",
		"收支时间v",
		"收支时间w",
		"收支时间x",
		"收支时间y",
		"收支时间z",
		"收支时间aa",
		"收支时间ab",
		"收支时间ac",
		"收支时间d",
		"收支时间e",
		"收支时间f",
		"收支时间g",
		"收支时间h",
		"收支时间i",
		"收支时间j",
		"收支时间k",
		"收支时间l",
		"收支时间m",
		"收支时间n",
		"收支时间o",
		"收支时间p",
		"收支时间q",
		"收支时间r",
	}
	dataList := [][]string{
		{"姓名", "2728252128003072", "充值", "0012200001452454", "现金支付", "1.29", "7.01", "600.01", "2021-09-09 13:49:48"},
		{"姓名", "2728252128003072", "充值", "0012200001452454", "现金支付", "5.29", "40.01", "55.01", "2021-09-09 13:49:48"},
		{"姓名", "2728252128003072", "充值", "0012200001452454", "现金支付", "20.29", "140.01", "2000.01", "2021-09-09 13:49:48"},
	}
	// 初始化一个Excel对象
	f := excelize.NewFile()
	addData := WriteExcelSampleList(f, sheetName, colWidth, titleList, dataList)
	for i := 0; i < 1000; i++ {
		dataList[0][0] = fmt.Sprintf("%s%d", "姓名", i*3+1)
		dataList[1][0] = fmt.Sprintf("%s%d", "姓名", i*3+2)
		dataList[2][0] = fmt.Sprintf("%s%d", "姓名", i*3+3)
		addData(dataList)
	}
	sheetName = "又是会员报表"
	WriteExcelSampleList(f, sheetName, colWidth, titleList, dataList)
	sheetName = "班组报表"
	WriteExcelSampleList(f, sheetName, colWidth, titleList, dataList)
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK==>" + f.Path)
}

type Order struct {
	OrderNo           string `excel:"title:'订单号',width:'25'"`
	TopCompanyName    string `excel:"index:'2',title:'公司名称',width:'35'"`
	MemberName        string `excel:"index:'1',title:'会员姓名'"`
	Phone             string `excel:"index:'1',title:'手机号'"`
	MemberStationName string
	Price             float32 `excel:"index:'5',title:'价格',width:'15'"`
	OrderStatus       int     `excel:"index:4,title:'状态',width:10"`
	OrderType         string
	NodeType          string
	NodeName          string
	StartDate         time.Time `excel:"title:'开始时间',width:'25',format:FormatDate"`
	EndDate           time.Time `excel:"title:'结束时间',width:'25',format:FormatDate"`
}

func (o Order) FormatDate() string {
	if o.StartDate.IsZero() {
		return ""
	}
	return o.StartDate.Format(time.DateTime)
}

// 通过结构体标签解析生成列表`excel:"index:'1',title:'标题',width:'25',format:FormatDate"`
func TestWriteExcelStructList(t *testing.T) {
	fileName := "D:\\usr\\local\\temp\\StructList.xlsx"
	sheetName := "订单报表"

	var orderList []*Order
	for i := 0; i < 100; i++ {
		iStr := strconv.Itoa(i)
		inOrder := &Order{
			OrderNo:        iStr + "#" + time.Now().Format("20060102150405"),
			TopCompanyName: "某某有限公司" + iStr,
			MemberName:     "爱尔纱" + iStr,
			Phone:          iStr + "12345678910",
			Price:          2.2*float32(i) + 1,
			OrderStatus:    i,
			StartDate:      time.Now(),
		}
		orderList = append(orderList, inOrder)
	}
	//fmt.Println(inOrderList)

	// 初始化一个Excel对象
	f := excelize.NewFile()
	addData := WriteExcelStructList(f, sheetName, orderList)
	addData(orderList)
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK==>" + f.Path)
}

// 报表生成示例
func TestWriteExcelStatistics(t *testing.T) {
	fileName := "D:\\usr\\local\\temp\\Statistics.xlsx"
	sheetName := "会员报表"

	colWidth := make(map[string]float64)
	colWidth["A"] = 25
	titleData := CellData{
		CellStart: "A1",
		CellEnd:   "G2",
		Value:     "支付方式汇总",
	}
	infoData := []CellData{
		{
			CellStart: "A3",
			CellEnd:   "E3",
			Value:     "加油站：71",
		},
		{
			CellStart: "A4",
			Value:     "班次号：202110110005",
		},
		{
			CellStart:  "B4",
			CellEnd:    "D4",
			Horizontal: "left",
			Value:      "班次开始时间：2021-10-11 15:27:09",
		},
		{
			CellStart:  "F3",
			CellEnd:    "G4",
			Horizontal: "right",
			Value:      "班次结束时间：2021-10-11 15:28:49",
		},
		{
			CellStart: "A5",
			CellEnd:   "G5",
			Value:     "说明：内容有点多！",
		},
		{
			CellStart: "D6",
			CellEnd:   "F6",
			Value:     "_____________________________________________________________________",
		},
		{
			CellStart: "A8",
			CellEnd:   "A8",
			Value:     "",
		},
	}
	var tableData = []TableData{
		{
			Data: [][]string{
				{"支付渠道", "支付方式", "数量", "销售升数(升)", "销售金额(元)", "实收金额(元)", "优惠金额(元)"},
				{"外部渠道", "现金支付", "1", "0.29", "2", "1.5", "0.5"},
				{"外部渠道", "现金支付", "2", "0.29", "2", "1.5", "0.5"},
				{"外部渠道", "现金支付", "3", "0.29", "2", "1.5", "0.5"},
			},
			SubTable: []TableData{
				{
					Value:     "小计",
					ValueCols: 2,
					Data: [][]string{
						{"1", "0.29", "2", "1.5", "0.5"},
						{"1", "0.29", "2", "1.5", "0.5"},
					},
					SubTable: []TableData{
						{
							Value:     "xxxxxxxx",
							ValueCols: 1,
							Data: [][]string{
								{"1", "0.29", "2", "1.5"},
								{"1", "0.29", "2", "1.5"},
							},
						},
						{
							Value:     "zzzzzzzzzz",
							ValueCols: 2,
							Data: [][]string{
								{"1", "0.29", "2"},
								{"1", "0.29", "2"},
							},
						},
					},
				},
			},
		},
		{
			Data: [][]string{
				{"支付渠道", "支付方式", "数量", "销售升数(升)", "销售金额(元)", "实收金额(元)", "优惠金额(元)"},
				{"外部渠道", "现金支付", "1", "0.29", "2", "1.5", "0.5"},
				{"外部渠道", "现金支付", "2", "0.29", "2", "1.5", "0.5"},
				{"外部渠道", "现金支付", "3", "0.29", "2", "1.5", "0.5"},
			},
			SubTable: []TableData{
				{
					Value:     "小计",
					ValueCols: 2,
					Data: [][]string{
						{"1", "0.29", "2", "1.5", "0.5"},
						{"1", "0.29", "2", "1.5", "0.5"},
					},
					SubTable: []TableData{
						{
							Value:     "xxxxxxxx",
							ValueCols: 1,
							Data: [][]string{
								{"1", "0.29", "2", "1.5"},
								{"1", "0.29", "2", "1.5"},
							},
						},
						{
							Value:     "zzzzzzzzzz",
							ValueCols: 2,
							Data: [][]string{
								{"1", "0.29", "2"},
								{"1", "0.29", "2"},
							},
						},
					},
				},
			},
		},
		{
			Value:     "合计",
			ValueCols: 2,
			Data: [][]string{
				{"1", "0.29", "2", "2", "0"},
				{"1", "0.29", "2", "2", "0"},
			},
		},
	}
	// 初始化一个Excel对象
	f := excelize.NewFile()
	WriteExcelStatistics(f, sheetName, colWidth, &titleData, &infoData, &tableData)
	WriteExcelStatistics(f, "第一个表格", colWidth, &titleData, &infoData, &tableData)
	colWidth["I"] = 25
	titleData.CellStart = "I1"
	titleData.CellEnd = "O2"
	titleData.Value = "横起一个表格"
	WriteExcelStatistics(f, sheetName, colWidth, &titleData, &infoData, &tableData)
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	titleData.CellStart = "A35"
	titleData.CellEnd = "G36"
	titleData.Value = "竖起一个表格"
	WriteExcelStatistics(f, sheetName, colWidth, &titleData, &infoData, &tableData)
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK==>" + f.Path)
}

// 读取表格示例
func TestReadExcel(t *testing.T) {
	fileName := "D:\\usr\\local\\temp\\SampleList.xlsx"
	rs, err := ReadExcelFile(fileName, "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rs)
	}
}

// "会员姓名":"姓名", "余额（元）":"600.01", "加油升数(升)":"1.29", "员工卡号":"0012200001452454", "支付方式":"现金支付", "支付流水号":"2728252128003072", "收支时间":"2021-09-09 13:49:48", "收支类型":"充值", "金额（元）":"7.01"
func TestReadObject(t *testing.T) {
	{
		var list []map[string]interface{}
		fileName := "D:\\usr\\local\\temp\\ReadList.xlsx"
		rs, err := ReadExcelFile(fileName, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ReadObject(rs, &list)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%#v\n", list)
		}
	}
	{
		type member struct {
			MemberName   string
			PayWay       string    `excel:"title:'支付方式'"`
			PayTime      time.Time `excel:"title:'收支时间'"`
			Balance      float32   `excel:"title:'余额（元）'"`
			EmployNumber string    `excel:"title:'员工卡号'"`
			OilV         int       `excel:"title:'加油升数(升)'"`
		}
		var list []member
		fileName := "D:\\usr\\local\\temp\\ReadList.xlsx"
		rs, err := ReadExcelFile(fileName, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ReadObject(rs, &list)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%#v\n", list)
		}
	}
}
