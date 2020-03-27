/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 20:33
 */
package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestOne(t *testing.T) {
	split := strings.Split("8.8.8.8", ".")
	log.Println(split)
	atoi1, err := strconv.Atoi(split[0])
	if err != nil {
		log.Fatalln(err)
	}
	atoi2, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatalln(err)
	}
	atoi3, err := strconv.Atoi(split[2])
	if err != nil {
		log.Fatalln(err)
	}
	atoi4, err := strconv.Atoi(split[3])
	if err != nil {
		log.Fatalln(err)
	}
	net.IPv4(byte(atoi1), byte(atoi2), byte(atoi3), byte(atoi4))

}
