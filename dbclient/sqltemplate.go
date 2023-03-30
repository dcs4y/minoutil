package dbclient

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

var sqlTemplate *template.Template

func LoadTemplate(path string) {
	tpl := template.New("sqlMapper")
	tpl.Funcs(template.FuncMap{
		"safe": func(p interface{}) interface{} {
			switch t := p.(type) {
			case string:
				// 字符串防sql注入
				t = strings.ReplaceAll(t, "'", "''")
				p = "'" + t + "'"
			}
			return p
		},
	})
	tpl, err := tpl.ParseGlob(path)
	if err != nil {
		log.Println(err)
	}
	sqlTemplate = tpl
}

// BuildSql 通过模板构建sql语句
func BuildSql(templateName string, param interface{}) (sql string, err error) {
	b := bytes.NewBuffer([]byte{})
	err = sqlTemplate.ExecuteTemplate(b, templateName, param)
	if err == nil {
		sql = b.String()
	}
	return sql, err
}
