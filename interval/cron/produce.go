package cron

import (
	"encoding/json"
	"falcon-mail-transmit/interval/config"
	"falcon-mail-transmit/lib/log"
	"falcon-mail-transmit/lib/queue/lmstfy"
	"falcon-mail-transmit/lib/redis"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ProduceFromRedis(cfg *config.Config) {
	// 目前采取一分钟一次的频率
	now := time.Now()
	next := now.Add(time.Minute * 1)
	//minute := next.Minute() - (next.Minute() % 5)
	next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
	timer := time.NewTimer(next.Sub(now))

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("produce from redis panic: ", err)
		}
	}()
	for {
		select {
		case <-timer.C:
			result, err := redis.SetNX("cron_produce_mail", true, time.Second*10)
			if err != nil {
				log.Logger.Error(err.Error())
			}
			if result {
				// 邮件队列
				problemMailKeys := redis.Keys("problem-mail-*")
				okMailKeys := redis.Keys("ok-mail-*")
				callbackMailKeys := redis.Keys("callback-mail-*")

				// 飞书队列
				problemFeishuKeys := redis.Keys("problem-feishu-*")
				okFeishuKeys := redis.Keys("ok-feishu-*")
				callbackFeishuKeys := redis.Keys("callback-feishu-*")

				// 处理邮件消息
				okMergedMails := GetMergedMails(okMailKeys, cfg.MailTypes, "恢复")
				problemMails := GetMergedMails(problemMailKeys, cfg.MailTypes, "故障")
				callbackMails := GetMergedMails(callbackMailKeys, cfg.MailTypes, "回调函数")

				fmt.Println("okMergedMails:", okMergedMails)
				fmt.Println("problemMails:", problemMails)
				fmt.Println("callbackMails:", callbackMails)

				for _, problemMail := range problemMails {
					err = lmstfy.ProduceProblemMail(problemMail)
					if err != nil {
						log.Logger.Error("failed to produce.go create new pool account event")
						continue
					}
				}
				for _, okMergedMail := range okMergedMails {
					err = lmstfy.ProduceOKMail(okMergedMail)
					if err != nil {
						log.Logger.Error("failed to produce.go create new pool account event")
						continue
					}
				}
				for _, callbackMail := range callbackMails {
					err = lmstfy.ProduceCallbackMail(callbackMail)
					if err != nil {
						log.Logger.Error("failed to produce.go create new pool account event")
						continue
					}
				}

				// 处理飞书消息
				okMergedFeishu := GetMergedFeishu(okFeishuKeys, cfg.MailTypes, "恢复")
				problemFeishu := GetMergedFeishu(problemFeishuKeys, cfg.MailTypes, "故障")
				callbackFeishu := GetMergedFeishu(callbackFeishuKeys, cfg.MailTypes, "回调函数")

				fmt.Println("okMergedFeishu:", okMergedFeishu)
				fmt.Println("problemFeishu:", problemFeishu)
				fmt.Println("callbackFeishu:", callbackFeishu)

				for _, feishu := range problemFeishu {
					err = lmstfy.ProduceProblemFeishu(feishu)
					if err != nil {
						log.Logger.Error("failed to produce problem feishu event")
						continue
					}
				}
				for _, feishu := range okMergedFeishu {
					err = lmstfy.ProduceOKFeishu(feishu)
					if err != nil {
						log.Logger.Error("failed to produce ok feishu event")
						continue
					}
				}
				for _, feishu := range callbackFeishu {
					err = lmstfy.ProduceCallbackFeishu(feishu)
					if err != nil {
						log.Logger.Error("failed to produce callback feishu event")
						continue
					}
				}
			}
		}
		now := time.Now()
		next := now.Add(time.Minute * 1)
		//minute := next.Minute() - (next.Minute() % 5)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		timer.Reset(next.Sub(time.Now()))
	}
}

type MergedMails struct {
	Sub string `json:"sub"`
}

func GetMergedMails(MailKeys []string, mailTypes []string, subjectType string) []lmstfy.Mail {
	var mails []lmstfy.Mail
	for _, MailKey := range MailKeys {
		mapMail := redis.HGetAll(MailKey)
		//fmt.Println("mapMail:", mapMail)
		var mail lmstfy.Mail
		jsonMail, err := json.Marshal(mapMail)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		//fmt.Println("jsonMail:", string(jsonMail))
		err = json.Unmarshal(jsonMail, &mail)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		//fmt.Println("Mail:", Mail)
		mails = append(mails, mail)
		err = redis.Del(MailKey)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
	}
	//fmt.Println("Mails:", Mails)

	//构造一个map，key为mail_type，value为[]Mail
	MailMap := make(map[string][]lmstfy.Mail)
	for _, t := range mailTypes {
		MailMap[t] = make([]lmstfy.Mail, 0)
	}

	var mergedMails []lmstfy.Mail

	//再把struct做筛选
	//根据mail_type合并相似的内容，再存入一个新的[]struct
	//如果不属于任何一类，则直接塞入mergedMails内
	for _, mail := range mails {
		isExistMailType := 0
		for mailType, _ := range MailMap {
			//fmt.Println(mailType)
			//fmt.Println(value)
			if strings.Contains(mailType, "*") {
				parts := strings.Split(mailType, "*")
				match := true
				for _, part := range parts {
					fmt.Println("part：", part)
					if !strings.Contains(mail.Sub, part) {
						match = false
						// 匹配不到mailType中的part时直接中断循环，从而减少代码运行的时间
						break
					}
				}
				if match {
					isExistMailType = 1
					MailMap[mailType] = append(MailMap[mailType], mail)
				}
			} else if strings.Contains(mail.Sub, mailType) {
				isExistMailType = 1
				MailMap[mailType] = append(MailMap[mailType], mail)
			}
		}
		if isExistMailType == 0 {
			mergedMails = append(mergedMails, mail)
		}
	}

	for mailType, mailTypeMails := range MailMap {
		if len(mailTypeMails) == 0 {
			continue
		}
		conAll := ""
		for i, mail := range mailTypeMails {
			conAll = conAll + "第" + strconv.Itoa(i+1) + "封:\n" + mail.Con + "\n"
		}
		//fmt.Println("conAll:", conAll)
		mergedMails = append(mergedMails, lmstfy.Mail{
			Tos: mailTypeMails[0].Tos,
			Sub: "[合并邮件][" + subjectType + "]" + mailType,
			Con: conAll,
		})
	}
	//fmt.Println("mergedMails:", mergedMails)

	return mergedMails
}

// GetMergedFeishu 处理飞书消息的合并逻辑
func GetMergedFeishu(feishuKeys []string, messageTypes []string, subjectType string) []lmstfy.Feishu {
	var feishuMessages []lmstfy.Feishu

	// 从 Redis 中读取所有飞书消息
	for _, feishuKey := range feishuKeys {
		mapFeishu := redis.HGetAll(feishuKey)
		var feishu lmstfy.Feishu
		jsonFeishu, err := json.Marshal(mapFeishu)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		err = json.Unmarshal(jsonFeishu, &feishu)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		feishuMessages = append(feishuMessages, feishu)
		err = redis.Del(feishuKey)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
	}

	// 构造一个map，key为message_type，value为[]Feishu
	feishuMap := make(map[string][]lmstfy.Feishu)
	for _, t := range messageTypes {
		feishuMap[t] = make([]lmstfy.Feishu, 0)
	}

	var mergedFeishu []lmstfy.Feishu

	// 根据message_type合并相似的内容
	for _, feishu := range feishuMessages {
		isExistMessageType := 0
		for messageType, _ := range feishuMap {
			if strings.Contains(messageType, "*") {
				parts := strings.Split(messageType, "*")
				match := true
				for _, part := range parts {
					if !strings.Contains(feishu.Sub, part) {
						match = false
						break
					}
				}
				if match {
					isExistMessageType = 1
					feishuMap[messageType] = append(feishuMap[messageType], feishu)
				}
			} else if strings.Contains(feishu.Sub, messageType) {
				isExistMessageType = 1
				feishuMap[messageType] = append(feishuMap[messageType], feishu)
			}
		}
		if isExistMessageType == 0 {
			mergedFeishu = append(mergedFeishu, feishu)
		}
	}

	// 合并同类型的飞书消息
	for messageType, typeFeishu := range feishuMap {
		if len(typeFeishu) == 0 {
			continue
		}
		contentAll := ""
		for i, feishu := range typeFeishu {
			contentAll = contentAll + "**第" + strconv.Itoa(i+1) + "条:**\n" + feishu.Con + "\n\n"
		}
		mergedFeishu = append(mergedFeishu, lmstfy.Feishu{
			Tos: typeFeishu[0].Tos,
			Sub: "[合并消息][" + subjectType + "]" + messageType,
			Con: contentAll,
		})
	}

	return mergedFeishu
}
