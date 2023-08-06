package geoutil

import (
	"encoding/json"
	"fmt"
	"github.com/dcs4y/minoutil/v2/netutil"
)

// AmapClient 高德地图工具
type AmapClient struct {
	Key string
}

// GetAddressByLocation 通过经纬度获取地址信息
// lng: 经度,lat: 纬度
func (mapClient *AmapClient) GetAddressByLocation(lng, lat float64) (string, error) {
	url := "https://restapi.amap.com/v3/geocode/regeo"
	b, err := netutil.Get(url, map[string]string{
		"key":      mapClient.Key,
		"location": fmt.Sprintf("%f,%f", lng, lat),
	})
	if err != nil {
		return "", err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(b, &result)
	if err != nil {
		return "", err
	}
	if result["status"] != "1" {
		return "", fmt.Errorf("%s", result["info"])
	}
	var address string
	formattedAddress := result["regeocode"].(map[string]interface{})["formatted_address"]
	switch formattedAddress.(type) {
	case string:
		address = formattedAddress.(string)
		streetNumber := result["regeocode"].(map[string]interface{})["addressComponent"].(map[string]interface{})["streetNumber"].(map[string]interface{})
		switch streetNumber["street"].(type) {
		case string:
			address = fmt.Sprintf("%s%s%s", address, streetNumber["street"], streetNumber["number"])
		}
	}
	return address, nil
}
