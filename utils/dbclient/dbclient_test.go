package dbclient

import (
	"fmt"
	"game/common"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 根据表结构生成go对象
func Test_generate(t *testing.T) {
	c := generate("db", "sys_", "sys_job")
	// 保存文件
	fmt.Println(c)
	saveFile(c)
}

// 根据表结构生成go对象
func generate(packageName, prefix string, tableNames ...string) (c string) {
	c = "package " + packageName + "\n\n"
	for _, tableName := range tableNames {
		t := &table{
			TableName: tableName,
			prefix:    prefix,
		}
		// 查询表信息
		getTableInfo(t)
		// 生成go对象
		c += generateStruct(t)
	}
	return
}

func getTableInfo(t *table) {
	{
		sql := "select TABLE_NAME,TABLE_COMMENT from information_schema.TABLES where TABLE_SCHEMA = 'game' and TABLE_NAME = '%s'"
		err := DB.Raw(fmt.Sprintf(sql, t.TableName)).Scan(&t).Error
		if err != nil {
			panic(err)
		}
	}
	{
		sql := "select COLUMN_NAME,DATA_TYPE,COLUMN_COMMENT from information_schema.COLUMNS where TABLE_SCHEMA = 'game' and TABLE_NAME = '%s'"
		err := DB.Raw(fmt.Sprintf(sql, t.TableName)).Scan(&t.columns).Error
		if err != nil {
			panic(err)
		}
	}
}

func generateStruct(t *table) string {
	entryName := t.getGoName()
	var c string
	if t.Comment != "" {
		c += "// " + entryName + " " + t.Comment + "\n"
	}
	c += "type " + entryName + " struct {\n"
	for _, column := range t.columns {
		c += "\t" + column.getGoField() + "\t" + column.getGoType()
		if column.Comment != "" {
			c += " // " + strings.ReplaceAll(column.Comment, "\n", "")
		}
		c += "\n"
	}
	c += "}\n\n"
	c += "func (t " + entryName + ") TableName() string {\n"
	c += "\treturn \"" + t.TableName + "\"\n"
	c += "}\n\n"
	return c
}

func saveFile(content string) {
	f, err := os.OpenFile(filepath.Join(common.DataPath, "db", "generate.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModeType)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
	}
}

func toGoName(s string) (goName string) {
	cs := strings.Split(s, "_")
	for _, s := range cs {
		s = strings.ToLower(s)
		goName += strings.Title(s)
	}
	return
}

type column struct {
	Field   string `gorm:"column:COLUMN_NAME"`
	Type    string `gorm:"column:DATA_TYPE"`
	Comment string `gorm:"column:COLUMN_COMMENT"`
}

func (c column) getGoField() string {
	return toGoName(c.Field)
}

func (c column) getGoType() string {
	var goType string
	switch c.Type {
	case "tinyint":
		goType = "int8"
	case "smallint":
		goType = "int"
	case "int":
		goType = "int"
	case "bigint":
		goType = "uint64"
	case "char":
		goType = "string"
	case "varchar":
		goType = "string"
	case "mediumtext":
		goType = "string"
	case "text":
		goType = "string"
	case "longtext":
		goType = "string"
	case "double":
		goType = "float32"
	case "decimal":
		goType = "float64"
	case "datetime":
		goType = "time.Time"
	case "timestamp":
		goType = "time.Time"
	}
	return goType
}

type table struct {
	columns   []column
	TableName string `gorm:"column:TABLE_NAME"`
	prefix    string
	Comment   string `gorm:"column:TABLE_COMMENT"`
}

func (table table) getGoName() string {
	tableName := strings.TrimPrefix(table.TableName, table.prefix)
	return toGoName(tableName)
}
