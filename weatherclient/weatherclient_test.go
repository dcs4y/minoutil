package weatherclient

import (
	"fmt"
	"testing"
)

func TestGetWeatherInfo(t *testing.T) {
	result, err := GetWeatherInfo("29.5689", "106.5577")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
