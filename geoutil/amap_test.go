package geoutil

import (
	"fmt"
	"testing"
)

func TestAMap(t *testing.T) {
	// 初始化高德地图
	amap := AmapClient{Key: "3a41ea13e5f5cbc90254568f10e9cc74"}
	// 根据经纬度获取地址
	address, err := amap.GetAddressByLocation(116.481488, 39.990464)
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
	address, err = amap.GetAddressByLocation(113.4022203113, 23.1378010917)
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}

func TestGetDistance(t *testing.T) {
	lng1 := 113.4022203113
	lat1 := 23.1378010917
	lng2 := 113.5826193044
	lat2 := 22.1191433172
	fmt.Println(GetDistance(lng1, lat1, lng2, lat2))
}
