package util

/*
 * @Description  工具类
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"golang.org/x/net/context"
	"kenaito-dns/config"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

// IsBlank 检查字符串是否空
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsValidName 判断字符串是否是有效的域名
func IsValidName(s string) bool {
	// 定义域名的正则表达式
	domainRegex := `^([a-zA-Z0-9][a-zA-Z0-9\-]{1,61}[a-zA-Z0-9]\.)+[a-zA-Z0-9]{2,6}$`
	re := regexp.MustCompile(domainRegex)
	return re.MatchString(s)
}

// IsIPv4 判断字符串是否是ipv4地址
func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ".")
}

// IsIPv6 判断字符串是否是ipv6地址
func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}

// IsValidDomain 判断域名是否正常解析
func IsValidDomain(domain string) bool {
	dnsServer := getLocalIP()
	if dnsServer == "" {
		dnsServer = "223.5.5.5"
	}
	dnsServer = dnsServer + ":53"
	_, err := lookupHostWithDNS(domain, dnsServer)
	if err != nil {
		return false
	}
	return true
}

func lookupHostWithDNS(host string, dnsServer string) ([]string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			conn, err := d.DialContext(ctx, network, dnsServer)
			if err != nil {
				fmt.Println("[app]  [error]  "+time.Now().Format(config.AppTimeFormat)+" [DNSTool] 连接到 DNS 服务器失败: ", err)
				return nil, err
			}
			return conn, nil
		},
	}
	ips, err := resolver.LookupHost(context.Background(), host)
	if err != nil {
		return nil, err
	}
	return ips, nil
}

func getLocalIP() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrList {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
