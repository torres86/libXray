package xray

import (
	"fmt"

	"github.com/xtls/libxray/nodep"
)

// Ping Xray config and find the delay and country code of its outbound.
// datDir means the dir which geosite.dat and geoip.dat are in.
// configPath means the config.json file path.
// timeout means how long the http request will be cancelled if no response, in units of seconds.
// url means the website we use to test speed. "https://www.google.com" is a good choice for most cases.
// proxy means the local http/socks5 proxy, like "socks5://[::1]:1080".
func Ping(datDir string, configPath string, timeout int, url string, proxy string) (int64, error) {
	InitEnv(datDir)
	server, err := StartXray(configPath)
	if err != nil {
		return nodep.PingDelayError, err
	}

	if err := server.Start(); err != nil {
		return nodep.PingDelayError, err
	}
	defer server.Close()

	delay, err := nodep.MeasureDelay(timeout, url, proxy)
	if err != nil {
		return delay, err
	}

	return delay, nil
}

// PingTCP 使用纯TCP连接测试Xray配置的延迟
// datDir means the dir which geosite.dat and geoip.dat are in.
// configPath means the config.json file path.
// timeout means how long the tcp connection will be cancelled if no response, in units of seconds.
// host means the target host to test tcp connection. "8.8.8.8" or "google.com" are good choices.
// port means the target port to test tcp connection. 80 for HTTP, 443 for HTTPS, 53 for DNS.
// proxy means the local http/socks5 proxy, like "socks5://[::1]:1080".
func PingTCP(datDir string, configPath string, timeout int, host string, port int, proxy string) (int64, error) {
	InitEnv(datDir)
	server, err := StartXray(configPath)
	if err != nil {
		return nodep.PingDelayError, err
	}

	if err := server.Start(); err != nil {
		return nodep.PingDelayError, err
	}
	defer server.Close()

	delay, err := nodep.MeasureTCPDelay(timeout, host, port, proxy)
	if err != nil {
		return delay, err
	}

	return delay, nil
}

// Connect 测试Xray代理的连接延迟，类似Shadowrocket的connect测速
// 这个方法专门测试到Xray代理服务器的连接建立时间，不涉及数据传输
// datDir means the dir which geosite.dat and geoip.dat are in.
// configPath means the config.json file path.
// timeout means how long the connection will be cancelled if no response, in units of seconds.
// proxyAddr means the proxy server address from config, like "proxy.example.com:1080".
func Connect(datDir string, configPath string, timeout int, proxyAddr string) (int64, error) {
	// 对于Connect测速，我们直接测试到代理服务器的连接，无需启动Xray实例
	// 这样更能反映纯代理连接性能

	if len(proxyAddr) == 0 {
		return nodep.PingDelayError, fmt.Errorf("proxy address cannot be empty")
	}

	delay, err := nodep.MeasureConnectDelay(timeout, proxyAddr)
	if err != nil {
		return delay, err
	}

	return delay, nil
}
