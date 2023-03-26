package weatherclient

import (
	"minoutil/netutil"
)

// https://openweathermap.org/

// 根据城市查询天气预报 https://openweathermap.org/api/geocoding-api

// GetWeatherInfo 根据经纬度获取天气预报 https://openweathermap.org/api/one-call-api
func GetWeatherInfo(lat, lon string) (string, error) {
	url := "https://api.openweathermap.org/data/2.5/onecall"
	param := make(map[string]string)
	param["lat"] = lat
	param["lon"] = lon
	param["exclude"] = "minutely"
	param["appid"] = "e1bfa8970f419e0e13f4ce62c5657b5e"
	param["units"] = "metric"
	param["lang"] = "zh_cn"
	b, err := netutil.Get(url, param)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
