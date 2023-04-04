package lmstfy

import (
	"bytes"
	"encoding/json"
	"falcon-mail-transmit/lib/log"
)

type Mail struct {
	Tos string `json:"tos"`
	Sub string `json:"sub"`
	Con string `json:"con"`
}

func Ping() (err error) {
	c := GetInstance()
	_, err = c.Client.Publish("test", []byte("this is a test message"), 5, 1, 0)
	return
}

func ProduceProblemMail(mail Mail) (err error) {
	m := new(bytes.Buffer)
	err = json.NewEncoder(m).Encode(mail)
	if err != nil {
		log.Logger.Error("failed to create pool account event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_problem_mail_event", m.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new problem mail event: ", jobId)
	return
}

func ProduceOKMail(mail Mail) (err error) {
	m := new(bytes.Buffer)
	err = json.NewEncoder(m).Encode(mail)
	if err != nil {
		log.Logger.Error("failed to create pool account event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_ok_mail_event", m.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new ok mail event: ", jobId)
	return
}
