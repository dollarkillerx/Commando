/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 20:07
 */
package engine

import (
	"commando/define"
	"log"
	"sync"
	"sync/atomic"
)

// 任务下发
func (e *engine) TaskRelease(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		loadInt32 := atomic.LoadInt32(&define.Over)
		if loadInt32 != 0 {
			close(ch)
			log.Println("Close task delivery")
			break
		} else {
			ch <- 1
		}
	}
}

// 任务执行
func (e *engine) TaskRun(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	wg2 := sync.WaitGroup{}
	task := make(chan int, e.config.Concurrency)
loop:
	for {
		select {
		case _, ok := <-ch:

			if !ok {
				wg2.Wait()
				break loop
			}
		}
	}
}
