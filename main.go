package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
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
	// 记录请求开始时间
	startTime := time.Now()

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

	// 获取并排序请求头键名
	headerKeys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)

	// 按排序后的顺序输出请求头
	sb.WriteString("\n=== 请求头 ===\n")
	for _, k := range headerKeys {
		for _, v := range r.Header[k] {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
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

	// 添加筛选后的环境变量信息
	sb.WriteString("\n=== 指定环境变量 ===\n")
	echoEnv := os.Getenv("ECHO_ENV")
	if echoEnv != "" {
		// 解析要显示的环境变量名
		envVarsToShow := strings.Split(echoEnv, ",")
		for _, envName := range envVarsToShow {
			envName = strings.TrimSpace(envName)
			if envValue, exists := os.LookupEnv(envName); exists {
				sb.WriteString(fmt.Sprintf("%s: %s\n", envName, envValue))
			} else {
				sb.WriteString(fmt.Sprintf("%s: (未设置)\n", envName))
			}
		}
	} else {
		sb.WriteString("ECHO_ENV 环境变量未设置，无指定环境变量可显示\n")
	}
	sb.WriteString("\n")

	// 计算处理耗时
	processingTime := time.Since(startTime)
	sb.WriteString(fmt.Sprintf("\n请求处理耗时: %v\n", processingTime))

	// 添加当前时间
	currentTime := time.Now().Format(time.RFC3339)
	sb.WriteString(fmt.Sprintf("当前时间: %s\n", currentTime))

	// 发送响应
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sb.String()))

	// 同时在控制台打印简要信息（可选）
	fmt.Printf("\n处理请求: %s %s 来自 %s\n", r.Method, r.URL.Path, r.RemoteAddr)
}
