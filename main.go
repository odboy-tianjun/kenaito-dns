package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"log"
	"net/http"
	"time"
)

func main() {
	initRestfulServer()
	// 注册 DNS 请求处理函数
	dns.HandleFunc(".", handleDNSRequest)
	// 设置服务器地址和协议
	server := &dns.Server{Addr: ":53", Net: "udp"}
	// 开始监听
	log.Printf("Starting DNS server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	// 将 DNS 响应标记为权威应答
	msg.Authoritative = true
	// 将 DNS 响应标记为递归可用
	// msg.RecursionAvailable = true
	// 遍历请求中的问题部分，生成相应的回答
	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA:
			handleARecord(question, msg)
			//case dns.TypeAAAA:
			//	handleAAAARecord(question, msg)
			//case dns.TypeCNAME:
			//	handleCNAMERecord(question, msg)
			//case dns.TypeMX:
			//	handleMXRecord(question, msg)
			//case dns.TypeTXT:
			//	handleTXTRecord(question, msg)
		}
	}
	// 发送响应
	w.WriteMsg(msg)
}

func initRestfulServer() {
	router := gin.Default()
	s := &http.Server{
		Addr:         ":18001",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	initRestFunc(router)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalln("restful服务监听失败")
		panic("restful服务监听失败")
	}
	log.Println("restful服务监听ing")
}
