# 深圳市星程国际旅行社 — 官方演示站 (Gin)

这是一个使用 Go + Gin 快速搭建的官网原型，供演示与本地开发使用。

快速开始（Windows PowerShell）：

```powershell
# 确保已安装 Go (>=1.20)
go version; go mod tidy
# 本地运行
go run main.go
# 或构建并运行
go build -o star-tour.exe .; .\star-tour.exe
```

默认监听 http://localhost:8080

文件说明：
- `main.go` - 程序入口，路由与模板加载
- `templates/` - HTML 模板
- `static/` - 静态资源 (CSS, images)

如需容器化，可使用提供的 `Dockerfile` 构建镜像。

配置邮件服务
1. 复制 `config.yaml.example` 为 `config.yaml` 并填写 SMTP 信息。
2. 可通过环境变量 `CONFIG_PATH` 指定其他配置文件路径。

示例 `config.yaml`:
```yaml
smtp:
	host: "smtp.example.com"
	port: 587
	user: "your@from.email"
	pass: "yourpassword"
	contact_to: "notify@yourcompany.com"
```

然后运行服务并访问 `/contact` 提交表单以测试邮件发送。