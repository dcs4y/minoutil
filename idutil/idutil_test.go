package idutil

import (
	"fmt"
	"testing"
)

func TestHashId_EncodeHashId(t *testing.T) {
	fmt.Println(EncodeHashId(123))
}

func TestSnowflake_GetId(t *testing.T) {
	// 生成节点实例
	for {
		fmt.Println(IdWorker.GetId())
	}
}
