这个是一个echoserver ai帮我写的

k8s 的调试工具，请求它，会打印出：

```bash
=== Echo Server 响应 ===

=== 请求信息 ===
方法: GET
路径: /
协议: HTTP/1.1
主机: localhost:8888
远程地址: [::1]:53930

=== 请求头 ===
Connection: keep-alive
Sec-Ch-Ua-Mobile: ?0
Sec-Ch-Ua-Platform: "macOS"
Upgrade-Insecure-Requests: 1
Cookie: wp-settings-time-1=1748330732; _ga=GA1.1.734620496.1752248056; Hm_lvt_50b516d50102c8c9ac5f80529b81ca17=1752248055,1752555511; _ga_YN0WWZGVYN=GS2.1.s1752555512$o2$g1$t1752555794$j60$l0$h0
Cache-Control: max-age=0
Sec-Ch-Ua: "Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Sec-Fetch-Dest: document
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6
Sec-Fetch-Mode: navigate
Sec-Fetch-User: ?1
Accept-Encoding: gzip, deflate, br, zstd
Sec-Fetch-Site: none

=== 查询参数 ===
无

=== 指定环境变量 ===
A: (未设置)
B: (未设置)
c: (未设置)
```

可以做到请求回显。配合k8s的envfrom可以做到打印出pod名字、pod IP等信息

```yml
          - name: NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
          - name: POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
```

## 支持的环境变量

SERVER_PORT=8888
ECHO_ENV=NODE_NAME,POD_NAME,POD_NAMESPACE,POD_IP
