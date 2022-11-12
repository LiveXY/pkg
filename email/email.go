package email

import (
	"crypto/tls"
	"etms/pkg/bytex"
	"etms/pkg/logx"
	"etms/pkg/strx"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

// SMTP配置
type SMTPConfig struct {
	Host     string `yaml:"host"`     // SMTP主机
	Account  string `yaml:"account"`  // 账号
	Password string `yaml:"password"` // 密码
	Display  string `yaml:"display"`  // 显示名
	From     string `yaml:"from"`     // 来自
}

// 发送邮件
func SendEmail(cfg SMTPConfig, to, subject, body string) (err error) {
	if len(cfg.Host) < 5 || len(cfg.Account) < 5 || len(cfg.Password) < 5 {
		return nil
	}
	if !strings.Contains(cfg.Host, ":") {
		cfg.Host = cfg.Host + ":25"
	}
	host := strings.Split(cfg.Host, ":")[0]
	port := strings.Split(cfg.Host, ":")[1]
	auth := smtp.PlainAuth("", cfg.Account, cfg.Password, host)
	contenttype := "Content-Type: text/plain; charset=UTF-8"
	display, from1, from2 := "", "", ""
	if len(cfg.Display) > 0 {
		display = cfg.Display + " "
	}
	if len(cfg.From) > 0 {
		if len(display) > 0 {
			from2 = cfg.From
			from1 = display + "<" + cfg.From + ">"
		} else {
			from2 = cfg.From
			from1 = cfg.From
		}
	} else {
		if len(display) > 0 {
			from2 = cfg.Account
			from1 = display + "<" + cfg.Account + ">"
		} else {
			from2 = cfg.Account
			from1 = cfg.Account
		}
	}
	msg := strx.ToBytes("To: " + to + "\r\nFrom: " + from1 + "\r\nSubject: " + subject + "\r\n" + contenttype + "\r\n\r\n" + body)
	sendto := strings.Split(to, ";")
	logx.Logger.Debug("发送邮件：", zap.String("from", from2), zap.String("to", to), zap.String("msg", bytex.ToStr(msg)), zap.Any("config", cfg))
	if port == "465" {
		conn, err := tls.Dial("tcp", cfg.Host, nil)
		if err != nil {
			logx.Error.Error("发送邮件错误dial：", zap.String("to", to), zap.String("body", body), zap.Error(err))
			return err
		}
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			logx.Error.Error("发送邮件错误new client：", zap.String("to", to), zap.String("body", body), zap.Error(err))
			return err
		}
		defer c.Close()
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				logx.Error.Error("发送邮件错误auth：", zap.String("to", to), zap.String("body", body), zap.Error(err))
				return err
			}
		}
		if err = c.Mail(from2); err != nil {
			logx.Error.Error("发送邮件错误mail：", zap.String("to", to), zap.String("body", body), zap.Error(err))
			return err
		}
		for _, addr := range sendto {
			if err = c.Rcpt(addr); err != nil {
				logx.Error.Error("发送邮件错误to：", zap.String("to", to), zap.String("body", body), zap.Error(err))
				return err
			}
		}
		w, err := c.Data()
		if err != nil {
			logx.Error.Error("发送邮件错误data：", zap.String("to", to), zap.String("body", body), zap.Error(err))
			return err
		}
		defer w.Close()
		if _, err = w.Write(msg); err != nil {
			logx.Error.Error("发送邮件错误write：", zap.String("to", to), zap.String("body", body), zap.Error(err))
			return err
		}
		c.Quit()
		err = nil
	} else {
		err := smtp.SendMail(cfg.Host, auth, from2, sendto, msg)
		if err != nil {
			logx.Error.Error("发送邮件错误：", zap.String("to", to), zap.String("body", body), zap.Error(err))
			return err
		}
	}
	logx.Logger.Debug("发送邮件成功：", zap.String("to", to), zap.String("body", body))
	return err
}
