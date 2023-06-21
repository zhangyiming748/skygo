package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"skygo_detection/guardian/app/sys_service"

	"skygo_detection/common"
	"skygo_detection/console"
	"skygo_detection/http"
	"skygo_detection/lib/common_lib/common_const"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		fmt.Println(c.String("name"), "===")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "http",         // 命令全称, 命令简写
			Usage: "http service", // 命令详细描述
			Flags: []cli.Flag{
				// 环境变量
				cli.StringFlag{
					Name:        "env, e",
					Value:       "dev",
					Usage:       "environment type",
					Destination: &common.CliFlagEnv,
				},
				// 指定配置文件路径
				cli.StringFlag{
					Name:        "config, c",
					Value:       "./config/dev/config.tml",
					Usage:       "config file path",
					Destination: &common.CliFlagConfigPath,
				},
				// 是否开启debug模式
				cli.BoolFlag{
					Name:        "debug, d",
					Usage:       "is used debug model",
					Destination: &common_const.CliFlagDebug,
				},
			},
			Action: func(c *cli.Context) { // 命令处理函数
				// 初始化监听配置
				sys_service.InitConfigWatcher("", common.CliFlagConfigPath)
				// 启动服务
				http.Main()
			},
		},
		{
			Name:  "console",         // 命令全称, 命令简写
			Usage: "console service", // 命令详细描述
			Flags: []cli.Flag{
				// 环境变量
				cli.StringFlag{
					Name:        "env, e",
					Value:       "dev",
					Usage:       "environment type",
					Destination: &common.CliFlagEnv,
				},
				// 指定配置文件路径
				cli.StringFlag{
					Name:        "config, c",
					Value:       "",
					Usage:       "config file path",
					Destination: &common.CliFlagConfigPath,
				},
				// 指定任务类型
				cli.StringFlag{
					Name:        "type, t",
					Value:       "",
					Usage:       "task type",
					Destination: &common.CliFlagConsoleTaskType,
				},
				// 指定任务名称
				cli.StringFlag{
					Name:        "name, n",
					Value:       "",
					Usage:       "task name",
					Destination: &common.CliFlagConsoleTaskName,
				},
			},
			Action: func(c *cli.Context) { // 命令处理函数
				// 初始化监听配置
				sys_service.InitConfigWatcher("", common.CliFlagConfigPath)
				// 启动服务
				console.Main(common.CliFlagConsoleTaskType, common.CliFlagConsoleTaskName)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
