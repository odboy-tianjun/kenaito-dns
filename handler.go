package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
)

// 构建 A 记录的函数 IPV4
func handleARecord(q dns.Question, msg *dns.Msg) {
	name := q.Name
	targetIp := "192.235.111.111"
	fmt.Printf("请求解析的域名：%s,解析的目标IP地址:%s\n", name, targetIp)
	ip := net.ParseIP(targetIp)
	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		A: ip,
	}
	msg.Answer = append(msg.Answer, rr)
}

//// 构建 A 记录的函数 IPV6
//func handleAAAARecord(q dns.Question, msg *dns.Msg) {
//	ip := net.ParseIP("rsdw::8888")
//	rr := &dns.AAAA{
//		Hdr: dns.RR_Header{
//			Name:   q.Name,
//			Rrtype: dns.TypeAAAA,
//			Class:  dns.ClassINET,
//			Ttl:    60,
//		},
//		AAAA: ip,
//	}
//	msg.Answer = append(msg.Answer, rr)
//}

//func handleCNAMERecord(q dns.Question, msg *dns.Msg) {
//	rr := &dns.CNAME{
//		Hdr: dns.RR_Header{
//			Name:   q.Name,
//			Rrtype: dns.TypeCNAME,
//			Class:  dns.ClassINET,
//			Ttl:    60,
//		},
//		Target: "example.com.",
//	}
//	msg.Answer = append(msg.Answer, rr)
//}
//
//func handleMXRecord(q dns.Question, msg *dns.Msg) {
//	rr := &dns.MX{
//		Hdr: dns.RR_Header{
//			Name:   q.Name,
//			Rrtype: dns.TypeMX,
//			Class:  dns.ClassINET,
//			Ttl:    60,
//		},
//		Preference: 10,
//		Mx:         "mail.example.com.",
//	}
//	msg.Answer = append(msg.Answer, rr)
//}
//
//func handleTXTRecord(q dns.Question, msg *dns.Msg) {
//	rr := &dns.TXT{
//		Hdr: dns.RR_Header{
//			Name:   q.Name,
//			Rrtype: dns.TypeTXT,
//			Class:  dns.ClassINET,
//			Ttl:    60,
//		},
//		Txt: []string{"v=spf1 include:_spf.example.com ~all"},
//	}
//	msg.Answer = append(msg.Answer, rr)
//}
