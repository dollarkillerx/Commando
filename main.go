/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 19:41
 */
package main

import (
	"commando/define"
	"commando/engine"
	"commando/utils"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Commando"
	app.Copyright = "For information security reconnaissance services(该程序仅用于信息防御技术  请勿用于其他用途  一切的法律责任与软件作者无关  用户使用该软件 默认自动同意以上条款！！！)"
	app.Version = "v0.1"

	taskConfig := define.TaskConfig{}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "t,target",
			Usage:       "Attack target",
			Required:    true,
			Destination: &taskConfig.Tag,
		},
		cli.IntFlag{
			Name:        "c,concurrency",
			Usage:       "Test the number of concurrency",
			Required:    true,
			Value:       100,
			Destination: &taskConfig.Concurrency,
		},
		cli.IntFlag{
			Name:        "p,pools",
			Usage:       "Number of task pools",
			Required:    true,
			Value:       100,
			Destination: &taskConfig.Pool,
		},
	}
	// 执行逻辑入口
	app.Action = func(ctx *cli.Context) error {
		engine.Engine.Run(&taskConfig)
		return nil
	}

	// 容错
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-h")
	}
	go utils.Signal()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
