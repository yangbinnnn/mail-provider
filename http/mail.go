package http

import (
	"crypto/tls"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/jordan-wright/email"
	"github.com/toolkits/web/param"
	"github.com/yangbinnnn/mail-provider/config"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)
		auth := smtp.PlainAuth("", cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Addr)
		msg := &email.Email{
			To:      strings.Split(tos, ","),
			From:    cfg.Smtp.From,
			Subject: subject,
			Text:    []byte(content),
			Headers: textproto.MIMEHeader{},
		}
		var err error
		if cfg.Smtp.TLS {
			err = msg.SendWithTLS(cfg.Smtp.Addr, auth, &tls.Config{})
		} else {
			err = msg.Send(cfg.Smtp.Addr, auth)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
