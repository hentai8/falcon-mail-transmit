package lmstfy

import (
	"context"
	"encoding/json"
	"falcon-mail-transmit/interval/config"
	"falcon-mail-transmit/lib/log"
	"fmt"
	"net/http"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// Mail消费函数
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

// 飞书消费函数
func ConsumeNewProblemFeishuEvent() {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("consume feishu queue panic: ", err)
		}
	}()

	c := GetInstance()
	for {
		job, err := c.Client.Consume("new_problem_feishu_event", 10, 5)
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		if job == nil {
			log.Logger.Info("no new_problem_feishu_event job, continue")
			continue
		}
		log.Logger.Info("consume new_problem_feishu_event, id: ", job.ID)

		var feishu Feishu
		err = json.Unmarshal(job.Data, &feishu)
		if err != nil {
			log.Logger.Error("failed to unmarshal problem feishu event job data")
			continue
		}

		err = sendFeishuMessage(feishu)
		if err != nil {
			log.Logger.Error("failed to send feishu message: ", err)
			continue
		}

		err1 := c.Client.Ack("new_problem_feishu_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}

func ConsumeNewOKFeishuEvent() {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("consume feishu queue panic: ", err)
		}
	}()

	c := GetInstance()
	for {
		job, err := c.Client.Consume("new_ok_feishu_event", 10, 5)
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		if job == nil {
			log.Logger.Info("no new_ok_feishu_event job, continue")
			continue
		}
		log.Logger.Info("consume new_ok_feishu_event, id: ", job.ID)

		var feishu Feishu
		err = json.Unmarshal(job.Data, &feishu)
		if err != nil {
			log.Logger.Error("failed to unmarshal ok feishu event job data")
			continue
		}

		err = sendFeishuMessage(feishu)
		if err != nil {
			log.Logger.Error("failed to send feishu message: ", err)
			continue
		}

		err1 := c.Client.Ack("new_ok_feishu_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}

func ConsumeNewCallbackFeishuEvent() {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("consume feishu queue panic: ", err)
		}
	}()

	c := GetInstance()
	for {
		job, err := c.Client.Consume("new_callback_feishu_event", 10, 5)
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		if job == nil {
			log.Logger.Info("no new_callback_feishu_event job, continue")
			continue
		}
		log.Logger.Info("consume new_callback_feishu_event, id: ", job.ID)

		var feishu Feishu
		err = json.Unmarshal(job.Data, &feishu)
		if err != nil {
			log.Logger.Error("failed to unmarshal callback feishu event job data")
			continue
		}

		err = sendFeishuMessage(feishu)
		if err != nil {
			log.Logger.Error("failed to send feishu message: ", err)
			continue
		}

		err1 := c.Client.Ack("new_callback_feishu_event", job.ID)
		if err1 != nil {
			log.Logger.Error("failed to ack job")
			continue
		}
	}
}

// 使用飞书官方SDK发送消息
func sendFeishuMessage(feishu Feishu) error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// 创建飞书 Client
	client := lark.NewClient(cfg.Feishu.AppID, cfg.Feishu.AppSecret)

	// 构造消息内容
	// 使用富文本格式
	content := fmt.Sprintf(`{
		"zh_cn": {
			"title": "%s",
			"content": [
				[
					{
						"tag": "text",
						"text": "%s"
					}
				]
			]
		}
	}`, feishu.Sub, feishu.Con)

	for _, tos := range feishu.Tos {
		// 创建请求对象
		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType("user_id").
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(tos).
				MsgType("post"). // 使用富文本格式
				Content(content).
				Build()).
			Build()

		// 发起请求
		resp, err := client.Im.V1.Message.Create(context.Background(), req)

		// 处理错误
		if err != nil {
			log.Logger.Error("failed to send to ", tos, ": ", err)

		}

		// 服务端错误处理
		if !resp.Success() {
			log.Logger.Error("failed to send to ", tos, ": ", resp.Msg)
			continue
		}

		log.Logger.Info("✓ Sent to ", tos)
	}

	return nil
}
