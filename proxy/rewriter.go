package proxy

import (
	"errors"
	"net/url"
	"strings"
)

// RewriteURL 从请求路径中解析出目标 URL
// 支持以下格式:
// 1. /https://github.com/user/repo -> https://github.com/user/repo
// 2. /github.com/user/repo -> https://github.com/user/repo (默认 https)
func RewriteURL(path string) (*url.URL, error) {
	// 去掉开头的 /
	targetStr := strings.TrimPrefix(path, "/")

	if targetStr == "" {
		return nil, errors.New("empty url")
	}

	// 检查是否包含 scheme
	// 简单的启发式检查：如果开头不是 http:// 或 https://，则默认添加 https://
	if !strings.HasPrefix(targetStr, "http://") && !strings.HasPrefix(targetStr, "https://") {
		// 这里可以根据需要添加更多复杂的判断，比如检查是否符合域名格式
		// 目前简单处理：默认补全 https://
		targetStr = "https://" + targetStr
	}

	// 解析 URL
	// 注意：如果原始路径中有 query parameters (e.g. ?foo=bar)，
	// 在 http.Handler 中应该从 r.URL.RawQuery 获取并拼接到这里，
	// 但 RewriteURL 函数只处理 path 部分，调用者需要负责拼接 Query。
	// 这里我们假设传入的 path 仅仅是路径部分，不包含 query。
	
	// 有一种情况：url 中包含多次 http://，比如 http://localhost/http://example.com
	// net/url 解析可能会有问题，需要仔细处理。
	
	u, err := url.Parse(targetStr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported scheme")
	}

	if u.Host == "" {
		return nil, errors.New("missing host")
	}

	return u, nil
}
