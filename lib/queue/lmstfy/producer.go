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

type Feishu struct {
	Tos []string `json:"tos"`
	Sub string   `json:"sub"`
	Con string   `json:"con"`
}

func Ping() (err error) {
	c := GetInstance()
	_, err = c.Client.Publish("test", []byte("this is a test message"), 5, 1, 0)
	return
}

// Mail相关的生产函数
func ProduceCallbackMail(mail Mail) (err error) {
	m := new(bytes.Buffer)
	err = json.NewEncoder(m).Encode(mail)
	if err != nil {
		log.Logger.Error("failed to create pool account event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_callback_mail_event", m.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new callback mail event: ", jobId)
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

// 飞书相关的生产函数
func ProduceCallbackFeishu(feishu Feishu) (err error) {
	m := new(bytes.Buffer)
	err = json.NewEncoder(m).Encode(feishu)
	if err != nil {
		log.Logger.Error("failed to create callback feishu event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_callback_feishu_event", m.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new callback feishu event: ", jobId)
	return
}

func ProduceProblemFeishu(feishu Feishu) (err error) {
	m := new(bytes.Buffer)
	err = json.NewEncoder(m).Encode(feishu)
	if err != nil {
		log.Logger.Error("failed to create problem feishu event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_problem_feishu_event", m.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new problem feishu event: ", jobId)
	return
}

func ProduceOKFeishu(feishu Feishu) (err error) {
	m := new(bytes.Buffer)
	err = json.NewEncoder(m).Encode(feishu)
	if err != nil {
		log.Logger.Error("failed to create ok feishu event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_ok_feishu_event", m.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new ok feishu event: ", jobId)
	return
}
