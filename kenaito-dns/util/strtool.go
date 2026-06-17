package util

/*
 * @Description  工具类
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"golang.org/x/net/context"
)

// IsBlank 检查字符串是否空
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsValidName 判断字符串是否是有效的域名
func IsValidName(s string) bool {
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

// IsValidDomain 判断域名格式是否合法（纯格式校验，不发起网络请求）
func IsValidDomain(domain string) bool {
	if IsBlank(domain) {
		return false
	}
	domain = strings.TrimSpace(domain)
	// 去除末尾的点（FQDN 格式）
	domain = strings.TrimSuffix(domain, ".")
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	// 通配符域名特殊处理
	if strings.HasPrefix(domain, "*.") {
		domain = domain[2:]
	}
	// 校验每个标签
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}
	labelRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if !labelRegex.MatchString(label) {
			return false
		}
	}
	// TLD 至少 2 个字符
	tld := labels[len(labels)-1]
	if len(tld) < 2 {
		return false
	}
	return true
}

// IsValidDomainResolvable 判断域名是否能实际解析（需要网络）
func IsValidDomainResolvable(domain string) bool {
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
				fmt.Println("[app]  [error]  "+NowStr()+" [DNSTool] 连接到 DNS 服务器失败: ", err)
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
		return ""
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
