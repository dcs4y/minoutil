package main

import (
	"fmt"
	"game/common"
	"game/timer"
	"game/utils/logutil"
	"game/web"
	"os"
)

var banner = `
			GGGGGGGGGGGGG                AAA                MMMMMMMM               MMMMMMMM EEEEEEEEEEEEEEEEEEEEEE
		 GGG::::::::::::G               A:::A               M:::::::M             M:::::::M E::::::::::::::::::::E
	   GG:::::::::::::::G              A:::::A              M::::::::M           M::::::::M E::::::::::::::::::::E
	  G:::::GGGGGGGG::::G             A:::::::A             M:::::::::M         M:::::::::M EE::::::EEEEEEEEE::::E
	 G:::::G       GGGGGG            A:::::::::A            M::::::::::M       M::::::::::M   E:::::E       EEEEEE
	G:::::G                         A:::::A:::::A           M:::::::::::M     M:::::::::::M   E:::::E             
	G:::::G                        A:::::A A:::::A          M:::::::M::::M   M::::M:::::::M   E::::::EEEEEEEEEE   
	G:::::G    GGGGGGGGGG         A:::::A   A:::::A         M::::::M M::::M M::::M M::::::M   E:::::::::::::::E   
	G:::::G    G::::::::G        A:::::A     A:::::A        M::::::M  M::::M::::M  M::::::M   E:::::::::::::::E   
	G:::::G    GGGGG::::G       A:::::AAAAAAAAA:::::A       M::::::M   M:::::::M   M::::::M   E::::::EEEEEEEEEE   
	G:::::G        G::::G      A:::::::::::::::::::::A      M::::::M    M:::::M    M::::::M   E:::::E             
	 G:::::G       G::::G     A:::::AAAAAAAAAAAAA:::::A     M::::::M     MMMMM     M::::::M   E:::::E       EEEEEE
	  G:::::GGGGGGGG::::G    A:::::A             A:::::A    M::::::M               M::::::M EE::::::EEEEEEEE:::::E
	   GG:::::::::::::::G   A:::::A               A:::::A   M::::::M               M::::::M E::::::::::::::::::::E
		 GGG::::::GGG:::G  A:::::A                 A:::::A  M::::::M               M::::::M E::::::::::::::::::::E
			GGGGGG   GGGG AAAAAAA                   AAAAAAA MMMMMMMM               MMMMMMMM EEEEEEEEEEEEEEEEEEEEEE

`

var log = logutil.Log

// 程度主入口
// 增加日志配置
func main() {
	fmt.Print(banner)
	log.Info("GAME主程序启动：", os.Args)
	log.Info("当前工作目录：", common.RootPath)
	fmt.Println("配置信息：")
	fmt.Printf("%#v\n", common.WS.ServerConfig)
	fmt.Printf("%#v\n", common.WS.WebConfig)
	fmt.Printf("%#v\n", common.WS.DBConfig)
	fmt.Printf("%#v\n", common.WS.RedisConfig)
	fmt.Printf("%#v\n", common.WS.EmailConfig)
	log.Info("操作系统：", os.Getenv("OS"))
	go timer.Start()
	web.Start()
}
