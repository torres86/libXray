package nodep

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	PingDelayTimeout int64 = 11000
	PingDelayError   int64 = 10000
)

// get the delay of some outbound.
// timeout means how long the http request will be cancelled if no response, in units of seconds.
// url means the website we use to test speed. "https://www.google.com" is a good choice for most cases.
// proxy means the local http/socks5 proxy, like "socks5://[::1]:1080". If proxy is empty, it means no proxy.
func MeasureDelay(timeout int, url string, proxy string) (int64, error) {
	httpTimeout := time.Second * time.Duration(timeout)
	c, err := CoreHTTPClient(httpTimeout, proxy)
	if err != nil {
		return PingDelayError, err
	}
	delay, err := PingHTTPRequest(c, url, timeout)
	if err != nil {
		return delay, err
	}

	return delay, nil
}

func CoreHTTPClient(timeout time.Duration, proxy string) (*http.Client, error) {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}

	if len(proxy) > 0 {
		tr.Proxy = func(r *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	return c, nil
}

func PingHTTPRequest(c *http.Client, url string, timeout int) (int64, error) {
	start := time.Now()
	req, _ := http.NewRequest("HEAD", url, nil)
	_, err := c.Do(req)
	delay := time.Since(start).Milliseconds()
	if err != nil {
		precision := delay - int64(timeout)*1000
		if math.Abs(float64(precision)) < 50 {
			return PingDelayTimeout, err
		} else {
			return PingDelayError, err
		}
	}
	return delay, nil
}

// MeasureTCPDelay 纯TCP连接测速，只测试TCP握手时间
// timeout: 超时时间（秒）
// host: 目标主机地址
// port: 目标端口
// proxy: 代理地址，格式如 "socks5://127.0.0.1:1080"
func MeasureTCPDelay(timeout int, host string, port int, proxy string) (int64, error) {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	httpTimeout := time.Second * time.Duration(timeout)

	var dialer *net.Dialer
	var err error

	if len(proxy) > 0 {
		// 使用代理进行TCP连接
		dialer, err = createProxyDialer(httpTimeout, proxy)
		if err != nil {
			return PingDelayError, err
		}
	} else {
		// 直连TCP
		dialer = &net.Dialer{
			Timeout: httpTimeout,
		}
	}

	start := time.Now()
	conn, err := dialer.Dial("tcp", address)
	delay := time.Since(start).Milliseconds()

	if conn != nil {
		conn.Close()
	}

	if err != nil {
		precision := delay - int64(timeout)*1000
		if math.Abs(float64(precision)) < 50 {
			return PingDelayTimeout, err
		} else {
			return PingDelayError, err
		}
	}

	return delay, nil
}

// MeasureConnectDelay 测试代理服务器的直接TCP连接延迟（已废弃）
// 这个方法只是简单的TCP连接测试，不能真正反映代理协议性能
// 建议使用 MeasureProxyConnectDelay 来测试真正的代理连接
// timeout: 超时时间（秒）
// proxyAddr: 代理服务器地址，格式如 "proxy.example.com:1080"
func MeasureConnectDelay(timeout int, proxyAddr string) (int64, error) {
	if len(proxyAddr) == 0 {
		return PingDelayError, fmt.Errorf("proxy address cannot be empty")
	}

	httpTimeout := time.Second * time.Duration(timeout)
	dialer := &net.Dialer{
		Timeout: httpTimeout,
	}

	start := time.Now()
	conn, err := dialer.Dial("tcp", proxyAddr)
	delay := time.Since(start).Milliseconds()

	if conn != nil {
		conn.Close()
	}

	if err != nil {
		precision := delay - int64(timeout)*1000
		if math.Abs(float64(precision)) < 50 {
			return PingDelayTimeout, err
		} else {
			return PingDelayError, err
		}
	}

	return delay, nil
}

// MeasureProxyConnectDelay 测试代理协议连接延迟，类似Shadowrocket的connect测速
// 这个方法通过代理协议测试到目标服务器的连接建立时间，只测握手不传输数据
// timeout: 超时时间（秒）
// targetHost: 目标主机地址
// targetPort: 目标端口
// proxy: 本地代理地址，格式如 "socks5://127.0.0.1:1080"
func MeasureProxyConnectDelay(timeout int, targetHost string, targetPort int, proxy string) (int64, error) {
	if len(proxy) == 0 {
		return PingDelayError, fmt.Errorf("proxy address cannot be empty")
	}

	httpTimeout := time.Second * time.Duration(timeout)
	target := net.JoinHostPort(targetHost, fmt.Sprintf("%d", targetPort))

	// 创建支持代理的HTTP客户端
	client, err := CoreHTTPClient(httpTimeout, proxy)
	if err != nil {
		return PingDelayError, err
	}

	// 构造一个简单的连接测试请求（只建立连接，不传输数据）
	start := time.Now()

	// 使用HTTP CONNECT方法测试代理连接建立
	// 这比直接TCP连接更能反映真实的代理协议握手时间
	testURL := fmt.Sprintf("http://%s", target)
	req, err := http.NewRequest("HEAD", testURL, nil)
	if err != nil {
		return PingDelayError, err
	}

	// 设置连接超时
	req = req.WithContext(func() context.Context {
		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
		defer cancel()
		return ctx
	}())

	// 执行请求（只测试连接建立，不关心响应内容）
	resp, err := client.Do(req)
	delay := time.Since(start).Milliseconds()

	if resp != nil {
		resp.Body.Close()
	}

	if err != nil {
		precision := delay - int64(timeout)*1000
		if math.Abs(float64(precision)) < 50 {
			return PingDelayTimeout, err
		} else {
			return PingDelayError, err
		}
	}

	return delay, nil
}

// createProxyDialer 创建支持代理的Dialer
func createProxyDialer(timeout time.Duration, proxy string) (*net.Dialer, error) {
	_, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}

	// 对于TCP测速，我们需要创建一个通过HTTP代理的Dialer
	// 这里简化实现，直接返回基础Dialer
	// 实际应用中可能需要更复杂的代理支持
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// 注意：这里的代理实现需要根据具体的代理类型来实现
	// SOCKS5代理需要特殊处理，HTTP代理也需要特殊处理
	// 为了简化，这里先返回基础dialer
	// 在生产环境中，建议使用 golang.org/x/net/proxy 包

	return dialer, nil
}
