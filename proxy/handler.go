package proxy

import (
	"io"
	"log"
	"net/http"
)

// ProxyHandler 处理代理请求
type ProxyHandler struct{}

// NewProxyHandler 创建一个新的代理处理器
func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{}
}

// ServeHTTP 实现 http.Handler 接口
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. 解析目标 URL
	targetURL, err := RewriteURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 拼接原始请求的 Query 参数
	if r.URL.RawQuery != "" {
		targetURL.RawQuery = r.URL.RawQuery
	}

	log.Printf("Proxying request: %s -> %s", r.URL.Path, targetURL.String())

	// 3. 创建发往目标服务器的新请求
	// method, url, body 保持一致
	outReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. 复制请求头
	copyHeader(outReq.Header, r.Header)

	// 5. 关键：设置 Host 头为目标 Host，否则很多服务会拒绝
	outReq.Host = targetURL.Host
	// 清除 RequestURI，因为它只在客户端请求中有意义
	outReq.RequestURI = ""

	// 6. 发起请求
	// 使用默认的 Client，支持 HTTP/2
	// 如果需要更细粒度的控制（如超时、重定向策略），可以自定义 http.Client
	client := &http.Client{
		// 不自动跟踪重定向，将重定向响应直接返回给客户端，让客户端决定
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(outReq)
	if err != nil {
		log.Printf("Request failed: %v", err)
		http.Error(w, "Request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 7. 复制响应头
	copyHeader(w.Header(), resp.Header)

	// 8. 设置响应状态码
	w.WriteHeader(resp.StatusCode)

	// 9. 流式复制响应体
	// io.Copy 会处理流式传输，适合大文件
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// 此时 Header 已经发送，无法再发送 http.Error
		log.Printf("Failed to copy response body: %v", err)
	}
}

// copyHeader 复制 http.Header
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		// 过滤掉一些 Hop-by-hop headers，虽然 net/http 也会自动处理一部分，但显式处理更安全
		// 这里简单起见全部复制，net/http 发送时会自动剔除 Connection 等头
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
	
	// 显式删除一些不应该转发的头（如果需要）
	// dst.Del("Connection")
	// dst.Del("Keep-Alive")
	// ... 但 net/http 的 Transport 通常会处理得很好
}

// IsRootPath 检查是否是根路径访问
func IsRootPath(path string) bool {
	return path == "/" || path == ""
}
