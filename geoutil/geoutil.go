package geoutil

import (
	"math"
)

// GetDistance 计算两个地址位置之间的距离。单位为：米。
// 参数：lng1,lat1,lng2,lat2，分别为两个地址的经度和纬度。
func GetDistance(lng1, lat1, lng2, lat2 float64) float64 {
	radius := 6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}
