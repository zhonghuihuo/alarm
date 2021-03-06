package sender

import (
	"log"
	"strings"
	"time"

	"github.com/toolkits/smtp"

	"github.com/urlooker/alarm/g"
)

func ConsumeMail() {
	queue := g.Config.Queue.Mail
	for {
		L := PopAllMail(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendMailList(L)
	}
}

func SendMailList(L []*g.Mail) {
	for _, mail := range L {
		if mail.Tos == "" || mail.Tos == "," || mail.Tos == ";" || mail.Content == "" {
			continue
		}

		toArr := strings.Split(mail.Tos, ",")
		log.Println("MailCount", len(toArr))

		MailWorkerChan <- 1
		go SendMail(mail)
	}
}

func SendMail(mail *g.Mail) {
	defer func() {
		<-MailWorkerChan
	}()

        if g.Config.Smtp.TLS {
                s := smtp.NewSMTP(g.Config.Smtp.Addr, g.Config.Smtp.Username, g.Config.Smtp.Password, g.Config.Smtp.TLS, g.Config.Smtp.Anonymous, g.Config.Smtp.SkipVerify)
                err := s.SendMail(g.Config.Smtp.From, strings.Replace(mail.Tos, ",", ";", -1), mail.Subject, mail.Content, "text")
                if err != nil {
                        log.Println(err, "tos:", mail.Tos)
                        //SendSmsToMaintainer("sender:" + err.Error())
                }
                if g.Config.Debug {
                        resp := "ok"
                        if err != nil {
                                resp = err.Error()
                        }
                        log.Println("==mail==>>>>", mail)
                        log.Println("<<<<==mail==", resp)
                }
        } else {
                s := smtp.New(g.Config.Smtp.Addr, g.Config.Smtp.Username, g.Config.Smtp.Password)
                err := s.SendMail(g.Config.Smtp.From, strings.Replace(mail.Tos, ",", ";", -1), mail.Subject, mail.Content, "text")
                if err != nil {
                        log.Println(err, "tos:", mail.Tos)
                        //SendSmsToMaintainer("sender:" + err.Error())
                }
                if g.Config.Debug {
                        resp := "ok"
                        if err != nil {
                                resp = err.Error()
                        }
                        log.Println("==mail==>>>>", mail)
                        log.Println("<<<<==mail==", resp)
                }
        }
}
