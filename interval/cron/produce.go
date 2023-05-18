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
				// 先全部从redis里取出来，转换为struct
				problemMailKeys := redis.Keys("problem-mail-*")
				okMailKeys := redis.Keys("ok-mail-*")
				callbackMailKeys := redis.Keys("callback-mail-*")
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
