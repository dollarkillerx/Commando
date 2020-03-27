/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 20:38
 */
package engine

import (
	"golang.org/x/net/ipv4"
	"net"
)

type pool struct {
	pool chan *poolItem
}

type poolItem struct {
	Local *net.IPAddr
	Dns   *net.IPConn
	Seed  *ipv4.RawConn
}

var Pool = &pool{pool: make(chan *poolItem, 0)}

func (p *pool) InitPool(num int) {

}
