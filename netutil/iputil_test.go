package netutil

import (
	"fmt"
	"testing"
)

func TestIpToLocation(t *testing.T) {
	{
		ips := GetRemoteIp()
		fmt.Println(ips)
	}
	{
		l, err := GetIpToLocation("")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v\n", l)
	}
	{
		l, err := GetIpToLocation("24.48.0.1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v\n", l)
	}
}
