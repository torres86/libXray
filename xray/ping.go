package xray

import (
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
// 这个方法通过代理协议测试到目标服务器的连接建立时间，不涉及数据传输
// datDir means the dir which geosite.dat and geoip.dat are in.
// configPath means the config.json file path.
// timeout means how long the connection will be cancelled if no response, in units of seconds.
// targetHost means the target host to test proxy connection, like "google.com" or "8.8.8.8".
// targetPort means the target port to test proxy connection, like 80 or 443.
// proxy means the local proxy address, like "socks5://127.0.0.1:1080".
func Connect(datDir string, configPath string, timeout int, targetHost string, targetPort int, proxy string) (int64, error) {
	// Connect测速需要启动Xray实例来测试真正的代理协议连接
	InitEnv(datDir)
	server, err := StartXray(configPath)
	if err != nil {
		return nodep.PingDelayError, err
	}

	if err := server.Start(); err != nil {
		return nodep.PingDelayError, err
	}
	defer server.Close()

	// 使用代理协议测试连接建立时间（只握手，不传输数据）
	delay, err := nodep.MeasureProxyConnectDelay(timeout, targetHost, targetPort, proxy)
	if err != nil {
		return delay, err
	}

	return delay, nil
}
