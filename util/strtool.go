package util

/*
 * @Description  工具类
 * @Author  www.odboy.cn
 * @Date  20241107
 */
import (
	"net"
	"regexp"
	"strings"
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
	_, err := net.LookupHost(domain)
	if err != nil {
		return false
	}
	return true
}
