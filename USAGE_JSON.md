# libXray JSON配置测速使用指南

## 概述

libXray 现在支持两种测速方式：
1. **配置文件方式（原有）**：需要将配置保存到文件
2. **JSON字符串方式（新增）**：直接传入JSON配置，无需文件

## 三种测速方法

### 1. HTTP测速 - PingFromJSON
完整的HTTP请求测速，包含DNS解析、TCP连接、TLS握手等。

### 2. TCP测速 - PingTCPFromJSON  
纯TCP连接测速，只测试网络层连接性能。

### 3. Connect测速 - ConnectFromJSON
代理协议连接测速，测试代理握手性能。

## Go语言使用示例

```go
package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    
    libXray "github.com/xtls/libxray"
)

func main() {
    // 示例Xray配置JSON
    xrayConfig := map[string]interface{}{
        "inbounds": []map[string]interface{}{
            {
                "port":     10808,
                "protocol": "socks",
                "settings": map[string]interface{}{
                    "auth": "noauth",
                },
            },
        },
        "outbounds": []map[string]interface{}{
            {
                "protocol": "vmess",
                "settings": map[string]interface{}{
                    "vnext": []map[string]interface{}{
                        {
                            "address": "example.com",
                            "port":    443,
                            "users": []map[string]interface{}{
                                {
                                    "id":       "your-uuid-here",
                                    "alterId":  0,
                                    "security": "auto",
                                },
                            },
                        },
                    },
                },
                "streamSettings": map[string]interface{}{
                    "network": "ws",
                    "security": "tls",
                    "wsSettings": map[string]interface{}{
                        "path": "/path",
                    },
                },
            },
        },
    }
    
    // 转换为JSON字符串
    configBytes, _ := json.Marshal(xrayConfig)
    configJSON := string(configBytes)
    
    // 测试参数
    datDir := "/path/to/dat"  // geosite.dat 和 geoip.dat 目录
    proxy := "socks5://127.0.0.1:10808"  // 本地代理地址
    
    // 1. HTTP测速
    testHTTPPing(datDir, configJSON, proxy)
    
    // 2. TCP测速  
    testTCPPing(datDir, configJSON, proxy)
    
    // 3. Connect测速
    testConnectPing(datDir, configJSON, proxy)
}

func testHTTPPing(datDir, configJSON, proxy string) {
    // 创建请求结构体
    request := map[string]interface{}{
        "datDir":     datDir,
        "configJSON": configJSON,
        "timeout":    10,
        "url":        "https://www.google.com",
        "proxy":      proxy,
    }
    
    // 编码为base64
    requestBytes, _ := json.Marshal(request)
    base64Request := base64.StdEncoding.EncodeToString(requestBytes)
    
    // 调用测速
    result := libXray.PingFromJSON(base64Request)
    
    // 解析结果
    decoded, _ := base64.StdEncoding.DecodeString(result)
    var response struct {
        Success bool   `json:"success"`
        Data    int64  `json:"data"`
        Err     string `json:"error"`
    }
    json.Unmarshal(decoded, &response)
    
    if response.Success {
        fmt.Printf("HTTP测速: %dms\n", response.Data)
    } else {
        fmt.Printf("HTTP测速失败: %s\n", response.Err)
    }
}

func testTCPPing(datDir, configJSON, proxy string) {
    request := map[string]interface{}{
        "datDir":     datDir,
        "configJSON": configJSON,
        "timeout":    5,
        "host":       "8.8.8.8",
        "port":       53,
        "proxy":      proxy,
    }
    
    requestBytes, _ := json.Marshal(request)
    base64Request := base64.StdEncoding.EncodeToString(requestBytes)
    result := libXray.PingTCPFromJSON(base64Request)
    
    // 解析结果
    decoded, _ := base64.StdEncoding.DecodeString(result)
    var response struct {
        Success bool   `json:"success"`
        Data    int64  `json:"data"`
        Err     string `json:"error"`
    }
    json.Unmarshal(decoded, &response)
    
    if response.Success {
        fmt.Printf("TCP测速: %dms\n", response.Data)
    } else {
        fmt.Printf("TCP测速失败: %s\n", response.Err)
    }
}

func testConnectPing(datDir, configJSON, proxy string) {
    request := map[string]interface{}{
        "datDir":     datDir,
        "configJSON": configJSON,
        "timeout":    8,
        "targetHost": "google.com",
        "targetPort": 443,
        "proxy":      proxy,
    }
    
    requestBytes, _ := json.Marshal(request)
    base64Request := base64.StdEncoding.EncodeToString(requestBytes)
    result := libXray.ConnectFromJSON(base64Request)
    
    // 解析结果
    decoded, _ := base64.StdEncoding.DecodeString(result)
    var response struct {
        Success bool   `json:"success"`
        Data    int64  `json:"data"`
        Err     string `json:"error"`
    }
    json.Unmarshal(decoded, &response)
    
    if response.Success {
        fmt.Printf("Connect测速: %dms\n", response.Data)
    } else {
        fmt.Printf("Connect测速失败: %s\n", response.Err)
    }
}
```

## Android/Java使用示例

```java
import org.json.JSONObject;
import android.util.Base64;

public class XrayTester {
    
    // HTTP测速
    public void testHTTPPingJSON() {
        try {
            // 创建请求
            JSONObject request = new JSONObject();
            request.put("datDir", "/android_asset/dat");
            request.put("configJSON", getXrayConfigJSON());  // 你的Xray配置JSON
            request.put("timeout", 10);
            request.put("url", "https://www.google.com");
            request.put("proxy", "socks5://127.0.0.1:10808");
            
            // 编码请求
            String jsonStr = request.toString();
            String base64Request = Base64.encodeToString(jsonStr.getBytes(), Base64.NO_WRAP);
            
            // 执行测速
            new Thread(() -> {
                String result = LibXray.pingFromJSON(base64Request);
                handleResult("HTTP", result);
            }).start();
            
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    
    // TCP测速
    public void testTCPPingJSON() {
        try {
            JSONObject request = new JSONObject();
            request.put("datDir", "/android_asset/dat");
            request.put("configJSON", getXrayConfigJSON());
            request.put("timeout", 5);
            request.put("host", "8.8.8.8");
            request.put("port", 53);
            request.put("proxy", "socks5://127.0.0.1:10808");
            
            String jsonStr = request.toString();
            String base64Request = Base64.encodeToString(jsonStr.getBytes(), Base64.NO_WRAP);
            
            new Thread(() -> {
                String result = LibXray.pingTCPFromJSON(base64Request);
                handleResult("TCP", result);
            }).start();
            
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    
    // Connect测速
    public void testConnectJSON() {
        try {
            JSONObject request = new JSONObject();
            request.put("datDir", "/android_asset/dat");
            request.put("configJSON", getXrayConfigJSON());
            request.put("timeout", 8);
            request.put("targetHost", "google.com");
            request.put("targetPort", 443);
            request.put("proxy", "socks5://127.0.0.1:10808");
            
            String jsonStr = request.toString();
            String base64Request = Base64.encodeToString(jsonStr.getBytes(), Base64.NO_WRAP);
            
            new Thread(() -> {
                String result = LibXray.connectFromJSON(base64Request);
                handleResult("Connect", result);
            }).start();
            
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    
    private void handleResult(String type, String result) {
        try {
            // 解析结果
            byte[] decoded = Base64.decode(result, Base64.DEFAULT);
            String jsonResult = new String(decoded);
            JSONObject response = new JSONObject(jsonResult);
            
            if (response.getBoolean("success")) {
                long delay = response.getLong("data");
                Log.i("PING_" + type, "延迟: " + delay + "ms");
            } else {
                String error = response.getString("error");
                Log.e("PING_" + type, "测速失败: " + error);
            }
        } catch (Exception e) {
            Log.e("PING_" + type, "解析失败", e);
        }
    }
    
    private String getXrayConfigJSON() {
        // 返回你的Xray配置JSON字符串
        return "{ \"inbounds\": [...], \"outbounds\": [...] }";
    }
}
```

## C/C++使用示例

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 假设已经链接了libxray动态库
extern char* CGoPingFromJSON(char* base64Text);
extern char* CGoPingTCPFromJSON(char* base64Text);
extern char* CGoConnectFromJSON(char* base64Text);

void test_http_ping() {
    // JSON请求（需要先进行base64编码）
    char* request = "{\"datDir\":\"/path/to/dat\",\"configJSON\":\"{...}\",\"timeout\":10,\"url\":\"https://www.google.com\",\"proxy\":\"socks5://127.0.0.1:10808\"}";
    
    // 这里需要进行base64编码，示例中省略
    char* base64_request = base64_encode(request);
    
    char* result = CGoPingFromJSON(base64_request);
    printf("HTTP测速结果: %s\n", result);
    
    // 释放内存
    free(result);
    free(base64_request);
}

void test_tcp_ping() {
    char* request = "{\"datDir\":\"/path/to/dat\",\"configJSON\":\"{...}\",\"timeout\":5,\"host\":\"8.8.8.8\",\"port\":53,\"proxy\":\"socks5://127.0.0.1:10808\"}";
    char* base64_request = base64_encode(request);
    
    char* result = CGoPingTCPFromJSON(base64_request);
    printf("TCP测速结果: %s\n", result);
    
    free(result);
    free(base64_request);
}
```

## 参数说明

### 公共参数
- `datDir`: geosite.dat 和 geoip.dat 文件所在目录
- `configJSON`: Xray配置的JSON字符串
- `timeout`: 超时时间（秒）
- `proxy`: 本地代理地址，格式如 "socks5://127.0.0.1:1080"

### HTTP测速特有参数
- `url`: 测试网址，推荐 "https://www.google.com"

### TCP测速特有参数  
- `host`: 目标主机，如 "8.8.8.8"
- `port`: 目标端口，如 53（DNS）、80（HTTP）、443（HTTPS）

### Connect测速特有参数
- `targetHost`: 目标主机，如 "google.com"
- `targetPort`: 目标端口，如 443

## 返回值说明

所有测速方法返回base64编码的JSON，解码后格式：

```json
{
    "success": true,
    "data": 150,     // 延迟毫秒数
    "error": ""
}
```

或失败时：

```json
{
    "success": false,
    "data": 10000,   // 错误码：10000(一般错误) 或 11000(超时)
    "error": "错误信息"
}
```

## 优势

1. **无需文件依赖**：直接传入JSON配置，不需要创建配置文件
2. **内存优化**：测速完成后自动释放内存
3. **性能优秀**：保持与原版本相同的测速性能
4. **完全兼容**：保留原有配置文件方式，新增JSON方式

## 注意事项

1. 确保 `datDir` 目录包含 `geosite.dat` 和 `geoip.dat` 文件
2. `configJSON` 必须是有效的Xray配置
3. `proxy` 地址要与配置中的入站端口匹配
4. 测速是阻塞操作，建议在后台线程执行
