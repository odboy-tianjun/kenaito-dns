package core

/*
 * @Description  DNS解析处理入口
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"kenaito-dns/cache"
	"kenaito-dns/config"
	"kenaito-dns/constant"
	"kenaito-dns/dao"
	"kenaito-dns/util"
	"net"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
)

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
	msg.RecursionAvailable = true

	if len(r.Question) == 0 {
		msg.Rcode = dns.RcodeFormatError
		w.WriteMsg(msg)
		return
	}

	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeMX, dns.TypeTXT:
			isFound := handleRecordByType(question, msg)
			if !isFound {
				forwardGlobalServer(question.Name, question.Qtype, msg)
			}
		default:
			// 不支持的类型，转发到上游 DNS
			forwardGlobalServer(question.Name, question.Qtype, msg)
		}
	}

	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] WriteMsg Failed, %v \n", util.NowStr(), err))
	}
}

// handleRecordByType 根据记录类型统一处理，避免重复代码
func handleRecordByType(q dns.Question, msg *dns.Msg) bool {
	name := q.Name
	queryName := name[0 : len(name)-1]

	// 根据 DNS Qtype 映射到常量类型
	var rrType string
	switch q.Qtype {
	case dns.TypeA:
		rrType = constant.R_A
	case dns.TypeAAAA:
		rrType = constant.R_AAAA
	case dns.TypeCNAME:
		rrType = constant.R_CNAME
	case dns.TypeMX:
		rrType = constant.R_MX
	case dns.TypeTXT:
		rrType = constant.R_TXT
	default:
		return false
	}

	cacheKey := fmt.Sprintf("%s-%s", queryName, rrType)
	var records []dao.ResolveRecord
	value, ok := cache.KeyResolveRecordMap.Load(cacheKey)
	if ok {
		records = value.([]dao.ResolveRecord)
	} else {
		records = dao.FindResolveRecordByNameType(queryName, rrType)
	}

	if len(records) == 0 {
		return false
	}

	for _, record := range records {
		switch q.Qtype {
		case dns.TypeA:
			ip := net.ParseIP(record.Value)
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    uint32(record.Ttl) * 60,
				},
				A: ip,
			}
			msg.Answer = append(msg.Answer, rr)
		case dns.TypeAAAA:
			ip := net.ParseIP(record.Value)
			rr := &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    uint32(record.Ttl) * 60,
				},
				AAAA: ip,
			}
			msg.Answer = append(msg.Answer, rr)
		case dns.TypeCNAME:
			rr := &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    uint32(record.Ttl) * 60,
				},
				Target: record.Value + ".",
			}
			msg.Answer = append(msg.Answer, rr)
		case dns.TypeMX:
			rr := &dns.MX{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeMX,
					Class:  dns.ClassINET,
					Ttl:    uint32(record.Ttl) * 60,
				},
				Preference: 10,
				Mx:         record.Value + ".",
			}
			msg.Answer = append(msg.Answer, rr)
		case dns.TypeTXT:
			rr := &dns.TXT{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    uint32(record.Ttl) * 60,
				},
				Txt: []string{record.Value},
			}
			msg.Answer = append(msg.Answer, rr)
		}
	}
	return true
}

var dnsServerIndex uint64

// forwardGlobalServer 发送DNS请求到公共DNS服务器
func forwardGlobalServer(name string, rrType uint16, msg *dns.Msg) {
	m := new(dns.Msg)
	m.Authoritative = true
	m.RecursionAvailable = true
	m.SetQuestion(name, rrType)

	dnsGlobalServers := config.ForwardDNSServers
	dnsGlobalServerCount := len(dnsGlobalServers)
	if dnsGlobalServerCount == 0 {
		fmt.Printf("[dns]  [error]  " + util.NowStr() + " [DNS] No forward DNS servers configured\n")
		return
	}

	var success bool
	var lastErr error

	for i := 0; i < dnsGlobalServerCount; i++ {
		idx := atomic.AddUint64(&dnsServerIndex, 1) % uint64(dnsGlobalServerCount)
		server := dnsGlobalServers[idx]

		c := new(dns.Client)
		c.Timeout = 2 * time.Second

		r, _, err := c.Exchange(m, server)
		if err != nil {
			lastErr = err
			fmt.Printf("[dns]  [warn]  "+util.NowStr()+" [DNS] Failed to query %s: %s\n", server, err.Error())
			continue
		}

		if r == nil {
			lastErr = fmt.Errorf("DNS server returned nil response")
			fmt.Printf("[dns]  [warn]  "+util.NowStr()+" [DNS] Server %s returned nil response\n", server)
			continue
		}

		if r.Rcode != dns.RcodeSuccess {
			lastErr = fmt.Errorf("DNS server returned error code: %d", r.Rcode)
			fmt.Printf("[dns]  [warn]  "+util.NowStr()+" [DNS] Server %s returned error code: %d\n", server, r.Rcode)
			continue
		}

		for _, record := range r.Answer {
			switch r := record.(type) {
			case *dns.A:
				msg.Answer = append(msg.Answer, r)
			case *dns.AAAA:
				msg.Answer = append(msg.Answer, r)
			case *dns.MX:
				msg.Answer = append(msg.Answer, r)
			case *dns.TXT:
				msg.Answer = append(msg.Answer, r)
			case *dns.CNAME:
				msg.Answer = append(msg.Answer, r)
			}
		}

		success = true
		fmt.Printf("[dns]  [info]  "+util.NowStr()+" [DNS] Successfully queried from %s\n", server)
		break
	}

	if !success {
		fmt.Printf("[dns]  [error]  "+util.NowStr()+" [DNS] All forward DNS servers failed. Last error: %v\n", lastErr)
	}
}
