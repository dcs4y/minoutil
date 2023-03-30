package dbclient

import (
	"fmt"
	"testing"
)

func TestBuildSql(t *testing.T) {
	LoadTemplate("templates/sql/*.gohtml")
	param := map[string]interface{}{
		"Name":      "dcs'",
		"Role":      true,
		"RoleId":    1,
		"PageSize":  20,
		"PageStart": 0,
	}
	sql, err := BuildSql("query_user_list", param)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sql)
}
