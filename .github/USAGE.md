# libXray Android 使用指南

## 获取AAR文件

### 方式1: 从GitHub Releases下载
1. 访问 [Releases页面](../../releases)
2. 下载最新的 `libXray.aar` 文件

### 方式2: 使用命令行下载
```bash
# 获取最新release的下载链接
LATEST_RELEASE=$(curl -s https://api.github.com/repos/YOUR_USERNAME/libXray/releases/latest | grep "browser_download_url.*libXray.aar" | cut -d '"' -f 4)

# 下载AAR文件
wget $LATEST_RELEASE
```

## 集成到Android项目

### 1. 添加AAR文件
```bash
# 将AAR文件复制到项目的libs目录
mkdir -p app/libs
cp libXray.aar app/libs/
```

### 2. 配置Gradle依赖
在 `app/build.gradle` 中添加：

```gradle
android {
    // ... 其他配置
}

dependencies {
    implementation files('libs/libXray.aar')
    // ... 其他依赖
}
```

### 3. 权限配置
在 `AndroidManifest.xml` 中添加必要权限：

```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

## 使用示例

### TCP测速 (新功能)

```java
import libXray.LibXray;
import org.json.JSONObject;
import android.util.Base64;

public class NetworkTester {
    
    // 创建TCP测速请求
    private String createTCPPingRequest(String datDir, String configPath, 
                                       int timeout, String host, int port, String proxy) {
        try {
            JSONObject request = new JSONObject();
            request.put("datDir", datDir);
            request.put("configPath", configPath);
            request.put("timeout", timeout);
            request.put("host", host);
            request.put("port", port);
            request.put("proxy", proxy);
            
            String jsonStr = request.toString();
            return Base64.encodeToString(jsonStr.getBytes(), Base64.NO_WRAP);
        } catch (Exception e) {
            e.printStackTrace();
            return "";
        }
    }
    
    // 执行TCP测速
    public void testTCPPing() {
        String request = createTCPPingRequest(
            "/android_asset/dat",        // 数据文件目录
            "/path/to/config.json",      // Xray配置文件
            10,                          // 10秒超时
            "8.8.8.8",                  // Google DNS
            53,                         // DNS端口
            "socks5://127.0.0.1:1080"   // 本地代理
        );
        
        // 异步执行测速
        new Thread(() -> {
            try {
                String result = LibXray.pingTCP(request);
                
                // 解析结果
                byte[] decoded = Base64.decode(result, Base64.DEFAULT);
                String jsonResult = new String(decoded);
                JSONObject response = new JSONObject(jsonResult);
                
                if (response.getBoolean("success")) {
                    long delay = response.getLong("data");
                    runOnUiThread(() -> {
                        // 更新UI - 显示延迟结果
                        Log.i("TCP_PING", "延迟: " + delay + "ms");
                    });
                } else {
                    String error = response.getString("error");
                    Log.e("TCP_PING", "测速失败: " + error);
                }
            } catch (Exception e) {
                Log.e("TCP_PING", "执行失败", e);
            }
        }).start();
    }
}
```

### HTTP测速 (原有功能)

```java
public void testHTTPPing() {
    String request = createHTTPPingRequest(
        "/android_asset/dat",
        "/path/to/config.json",
        10,
        "https://connectivitycheck.gstatic.com/generate_204",
        "socks5://127.0.0.1:1080"
    );
    
    new Thread(() -> {
        String result = LibXray.ping(request);
        // 处理结果...
    }).start();
}
```

## 测速对比

| 测试类型 | 方法 | 适用场景 | 优势 |
|---------|------|---------|------|
| **TCP测速** | `LibXray.pingTCP()` | 纯网络连接测试 | 更快、更纯粹的网络性能 |
| **HTTP测速** | `LibXray.ping()` | Web访问体验测试 | 更接近真实使用场景 |

## 常用测试目标

```java
// DNS服务器测试
testTCP("8.8.8.8", 53);        // Google DNS
testTCP("1.1.1.1", 53);        // Cloudflare DNS

// Web服务器测试
testTCP("google.com", 80);      // HTTP
testTCP("google.com", 443);     // HTTPS

// 自定义服务器测试
testTCP("your-server.com", 8080);
```

## 故障排除

### 1. AAR导入失败
- 确保AAR文件放在正确的 `app/libs/` 目录
- 检查Gradle配置是否正确
- 清理并重新构建项目: `./gradlew clean build`

### 2. 运行时错误
- 确保应用有网络权限
- 检查配置文件路径是否正确
- 验证代理地址格式

### 3. 测速结果异常
- 检查网络连接
- 验证目标主机和端口是否可达
- 确认代理服务器是否正常运行

## 更多资源

- [libXray GitHub仓库](https://github.com/XTLS/libXray)
- [Xray-core文档](https://xtls.github.io/)
- [问题反馈](../../issues)
