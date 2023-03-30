package lmstfy

import (
	"bytes"
	"encoding/json"
	"falcon-mail-transmit/lib/log"
)

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func Ping() (err error) {
	c := GetInstance()
	_, err = c.Client.Publish("test", []byte("this is a test message"), 5, 1, 0)
	return
}

func ProduceMail(params Mail) (err error) {
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(params)
	if err != nil {
		log.Logger.Error("failed to create pool account event")
		return
	}
	c := GetInstance()
	jobId, err := c.Client.Publish("new_create_pool_account_event", b.Bytes(), 0, 3, 5)
	log.Logger.Info("lmstfy new create pool account event: ", jobId)
	return
}
