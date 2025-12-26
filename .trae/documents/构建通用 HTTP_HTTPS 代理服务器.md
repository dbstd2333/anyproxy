# 构建基于 URL 路径重写的 HTTP 加速代理

根据您的最新需求，该程序将实现一个反向代理服务：通过解析用户请求的 URL 路径作为目标地址（例如 `https://yoursite.com/github.com/user/repo`），并在后端请求该目标地址后将内容流式返回给用户。这常用于资源加速下载。

## 项目结构
- `main.go`: 程序入口，启动 Web 服务器。
- `proxy/`: 代理核心包。
  - `handler.go`: 核心请求处理逻辑。
  - `rewriter.go`: URL 解析与重写逻辑。
- `go.mod`: 模块定义。

## 技术实现细节

### 1. 初始化
- `go mod init anyproxy`

### 2. URL 解析与重写 (`proxy/rewriter.go`)
- 实现逻辑从 `Request.URL.Path` 中提取目标 URL。
- **智能补全**：
  - 支持 `https://domain/http://target.com/resource`（显式协议）。
  - 支持 `https://domain/target.com/resource`（隐式协议，默认补全 `https://`）。
- 验证目标 URL 的合法性。

### 3. 核心代理逻辑 (`proxy/handler.go`)
- 实现 `ProxyHandler`。
- **请求构建**：
  - 创建指向目标 URL 的新请求。
  - **关键点**：必须重写 `Host` 请求头为目标域名，否则大多云服务（如 GitHub）会拒绝请求。
  - 复制客户端的原始 Header（如 `Range`, `User-Agent`, `Authorization` 等）到新请求，确保 Git 协议或断点续传正常工作。
  - 处理 `X-Forwarded-For` 等代理头。
- **响应转发**：
  - 发起请求获取响应。
  - 将目标响应的 Header（如 `Content-Type`, `Content-Length`, `Content-Disposition`）复制回客户端。
  - 写入状态码。
  - 使用 `io.Copy` 将响应体流式传输给客户端，支持大文件传输。

### 4. 入口文件 (`main.go`)
- 路由处理：
  - `/` (根路径): 显示简易的使用说明 HTML。
  - `/*`: 通配符路径，交由 `ProxyHandler` 处理。

## 验证计划
- 启动服务。
- 访问 `http://localhost:8080/github.com/torvalds/linux` 验证是否会被重定向或代理到 GitHub。
- 验证 Header 转发是否正确（如 `Host` 头修改）。

这个方案完全符合您描述的 `https://domain/github.com/aaa/bbb.git` 访问模式。