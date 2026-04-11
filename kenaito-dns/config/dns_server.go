package config

/*
 * @Description  DNS服务配置
 * @Author  https://www.odboy.cn
 * @Date  20241108
 */
const (
	DNSServerPort = ":53"
)

var ForwardDNSServers = []string{
	"8.8.8.8:53",   // Google DNS
	"8.8.4.4:53",   // Google DNS Backup
	"223.5.5.5:53", // AliYun DNS
	"223.6.6.6:53", // AliYun DNS Backup
	"1.1.1.1:53",   // Cloudflare DNS
	"1.0.0.1:53",   // Cloudflare DNS Backup
}
