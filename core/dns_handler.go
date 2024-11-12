package core

/*
 * @Description  DNS解析处理入口
 * @Author  www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"github.com/miekg/dns"
	"kenaito-dns/cache"
	"kenaito-dns/config"
	"kenaito-dns/constant"
	"kenaito-dns/dao"
	"net"
	"time"
)

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	// 将 DNS 响应标记为权威应答
	msg.Authoritative = true
	// 将 DNS 响应标记为递归可用
	msg.RecursionAvailable = true
	// 遍历请求中的问题部分，生成相应的回答
	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA:
			isFound := handleARecord(question, msg)
			if !isFound {
				forwardGlobalServer(question.Name, dns.TypeA, msg)
			}
			break
		case dns.TypeAAAA:
			isFound := handleAAAARecord(question, msg)
			if !isFound {
				forwardGlobalServer(question.Name, dns.TypeAAAA, msg)
			}
			break
		case dns.TypeCNAME:
			isFound := handleCNAMERecord(question, msg)
			if !isFound {
				forwardGlobalServer(question.Name, dns.TypeCNAME, msg)
			}
			break
		case dns.TypeMX:
			isFound := handleMXRecord(question, msg)
			if !isFound {
				forwardGlobalServer(question.Name, dns.TypeMX, msg)
			}
			break
		case dns.TypeTXT:
			isFound := handleTXTRecord(question, msg)
			if !isFound {
				forwardGlobalServer(question.Name, dns.TypeTXT, msg)
			}
			break
		}
	}
	// 发送响应
	err := w.WriteMsg(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] WriteMsg Failed, %v \n", time.Now().Format(config.AppTimeFormat), err))
	}
}

// 发送DNS请求到公共DNS服务器
func forwardGlobalServer(name string, rrType uint16, msg *dns.Msg) {
	m := new(dns.Msg)
	m.Authoritative = true
	m.RecursionAvailable = true
	m.SetQuestion(name, rrType)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, config.ForwardDNServer)
	if err != nil {
		fmt.Printf("[dns]  [error]  "+time.Now().Format(config.AppTimeFormat)+" [DNS] Lookup Failed: %s\n", err.Error())
	}
	if r.Rcode != dns.RcodeSuccess {
		fmt.Printf("[dns]  [error]  " + time.Now().Format(config.AppTimeFormat) + " [DNS] Lookup Failed \n")
	}
	for _, record := range r.Answer {
		switch r := record.(type) {
		case *dns.A:
			msg.Answer = append(msg.Answer, r)
			break
		case *dns.AAAA:
			msg.Answer = append(msg.Answer, r)
			break
		case *dns.MX:
			msg.Answer = append(msg.Answer, r)
			break
		case *dns.TXT:
			msg.Answer = append(msg.Answer, r)
			break
		case *dns.CNAME:
			msg.Answer = append(msg.Answer, r)
			break
		}
	}
}

// 构建 A 记录 IPV4
func handleARecord(q dns.Question, msg *dns.Msg) bool {
	name := q.Name
	queryName := name[0 : len(name)-1]
	var records []dao.ResolveRecord
	cacheKey := fmt.Sprintf("%s-%s", queryName, constant.R_A)
	value, ok := cache.KeyResolveRecordMap.Load(cacheKey)
	if ok {
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache start")
		records = value.([]dao.ResolveRecord)
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache end")
	} else {
		//// 支持pod内置域名解析
		//similarityValue := cache.NameSet.GetSimilarityValue(queryName)
		//if similarityValue != "" {
		//	records = dao.FindResolveRecordByNameType(similarityValue, constant.R_A)
		//} else {
		records = dao.FindResolveRecordByNameType(queryName, constant.R_A)
		//}
	}
	if len(records) > 0 {
		for _, record := range records {
			fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] A记录, 主机名称: %s, 目标值: %s \n", time.Now().Format(config.AppTimeFormat), name, record.Value))
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
		}
		return true
	}
	return false
}

// 构建 AAAA 记录 IPV6
func handleAAAARecord(q dns.Question, msg *dns.Msg) bool {
	name := q.Name
	queryName := name[0 : len(name)-1]
	var records []dao.ResolveRecord
	cacheKey := fmt.Sprintf("%s-%s", queryName, constant.R_AAAA)
	value, ok := cache.KeyResolveRecordMap.Load(cacheKey)
	if ok {
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache start")
		records = value.([]dao.ResolveRecord)
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache end")
	} else {
		//// 支持pod内置域名解析
		//similarityValue := cache.NameSet.GetSimilarityValue(queryName)
		//if similarityValue != "" {
		//	records = dao.FindResolveRecordByNameType(similarityValue, constant.R_AAAA)
		//} else {
		records = dao.FindResolveRecordByNameType(queryName, constant.R_AAAA)
		//}
	}
	if len(records) > 0 {
		for _, record := range records {
			fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] AAAA记录, 主机名称: %s, 目标值: %s \n", time.Now().Format(config.AppTimeFormat), name, record.Value))
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
		}
		return true
	}
	return false
}

// 构建 CNAME 记录
func handleCNAMERecord(q dns.Question, msg *dns.Msg) bool {
	name := q.Name
	queryName := name[0 : len(name)-1]
	var records []dao.ResolveRecord
	cacheKey := fmt.Sprintf("%s-%s", queryName, constant.R_CNAME)
	value, ok := cache.KeyResolveRecordMap.Load(cacheKey)
	if ok {
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache start")
		records = value.([]dao.ResolveRecord)
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache end")
	} else {
		records = dao.FindResolveRecordByNameType(queryName, constant.R_CNAME)
	}
	if len(records) > 0 {
		for _, record := range records {
			fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] CNAME记录, 主机名称: %s, 目标值: %s \n", time.Now().Format(config.AppTimeFormat), name, record.Value))
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
		}
		return true
	}
	return false
}

// 构建 MX 记录
func handleMXRecord(q dns.Question, msg *dns.Msg) bool {
	name := q.Name
	queryName := name[0 : len(name)-1]
	var records []dao.ResolveRecord
	cacheKey := fmt.Sprintf("%s-%s", queryName, constant.R_MX)
	value, ok := cache.KeyResolveRecordMap.Load(cacheKey)
	if ok {
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache start")
		records = value.([]dao.ResolveRecord)
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache end")
	} else {
		records = dao.FindResolveRecordByNameType(queryName, constant.R_MX)
	}
	if len(records) > 0 {
		for _, record := range records {
			fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] MX记录, 主机名称: %s, 目标值: %s \n", time.Now().Format(config.AppTimeFormat), name, record.Value))
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
		}
		return true
	}
	return false
}

// 构建 TXT 记录
func handleTXTRecord(q dns.Question, msg *dns.Msg) bool {
	name := q.Name
	queryName := name[0 : len(name)-1]
	var records []dao.ResolveRecord
	cacheKey := fmt.Sprintf("%s-%s", queryName, constant.R_TXT)
	value, ok := cache.KeyResolveRecordMap.Load(cacheKey)
	if ok {
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache start")
		records = value.([]dao.ResolveRecord)
		fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Query cache end")
	} else {
		records = dao.FindResolveRecordByNameType(queryName, constant.R_TXT)
	}
	if len(records) > 0 {
		for _, record := range records {
			fmt.Println(fmt.Sprintf("[app]  [info]  %s [DNS] TXT记录, 主机名称: %s, 目标值: %s \n", time.Now().Format(config.AppTimeFormat), name, record.Value))
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
		return true
	}
	return false
}
