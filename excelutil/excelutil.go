package excelutil

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var HeaderStyle = excelize.Style{
	Border: []excelize.Border{
		{Type: "left", Color: "#000000", Style: 1},
		{Type: "top", Color: "#000000", Style: 1},
		{Type: "bottom", Color: "#000000", Style: 1},
		{Type: "right", Color: "#000000", Style: 1},
	},
	Font:      &excelize.Font{Bold: true, Family: "微软雅黑"},
	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
}

var DataStyle = excelize.Style{
	Border: []excelize.Border{
		{Type: "left", Color: "#000000", Style: 1},
		{Type: "top", Color: "#000000", Style: 1},
		{Type: "bottom", Color: "#000000", Style: 1},
		{Type: "right", Color: "#000000", Style: 1},
	},
	Font:      &excelize.Font{Bold: false, Family: "微软雅黑"},
	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
}

var HeaderStyle17 = excelize.Style{
	Font:      &excelize.Font{Bold: true, Family: "微软雅黑", Size: 17},
	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
}

var InfoStyle = excelize.Style{
	Font:      &excelize.Font{Bold: true, Family: "微软雅黑", Size: 10},
	Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
}

var DataStyleBold = excelize.Style{
	Border: []excelize.Border{
		{Type: "left", Color: "#000000", Style: 1},
		{Type: "top", Color: "#000000", Style: 1},
		{Type: "bottom", Color: "#000000", Style: 1},
		{Type: "right", Color: "#000000", Style: 1},
	},
	Font:      &excelize.Font{Bold: true, Family: "微软雅黑"},
	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
}

// WriteExcelSampleList
// 简单的表格数据生成。可重复调用添加一个Sheet到同一个文件。文件中包含多Sheet时，参数sheetName必填。
// 返回函数可追加数据到同一Sheet。
func WriteExcelSampleList(f *excelize.File, sheetName string, colWidth map[interface{}]float64, titleList []string, dataList [][]string) func(dataList [][]string) {
	// 重置Sheet名称
	if sheetName != "" {
		if f.GetSheetName(0) == "Sheet1" {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			f.NewSheet(sheetName)
		}
	} else {
		sheetName = f.GetSheetName(0)
	}
	// 设置标题行样式
	styleId, err := f.NewStyle(&HeaderStyle)
	if err != nil {
		fmt.Println(err)
	}
	// 当前sheet的行索引
	rowIndex := 1
	maxColName := getColName(len(titleList))
	f.SetCellStyle(sheetName, "A1", maxColName+"1", styleId)
	// 设置标题行高
	f.SetRowHeight(sheetName, rowIndex, 23)
	// 设置默认列宽
	f.SetColWidth(sheetName, "A", maxColName, 20)
	// 设置个别列宽
	if colWidth != nil && len(colWidth) > 0 {
		for k, v := range colWidth {
			switch c := k.(type) {
			case int:
				col, err := excelize.ColumnNumberToName(c)
				if err == nil {
					f.SetColWidth(sheetName, col, col, v)
				}
			case string:
				f.SetColWidth(sheetName, c, c, v)
			}
		}
	}
	// 填充标题行
	var colNameList []string
	for i, title := range titleList {
		colNameList = append(colNameList, getColName(i+1))
		f.SetCellValue(sheetName, colNameList[i]+"1", title)
	}

	// 设置数据区域样式
	styleId, err = f.NewStyle(&DataStyle)
	if err != nil {
		fmt.Println(err)
	}
	// 填充数据
	for _, r := range dataList {
		rowIndex++
		f.SetCellStyle(sheetName, getCellName("A", rowIndex), getCellName(maxColName, rowIndex), styleId)
		f.SetRowHeight(sheetName, rowIndex, 20)
		for i, c := range r {
			f.SetCellValue(sheetName, getCellName(colNameList[i], rowIndex), c)
		}
	}

	return func(dataList [][]string) {
		// 继续填充数据
		for _, r := range dataList {
			rowIndex++
			f.SetCellStyle(sheetName, getCellName("A", rowIndex), getCellName(maxColName, rowIndex), styleId)
			f.SetRowHeight(sheetName, rowIndex, 20)
			for i, c := range r {
				f.SetCellValue(sheetName, getCellName(colNameList[i], rowIndex), c)
			}
		}
	}
}

type titleData struct {
	title     string
	col       int64
	fieldName string
	colWidth  float64
}

// 根据col从小到大排序
func sortTitleData(titleDataList []titleData) {
	length := len(titleDataList)
	iLength := length - 1
	for i := 0; i < iLength; i++ {
		for j := i + 1; j < length; j++ {
			if titleDataList[i].col > titleDataList[j].col {
				titleDataList[i], titleDataList[j] = titleDataList[j], titleDataList[i]
			}
		}
	}
}

// WriteExcelStructList 通过结构体标签解析生成列表`excel:"index:1,title:标题,width:25,format:MethodName【返回string】"`
// MethodName必须为可导出的。
func WriteExcelStructList(f *excelize.File, sheetName string, dataList interface{}) func(dataList interface{}) {
	colWidth := make(map[interface{}]float64)
	var titleList []string
	var titleDataList []titleData
	methodMap := make(map[string]string)
	oList := reflect.ValueOf(dataList)
	if oList.Len() > 0 {
		o := oList.Index(0)
		if o.Kind() == reflect.Ptr {
			o = o.Elem()
		}
		fieldNum := o.NumField()
		for i, sf := range reflect.VisibleFields(o.Type()) {
			tagInfo, ok := sf.Tag.Lookup("excel")
			if ok {
				title := titleData{
					col:       int64(i*fieldNum + fieldNum),
					fieldName: sf.Name,
				}
				tagArray := strings.Split(tagInfo, ",")
				for _, tagStr := range tagArray {
					tags := strings.Split(tagStr, ":")
					value := strings.ReplaceAll(tags[1], "'", "")
					switch tags[0] {
					case "index":
						index, err := strconv.Atoi(value)
						if err != nil {
							fmt.Println(err)
						} else {
							title.col = int64((index-1)*fieldNum + i)
						}
					case "title":
						title.title = value
					case "width":
						colWidth, err := strconv.ParseFloat(value, 64)
						if err != nil {
							fmt.Println(err)
						} else {
							title.colWidth = colWidth
						}
					case "format":
						methodMap[title.fieldName] = value
					}
				}
				titleDataList = append(titleDataList, title)
			}
		}
	}
	sortTitleData(titleDataList)
	for i, titleObject := range titleDataList {
		titleList = append(titleList, titleObject.title)
		if titleObject.colWidth > 0 {
			col, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				fmt.Println(err)
			} else {
				colWidth[col] = titleObject.colWidth
			}
		}
	}
	// 对象转数组
	returnFunc := func(dataList reflect.Value) (valueList [][]string) {
		for i := 0; i < oList.Len(); i++ {
			o := oList.Index(i)
			if o.Kind() == reflect.Ptr {
				// pointer转struct
				o = o.Elem()
			}
			var values []string
			for _, titleObject := range titleDataList {
				v := o.FieldByName(titleObject.fieldName)
				methodName, ok := methodMap[titleObject.fieldName]
				if ok {
					method := o.MethodByName(methodName)
					if method.IsValid() {
						vs := method.Call(nil)
						if len(vs) > 0 {
							values = append(values, vs[0].String())
							continue
						}
					} else {
						// 方法调用失败后，尝试调用指针的方法
						// struct转pointer
						method = o.Addr().MethodByName(methodName)
						if method.IsValid() {
							vs := method.Call(nil)
							if len(vs) > 0 {
								values = append(values, vs[0].String())
								continue
							}
						}
					}
				}
				switch v.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					values = append(values, fmt.Sprintf("%d", v.Int()))
				case reflect.Float32, reflect.Float64:
					values = append(values, fmt.Sprintf("%0.2f", v.Float()))
				case reflect.String:
					values = append(values, v.String())
				case reflect.Struct:
					values = append(values, v.String())
				default:
					values = append(values, v.String())
				}
			}
			valueList = append(valueList, values)
		}
		return
	}
	// 生成数据并生成Excel
	valueList := returnFunc(oList)
	fn := WriteExcelSampleList(f, sheetName, colWidth, titleList, valueList)
	// 注册回调函数，方便重复添加数据
	return func(dataList interface{}) {
		oList := reflect.ValueOf(dataList)
		valueList := returnFunc(oList)
		fn(valueList)
	}
}

// CellData 单元格数据结构
type CellData struct {
	CellStart  string
	CellEnd    string
	Horizontal string // 水平对齐方式[left|center|right]
	Value      string
}

// TableData 表格数据结构
type TableData struct {
	Value     string
	ValueCols int
	Data      [][]string
	SubTable  []TableData
}

// WriteExcelStatistics
// 生成统计报表
func WriteExcelStatistics(f *excelize.File, sheetName string, colWidth map[string]float64, titleData *CellData, infoDataList *[]CellData, tableDataList *[]TableData) {
	// 重置Sheet名称
	if sheetName != "" {
		if f.GetSheetName(0) == "Sheet1" {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			f.NewSheet(sheetName)
		}
	} else {
		sheetName = f.GetSheetName(0)
	}
	// 设置标题行样式
	styleId, err := f.NewStyle(&HeaderStyle17)
	if err != nil {
		fmt.Println(err)
	}
	// 当前sheet的行列索引
	titleMinCol, rowIndex, err := excelize.SplitCellName(titleData.CellStart)
	if err != nil {
		fmt.Println(err)
	}
	row0 := rowIndex - 1
	colIndex, err := excelize.ColumnNameToNumber(titleMinCol)
	if err != nil {
		fmt.Println(err)
	}
	col0 := colIndex - 1
	// 合并标题行单元格
	f.MergeCell(sheetName, titleData.CellStart, titleData.CellEnd)
	f.SetCellStyle(sheetName, titleData.CellStart, titleData.CellStart, styleId)
	f.SetCellValue(sheetName, titleData.CellStart, titleData.Value)
	// 设置默认列宽
	titleMaxCol, rowIndex, err := excelize.SplitCellName(titleData.CellEnd)
	if err != nil {
		fmt.Println(err)
	}
	// 设置标题行高
	for r := row0 + 1; r <= rowIndex; r++ {
		f.SetRowHeight(sheetName, r, 26)
	}
	maxCol, err := excelize.ColumnNameToNumber(titleMaxCol)
	if err != nil {
		fmt.Println(err)
	}
	maxCol -= col0
	f.SetColWidth(sheetName, titleMinCol, titleMaxCol, 20)
	// 设置个别列宽
	if colWidth != nil && len(colWidth) > 0 {
		for k, v := range colWidth {
			f.SetColWidth(sheetName, k, k, v)
		}
	}
	// 设置信息行默认样式
	styleId, err = f.NewStyle(&InfoStyle)
	if err != nil {
		fmt.Println(err)
	}
	// 填充信息行的数据
	if infoDataList != nil {
		for _, cellData := range *infoDataList {
			// 根据titleData的起始坐标重置
			if cellData.CellStart != "" {
				cellData.CellStart = addCell(cellData.CellStart, col0, row0)
			}
			if cellData.CellEnd != "" {
				cellData.CellEnd = addCell(cellData.CellEnd, col0, row0)
			}
			if cellData.CellEnd != "" && cellData.CellStart != cellData.CellEnd {
				f.MergeCell(sheetName, cellData.CellStart, cellData.CellEnd)
			}
			infoStyleId := styleId
			if cellData.Horizontal != "left" {
				infoStyleId, err = f.NewStyle(&excelize.Style{
					Font:      &excelize.Font{Bold: true, Family: "微软雅黑", Size: 10},
					Alignment: &excelize.Alignment{Horizontal: cellData.Horizontal, Vertical: "center", WrapText: true},
				})
				if err != nil {
					fmt.Println(err)
				}
			}
			f.SetCellStyle(sheetName, cellData.CellStart, cellData.CellStart, infoStyleId)
			_, row, err := excelize.SplitCellName(cellData.CellStart)
			if err != nil {
				fmt.Println(err)
			}
			if rowIndex < row {
				rowIndex = row
			}
			f.SetRowHeight(sheetName, row, 20)
			f.SetCellValue(sheetName, cellData.CellStart, cellData.Value)
		}
	}
	// 设置表格样式
	styleId, err = f.NewStyle(&DataStyleBold)
	if err != nil {
		fmt.Println(err)
	}
	// 填充表格的数据
	if tableDataList != nil {
		var tableRows int
		for _, tableData := range *tableDataList {
			tableRows += addSubTableData(f, sheetName, rowIndex+tableRows+1, colIndex, &tableData)
		}
		// 设置表格单元格样式
		f.SetCellStyle(sheetName, getCellName(getColName(colIndex), rowIndex+1), getCellName(getColName(colIndex+maxCol-1), rowIndex+tableRows), styleId)
		for i := 1; i <= tableRows; i++ {
			f.SetRowHeight(sheetName, rowIndex+i, 20)
		}
	}
}

func addSubTableData(f *excelize.File, sheetName string, rowIndex int, colIndex int, tableData *TableData) (totalRows int) {
	if tableData.Value != "" {
		// 设置当前单元格的值
		f.SetCellValue(sheetName, getCellName(getColName(colIndex), rowIndex), tableData.Value)
	}
	if tableData.Data != nil {
		// 填充表格数据
		for rno, r := range tableData.Data {
			for cno, c := range r {
				f.SetCellValue(sheetName, getCellName(getColName(colIndex+tableData.ValueCols+cno), rowIndex+rno), c)
			}
		}
		totalRows += len(tableData.Data)
	}
	if tableData.SubTable != nil {
		for _, subTableData := range tableData.SubTable {
			totalRows += addSubTableData(f, sheetName, rowIndex+totalRows, colIndex+tableData.ValueCols, &subTableData)
		}
	}
	if tableData.ValueCols != 0 {
		// 合并单元格
		f.MergeCell(sheetName, getCellName(getColName(colIndex), rowIndex), getCellName(getColName(colIndex+tableData.ValueCols-1), rowIndex+totalRows-1))
	}
	return
}

// 索引转列名。index >= 1。
func getColName(index int) (colName string) {
	if index <= 26 {
		colName = string(rune(64 + index))
	} else {
		colName, _ = excelize.ColumnNumberToName(index)
	}
	return
}

// 获取单元格名称
func getCellName(cell string, row int) string {
	return cell + strconv.Itoa(row)
}

// 为单元格增加行列数
func addCell(cell string, cols int, rows int) string {
	col, row, err := excelize.SplitCellName(cell)
	if err != nil {
		fmt.Println(err)
	}
	colInt, err := excelize.ColumnNameToNumber(col)
	if err != nil {
		fmt.Println(err)
	}
	return getCellName(getColName(colInt+cols), row+rows)
}

// ReadExcelFile 读取Excel文件内容。
func ReadExcelFile(fileName string, sheetName string) (rs [][]string, err error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return
	}
	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}
	rs, err = f.GetRows(sheetName)
	if err != nil {
		return
	}
	return
}

// ReadExcel 读取Excel内容。
func ReadExcel(r io.Reader, sheetName string) (rs [][]string, err error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return
	}
	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}
	rs, err = f.GetRows(sheetName)
	if err != nil {
		return
	}
	return
}

// ReadObject 将读取的Excel文件内容，转化为map或对象的数组。list为数组指针。
func ReadObject(rs [][]string, list interface{}) error {
	if len(rs) == 0 {
		return errors.New("数据集不能为空！")
	}
	titleList := rs[0]
	rs = rs[1:]
	switch l := list.(type) {
	case *[]map[string]interface{}:
		for _, r := range rs {
			data := make(map[string]interface{})
			for i, c := range r {
				data[titleList[i]] = c
			}
			*l = append(*l, data)
		}
	default:
		listReflect := reflect.ValueOf(list).Elem()
		switch listReflect.Kind() {
		case reflect.Slice, reflect.Array:
			{
				objectReflect := listReflect.Type().Elem()
				switch objectReflect.Kind() {
				case reflect.Struct:
					listReflect.Set(reflect.MakeSlice(listReflect.Type(), 0, 20))
					fields := make([]string, len(titleList))
					num := objectReflect.NumField()
					for i := 0; i < num; i++ {
						sf := objectReflect.Field(i)
						tagInfo, b := sf.Tag.Lookup("excel")
						if b {
							// 1.通过标签查询列头
							tagArray := strings.Split(tagInfo, ",")
							for _, tagStr := range tagArray {
								tags := strings.Split(tagStr, ":")
								value := strings.ReplaceAll(tags[1], "'", "")
								switch tags[0] {
								case "title":
									for titleIndex, title := range titleList {
										if title == value {
											fields[titleIndex] = sf.Name
										}
									}
								}
							}
						} else {
							// 2.通过struct属性名查找列头
							for titleIndex, title := range titleList {
								if title == sf.Name && fields[titleIndex] == "" {
									fields[titleIndex] = title
								}
							}
						}
					}
					for _, r := range rs {
						o := reflect.New(objectReflect).Elem()
						for i, c := range r {
							fieldName := fields[i]
							if fieldName != "" {
								f := o.FieldByName(fieldName)
								if f.CanSet() {
									switch f.Kind() {
									case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
										n, err := strconv.ParseInt(c, 10, parseKindNumber(fieldName, f.Kind().String(), "int"))
										if err != nil {
											fmt.Println(fieldName, err)
										} else {
											f.SetInt(n)
										}
									case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
										n, err := strconv.ParseUint(c, 10, parseKindNumber(fieldName, f.Kind().String(), "uint"))
										if err != nil {
											fmt.Println(fieldName, err)
										} else {
											f.SetUint(n)
										}
									case reflect.Float32, reflect.Float64:
										n, err := strconv.ParseFloat(c, parseKindNumber(fieldName, f.Kind().String(), "float"))
										if err != nil {
											fmt.Println(fieldName, err)
										} else {
											f.SetFloat(n)
										}
									case reflect.String:
										f.SetString(c)
									default:
										if f.Type().String() == "time.Time" {
											t, err := time.Parse(time.DateTime, c)
											if err != nil {
												fmt.Println(fieldName, err)
												t, err = time.Parse(time.DateOnly, c)
												if err != nil {
													fmt.Println(fieldName, err)
												} else {
													f.Set(reflect.ValueOf(t))
												}
											} else {
												f.Set(reflect.ValueOf(t))
											}
										} else {
											fmt.Println(fieldName, "暂不支持的类型！")
										}
									}
								}
							}
						}
						listReflect.Set(reflect.Append(listReflect, o))
					}
				default:
					return errors.New("不支持的list数据结构！")
				}
			}
		default:
			return errors.New("不支持的list类型！")
		}
	}
	return nil
}

// 解析数字
func parseKindNumber(fieldName, kind, numberType string) int {
	k := strings.TrimPrefix(kind, numberType)
	if k == "" {
		k = "32"
	}
	i, err := strconv.ParseInt(k, 10, 32)
	if err != nil {
		fmt.Println(fieldName, err)
		return 32
	} else {
		return int(i)
	}
}
