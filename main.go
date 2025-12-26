package main

import (
	"anyproxy/proxy"
	"fmt"
	"log"
	"net/http"
	"os"
)

const helpHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>AnyProxy - URL Accelerator</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; line-height: 1.6; }
        h1 { color: #333; }
        code { background: #f4f4f4; padding: 2px 5px; border-radius: 3px; }
        pre { background: #f4f4f4; padding: 15px; border-radius: 5px; overflow-x: auto; }
        .example { margin-top: 20px; }
    </style>
</head>
<body>
    <h1>AnyProxy 加速服务</h1>
    <p>这是一个简单的反向代理服务，用于加速访问目标资源。</p>
    
    <h2>使用方法</h2>
    <p>将目标 URL 附加到本服务地址后面：</p>
    <pre>http://localhost:8080/{target_url}</pre>
    
    <div class="example">
        <h3>示例：</h3>
        <p>加速访问 GitHub 仓库：</p>
        <code>http://localhost:8080/github.com/torvalds/linux</code>
        <br><br>
        <p>显式指定协议：</p>
        <code>http://localhost:8080/https://www.google.com</code>
    </div>
</body>
</html>
`

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	proxyHandler := proxy.NewProxyHandler()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 处理根路径，显示帮助信息
		if proxy.IsRootPath(r.URL.Path) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, helpHTML)
			return
		}

		// 其他路径交给代理处理器
		proxyHandler.ServeHTTP(w, r)
	})

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
