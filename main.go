package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// 设置路由处理函数
	http.HandleFunc("/", echoHandler)

	// 启动服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // 默认端口
	}
	fmt.Printf("Echo Server 正在监听端口 %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// 创建响应构建器
	var sb strings.Builder

	sb.WriteString("=== Echo Server 响应 ===\n\n")

	// 添加请求基本信息
	sb.WriteString("=== 请求信息 ===\n")
	sb.WriteString(fmt.Sprintf("方法: %s\n", r.Method))
	sb.WriteString(fmt.Sprintf("路径: %s\n", r.URL.Path))
	sb.WriteString(fmt.Sprintf("协议: %s\n", r.Proto))
	sb.WriteString(fmt.Sprintf("主机: %s\n", r.Host))
	sb.WriteString(fmt.Sprintf("远程地址: %s\n", r.RemoteAddr))

	// 添加请求头
	sb.WriteString("\n=== 请求头 ===\n")
	for name, headers := range r.Header {
		for _, h := range headers {
			sb.WriteString(fmt.Sprintf("%v: %v\n", name, h))
		}
	}

	// 添加查询参数
	sb.WriteString("\n=== 查询参数 ===\n")
	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		sb.WriteString("无\n")
	} else {
		for key, values := range queryParams {
			for _, value := range values {
				sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
			}
		}
	}

	// 处理请求体
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		// 读取请求体内容
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			sb.WriteString(fmt.Sprintf("\n读取请求体错误: %v\n", err))
		} else {
			// 将读取的内容重新放回Body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			// 添加原始请求体
			sb.WriteString("\n=== 原始请求体 ===\n")
			if len(bodyBytes) == 0 {
				sb.WriteString("空\n")
			} else {
				sb.WriteString(string(bodyBytes) + "\n")
			}

			// 尝试解析表单数据
			if err := r.ParseForm(); err != nil {
				sb.WriteString(fmt.Sprintf("\n解析表单错误: %v\n", err))
			} else {
				// 添加表单数据
				sb.WriteString("\n=== 表单数据 ===\n")
				if len(r.PostForm) == 0 {
					sb.WriteString("无\n")
				} else {
					for key, values := range r.PostForm {
						for _, value := range values {
							sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
						}
					}
				}
			}
		}
	}


	// 添加环境变量信息
	sb.WriteString("\n=== 环境变量 ===\n")
	for _, env := range os.Environ() {
		sb.WriteString(env + "\n")
	}
	sb.WriteString("\n")


	// 发送响应
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sb.String()))

	// 同时在控制台打印简要信息（可选）
	fmt.Printf("\n处理请求: %s %s 来自 %s\n", r.Method, r.URL.Path, r.RemoteAddr)
}