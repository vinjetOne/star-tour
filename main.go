package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	gomail "gopkg.in/gomail.v2"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	// 先读取配置文件（默认 config.yaml，或使用 CONFIG_PATH 环境变量覆盖）
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		// 在部署环境（如 Render）我们通常使用环境变量来配置敏感信息，
		// 所以这里不强制终止，改为记录并继续；后续会用环境变量覆盖或补足配置。
		log.Printf("加载配置失败 (可忽略): %v，继续使用环境变量覆盖或默认值", err)
		cfg = &Config{}
	}

	// 用环境变量覆盖配置（优先使用环境变量）
	applyEnvOverrides(cfg)

	r := gin.Default()
	// 静态资源
	r.Static("/static", "./static")

	// 模板
	r.LoadHTMLGlob("templates/*")

	// 路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"Title":   "深圳市星程国际旅行社",
			"Contact": cfg.Contact,
		})
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about", gin.H{"Title": "关于我们 - 深圳市星程国际旅行社", "Contact": cfg.Contact})
	})

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact", gin.H{"Title": "联系我们 - 深圳市星程国际旅行社", "Contact": cfg.Contact})
	})

	// 处理联系表单提交（支持 JSON 和 form 数据）
	r.POST("/contact", func(c *gin.Context) {
		var name, email, phone, message string

		// 读取原始请求体，优先尝试 JSON 解析；若失败回退为表单
		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("读取请求体失败: %v", err)
		}
		log.Printf("请求 Content-Type: %s", c.GetHeader("Content-Type"))
		log.Printf("请求原始体: %s", string(raw))

		// 尝试解析 JSON
		var payload struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Phone   string `json:"phone"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(raw, &payload); err == nil {
			name = payload.Name
			email = payload.Email
			phone = payload.Phone
			message = payload.Message
		} else {
			// 恢复请求体以便后续 ParseForm/PostForm 使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
			// 解析表单数据（支持 x-www-form-urlencoded / multipart）
			if err := c.Request.ParseForm(); err == nil {
				name = c.Request.PostFormValue("name")
				email = c.Request.PostFormValue("email")
				phone = c.Request.PostFormValue("phone")
				message = c.Request.PostFormValue("message")
			} else {
				log.Printf("ParseForm 失败: %v", err)
			}
		}

		if name == "" || email == "" || message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必填字段"})
			return
		}

		// 从启动时加载的配置中读取 SMTP 配置
		smtpHost := cfg.SMTP.Host
		smtpPort := cfg.SMTP.Port
		smtpUser := cfg.SMTP.User
		smtpPass := cfg.SMTP.Pass
		contactTo := cfg.SMTP.ContactTo

		if smtpHost == "" || smtpPort == 0 || smtpUser == "" || smtpPass == "" || contactTo == "" {
			c.String(http.StatusInternalServerError, "邮件服务未配置，请联系管理员")
			return
		}

		// 构建邮件内容
		subject := fmt.Sprintf("新留言 - %s", name)
		body := fmt.Sprintf("姓名: %s\n邮箱: %s\n电话: %s\n\n留言:\n%s", name, email, phone, message)

		m := gomail.NewMessage()
		m.SetHeader("From", smtpUser)
		m.SetHeader("To", contactTo)
		m.SetHeader("Subject", subject)
		m.SetBody("text/plain", body)

		d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

		if err := d.DialAndSend(m); err != nil {
			// 记录错误并返回（不泄露密码）
			log.Printf("邮件发送失败: to=%s subject=%s err=%v", contactTo, subject, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发送邮件失败"})
			return
		}

		log.Printf("邮件已发送: from=%s to=%s subject=%s", smtpUser, contactTo, subject)
		c.JSON(http.StatusOK, gin.H{"message": "感谢您的留言，我们已收到并会尽快联系您。"})
	})

	// 未命中路由返回 404 简单页面，避免渲染首页内容
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{"Title": "404 - 页面未找到"})
	})

	// 在 PaaS（例如 Render）上，端口通常由环境变量提供（PORT）。默认为 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// applyEnvOverrides 会将常见的环境变量覆盖到 cfg 中，便于在部署平台上配置
func applyEnvOverrides(cfg *Config) {
	if cfg == nil {
		return
	}
	if v := os.Getenv("SMTP_HOST"); v != "" {
		cfg.SMTP.Host = v
	}
	if v := os.Getenv("SMTP_PORT"); v != "" {
		// 尝试解析为 int
		var p int
		if _, err := fmt.Sscanf(v, "%d", &p); err == nil {
			cfg.SMTP.Port = p
		}
	}
	if v := os.Getenv("SMTP_USER"); v != "" {
		cfg.SMTP.User = v
	}
	if v := os.Getenv("SMTP_PASS"); v != "" {
		cfg.SMTP.Pass = v
	}
	if v := os.Getenv("SMTP_CONTACT_TO"); v != "" {
		cfg.SMTP.ContactTo = v
	}

	if v := os.Getenv("CONTACT_PHONE"); v != "" {
		cfg.Contact.Phone = v
	}
	if v := os.Getenv("CONTACT_EMAIL"); v != "" {
		cfg.Contact.Email = v
	}
	if v := os.Getenv("CONTACT_WECHAT"); v != "" {
		cfg.Contact.Wechat = v
	}
	if v := os.Getenv("CONTACT_ADDRESS"); v != "" {
		cfg.Contact.Address = v
	}
}

// Config 定义用于读取 SMTP 等配置
type Config struct {
	SMTP struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		User      string `yaml:"user"`
		Pass      string `yaml:"pass"`
		ContactTo string `yaml:"contact_to"`
	} `yaml:"smtp"`
	Contact struct {
		Phone   string `yaml:"phone"`
		Email   string `yaml:"email"`
		Wechat  string `yaml:"wechat"`
		Address string `yaml:"address"`
	} `yaml:"contact"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
