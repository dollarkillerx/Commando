/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 20:38
 */
package engine

import (
	"errors"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"time"
)

type pool struct {
	pool chan *PoolItem
}

type PoolItem struct {
	Local *net.IPAddr
	Dns   *net.IPConn
	Seed  *ipv4.RawConn
}

var Pool = &pool{}

func (p *pool) InitPool(num int) {
	p.pool = make(chan *PoolItem, num)
	for i := 0; i < num; i++ {
		p.newPoolItem()
	}
}

func (p *pool) newPoolItem() {
	item := &PoolItem{}
	laddr, err := net.ResolveIPAddr("ip4", "0.0.0.0")
	if err != nil {
		log.Fatalln("Please use root permissions for Linux? ", err)
	}
	item.Local = laddr
	c, err := net.ListenIP("ip4:udp", laddr)
	if err != nil {
		log.Fatalln("Please use root permissions for Linux? ", err)
	}
	item.Dns = c
	conn, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatalln("Please use root permissions for Linux? ", err)
	}
	item.Seed = conn

	p.pool <- item
}

func (p *pool) Get() (*PoolItem, error) {
	for {
		select {
		case data := <-p.pool:
			return data, nil
		case <-time.After(time.Second):
			return nil, errors.New("time out")
		}
	}
}

func (p *pool) Put(item *PoolItem) error {
	for {
		select {
		case p.pool <- item:
			return nil
		case <-time.After(time.Second):
			return errors.New("time out")
		}
	}
}

func (p *PoolItem) Send(header *ipv4.Header, body []byte) error {
	return p.Seed.WriteTo(header, body, nil)
}
