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
		job, err := c.Client.Consume("new_problem_mail_event", 10, 2)
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
		req, err := http.NewRequest("GET", "http://127.0.0.1:4000/sender/mail", nil)
		q := req.URL.Query()
		q.Add("tos", mail.Tos)
		q.Add("subject", mail.Sub)
		q.Add("content", mail.Con)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		resp.Body.Close()

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
		job, err := c.Client.Consume("new_ok_mail_event", 10, 2)
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
			log.Logger.Error("failed to unmarshal create ok mai event job data")
			continue
		}

		fmt.Println(mail.Con)
		req, err := http.NewRequest("GET", "http://127.0.0.1:4000/sender/mail", nil)
		q := req.URL.Query()
		q.Add("tos", mail.Tos)
		q.Add("subject", mail.Sub)
		q.Add("content", mail.Con)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		resp.Body.Close()

		err1 := c.Client.Ack("new_ok_mail_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}
