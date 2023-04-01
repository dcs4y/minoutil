package baseModel

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

// MarshalJSON 为 Time 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
func (t Time) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(time.DateTime))
	return []byte(output), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	value := string(b)
	if value == `""` || value == "null" || value == "0" {
		return nil
	}
	if strings.HasPrefix(value, `"`) {
		value = value[1 : len(value)-1]
		var layout string
		if len(value) <= 10 {
			layout = time.DateOnly
		} else {
			layout = time.DateTime
		}
		time, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			t.Time = time
		}
		return err
	} else {
		// 默认秒
		unixTime, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			if unixTime <= 9999999999 {
				unixTime *= 1000
			}
			// 需要毫秒
			t.Time = time.UnixMilli(unixTime)
		}
		return err
	}
}

// Value 为 Time 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 为 Time 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// String 方法
func (t Time) String() string {
	return t.Time.Format(time.DateTime)
}
