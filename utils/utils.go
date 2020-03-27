/**
*@program: Commando
*@description: https://github.com/dollarkillerx
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-27 19:39
 */
package utils

import (
	"encoding/binary"
	"github.com/miekg/dns"
	"log"
	"net"
)

func CheckSum(msg []byte) uint16 {
	var (
		sum    uint32
		length int = len(msg)
		index  int
	)

	for length > 1 {
		sum += uint32(msg[index])<<8 + uint32(msg[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(msg[index]) << 8
	}
	sum += (sum >> 16)
	return uint16(^sum)
}

func NewMsg() []byte {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
		Name:   "qq.com.",
	}
	b, err := msg.Pack()
	if err != nil {
		log.Fatal(err)
	}
	return b
}

type PseudoHeader struct {
	SrcIP  net.IP
	DstIP  net.IP
	Zero   byte
	Proto  byte
	Length uint16
	UDPHeader
}

func (p PseudoHeader) Bytes() []byte {
	var b = make([]byte, 20)
	b[0], b[1], b[2], b[3] = p.SrcIP[12], p.SrcIP[13], p.SrcIP[14], p.SrcIP[15]
	b[4], b[5], b[6], b[7] = p.DstIP[12], p.DstIP[13], p.DstIP[14], p.DstIP[15]
	b[8] = p.Zero
	b[9] = p.Proto
	binary.BigEndian.PutUint16(b[10:12], p.Length)
	binary.BigEndian.PutUint16(b[12:], p.SPort)
	binary.BigEndian.PutUint16(b[14:], p.DPort)
	binary.BigEndian.PutUint16(b[16:], p.Length)
	binary.BigEndian.PutUint16(b[18:], p.CheckSum)
	return b
}

type UDPHeader struct {
	SPort    uint16
	DPort    uint16
	Length   uint16
	CheckSum uint16
}

func (p UDPHeader) Bytes() []byte {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint16(b[:2], p.SPort)
	binary.BigEndian.PutUint16(b[2:4], p.DPort)
	binary.BigEndian.PutUint16(b[4:6], p.Length)
	binary.BigEndian.PutUint16(b[6:8], p.CheckSum)
	return b
}
