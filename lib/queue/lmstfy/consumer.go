package lmstfy

import (
	"encoding/json"
	"falcon-mail-transmit/lib/log"
	"fmt"
	"net/http"
	"time"
)

func ConsumeNewProblemMailEvent() {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("consume queue panic: ", err)
		}
	}()

	c := GetInstance()
	for {
		job, err := c.Client.Consume("new_problem_mail_event", 10, 5)
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		if job == nil {
			log.Logger.Info("no new_problem_mail_event job, continue")
			continue
		}
		log.Logger.Info("consume new_problem_mail_event, id: ", job.ID)

		var mail Mail
		err = json.Unmarshal(job.Data, &mail)
		if err != nil {
			log.Logger.Error("failed to unmarshal create problem mail event job data")
			continue
		}

		fmt.Println(mail.Con)
		req, err := http.NewRequest("POST", "http://127.0.0.1:4000/sender/mail", nil)
		q := req.URL.Query()
		q.Add("tos", mail.Tos)
		q.Add("subject", mail.Sub)
		q.Add("content", mail.Con)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			if err != nil {
				log.Logger.Error("err:", err)
			}
			log.Logger.Error("failed to send mail by main email, try to send by backup email")

			req2, err := http.NewRequest("POST", "http://127.0.0.1:4001/sender/mail", nil)
			q2 := req2.URL.Query()
			q2.Add("tos", mail.Tos)
			q2.Add("subject", "[主邮箱错误]"+mail.Sub)
			q2.Add("content", mail.Con)
			req2.URL.RawQuery = q2.Encode()
			resp2, err := http.DefaultClient.Do(req2)
			if err != nil || resp2.StatusCode != 200 {
				if err != nil {
					log.Logger.Error(err)
				}
				log.Logger.Error("failed to send mail by backup email!!!")
				continue
			} else {
				log.Logger.Info("success send mail by backup email success")
			}
			resp2.Body.Close()
		} else {
			log.Logger.Info("send mail by main email success")
		}
		if resp != nil {
			resp.Body.Close()
		}

		err1 := c.Client.Ack("new_problem_mail_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}

func ConsumeNewOKMailEvent() {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("consume queue panic: ", err)
		}
	}()

	c := GetInstance()
	for {
		job, err := c.Client.Consume("new_ok_mail_event", 10, 5)
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		if job == nil {
			log.Logger.Info("no new_ok_mail_event job, continue")
			continue
		}
		log.Logger.Info("consume new_ok_mail_event, id: ", job.ID)

		var mail Mail
		err = json.Unmarshal(job.Data, &mail)
		if err != nil {
			log.Logger.Error("failed to unmarshal create ok mail event job data")
			continue
		}

		fmt.Println(mail.Con)
		req, err := http.NewRequest("POST", "http://127.0.0.1:4000/sender/mail", nil)
		q := req.URL.Query()
		q.Add("tos", mail.Tos)
		q.Add("subject", mail.Sub)
		q.Add("content", mail.Con)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			if err != nil {
				log.Logger.Error("err:", err)
			}
			log.Logger.Error("failed to send mail by main email, try to send by backup email")

			req2, err := http.NewRequest("POST", "http://127.0.0.1:4001/sender/mail", nil)
			q2 := req2.URL.Query()
			q2.Add("tos", mail.Tos)
			q2.Add("subject", "[主邮箱错误]"+mail.Sub)
			q2.Add("content", mail.Con)
			req2.URL.RawQuery = q2.Encode()
			resp2, err := http.DefaultClient.Do(req2)
			if err != nil || resp2.StatusCode != 200 {
				if err != nil {
					log.Logger.Error(err)
				}
				log.Logger.Error("failed to send mail by backup email!!!")
				continue
			} else {
				log.Logger.Info("success send mail by backup email success")
			}
			resp2.Body.Close()
		} else {
			log.Logger.Info("send mail by main email success")
		}
		if resp != nil {
			resp.Body.Close()
		}

		err1 := c.Client.Ack("new_ok_mail_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}

func ConsumeNewCallbackMailEvent() {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("consume queue panic: ", err)
		}
	}()

	c := GetInstance()
	for {
		job, err := c.Client.Consume("new_callback_mail_event", 10, 5)
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		if job == nil {
			log.Logger.Info("no new_callback_mail_event job, continue")
			continue
		}
		log.Logger.Info("consume new_callback_mail_event, id: ", job.ID)

		var mail Mail
		err = json.Unmarshal(job.Data, &mail)
		if err != nil {
			log.Logger.Error("failed to unmarshal create callback mail event job data")
			continue
		}

		fmt.Println(mail.Con)
		req, err := http.NewRequest("POST", "http://127.0.0.1:4000/sender/mail", nil)
		q := req.URL.Query()
		q.Add("tos", mail.Tos)
		q.Add("subject", mail.Sub)
		q.Add("content", mail.Con)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			if err != nil {
				log.Logger.Error("err:", err)
			}
			log.Logger.Error("failed to send mail by main email, try to send by backup email")
			req2, err := http.NewRequest("POST", "http://127.0.0.1:4001/sender/mail", nil)
			q2 := req2.URL.Query()
			q2.Add("tos", mail.Tos)
			q2.Add("subject", "[主邮箱错误]"+mail.Sub)
			q2.Add("content", mail.Con)
			req2.URL.RawQuery = q2.Encode()
			resp2, err := http.DefaultClient.Do(req2)
			if err != nil || resp2.StatusCode != 200 {
				if err != nil {
					log.Logger.Error(err)
				}
				log.Logger.Error("failed to send mail by backup email!!!")
				continue
			}
			resp2.Body.Close()
		} else {
			log.Logger.Info("send mail by main email success")
		}
		if resp != nil {
			resp.Body.Close()
		}

		err1 := c.Client.Ack("new_callback_mail_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}
