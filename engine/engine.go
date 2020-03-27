/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 19:50
 */
package engine

import (
	"commando/define"
	"golang.org/x/net/ipv4"
	"sync"
)

type engine struct {
	config *define.TaskConfig
	header *ipv4.Header
	body   []byte
}

var Engine = &engine{}

func (e *engine) Run(config *define.TaskConfig) {
	e.config = config
	e.initConfig() // 初始化系统
	chTask := make(chan int, 1000)
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 任务下发
	go e.TaskRelease(chTask, &wg)
	// 任务执行

	wg.Wait()
}
