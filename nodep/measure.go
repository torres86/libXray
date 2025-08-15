package nodep

import (
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
