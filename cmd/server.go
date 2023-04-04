package main

import (
	"falcon-mail-transmit/interval/config"
	"falcon-mail-transmit/interval/cron"
	"falcon-mail-transmit/interval/utils"
	log2 "falcon-mail-transmit/lib/log"
	"falcon-mail-transmit/lib/queue/lmstfy"
	"falcon-mail-transmit/lib/redis"
	"flag"
	"fmt"
	"github.com/bitleak/lmstfy/client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	math "math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type QueueLmstfyClient struct {
	Client *client.LmstfyClient
}

var debugMode bool
var logInfo *rotatelogs.RotateLogs

func main() {
	parseFlag()
	cfg, err := config.Load()
	if err != nil {
		log.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	logInfo = log2.Init(cfg.LogDir)
	if debugMode {
		log2.Logger.SetLevel(logrus.DebugLevel)
	}

	err = lmstfy.Ping()
	if err != nil {
		log2.Logger.Fatal("failed to publish lmstfy test message: ", err)
		os.Exit(-1)
	}

	go cron.ProduceFromRedis(cfg)

	utils.HandlePanic(lmstfy.ConsumeNewProblemMailEvent, lmstfy.ConsumeNewOKMailEvent)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "falcon-mail-transmit ( ￣ー￣)人(^▽^ )")
	})
	e.POST("/sender/mail", save)
	e.Logger.Fatal(e.Start(":1323"))
}

func save(c echo.Context) error {
	tos := c.FormValue("tos")
	sub := c.FormValue("subject")
	con := c.FormValue("content")

	fmt.Println("tos:", tos)
	fmt.Println("subject:", sub)
	fmt.Println("content:", con)

	rSub1 := regexp.MustCompile("\\[P0\\]")
	rSub2p := regexp.MustCompile("\\[PROBLEM\\]")
	rSub2o := regexp.MustCompile("\\[OK\\]")
	rSub3 := regexp.MustCompile("\\[\\]")
	rSub4 := regexp.MustCompile("。[^\\]]*\\]")
	rSub5 := regexp.MustCompile("\\[[^\\]]*\\]$")

	sub = rSub1.ReplaceAllString(sub, "[海外矿池]")
	sub = rSub2p.ReplaceAllString(sub, "[故障]")
	sub = rSub2o.ReplaceAllString(sub, "[恢复]")
	sub = rSub3.ReplaceAllString(sub, "")
	sub = rSub4.ReplaceAllString(sub, "]")
	sub = rSub5.ReplaceAllString(sub, "")
	fmt.Println(sub)

	rCon1 := regexp.MustCompile("\nP0")
	rCon2p := regexp.MustCompile("PROBLEM")
	rCon2o := regexp.MustCompile("OK")
	rCon3 := regexp.MustCompile("\nEndpoint")
	rCon4 := regexp.MustCompile("\nMetric")
	rCon5 := regexp.MustCompile("\nTags")
	rCon6 := regexp.MustCompile("\nNote")
	rCon7m := regexp.MustCompile("\nMax")
	rCon7c := regexp.MustCompile(", Current")
	rCon8 := regexp.MustCompile("\nTimestamp")

	con = rCon1.ReplaceAllString(con, "\n海外矿池")
	con = rCon2p.ReplaceAllString(con, "故障")
	con = rCon2o.ReplaceAllString(con, "恢复")
	con = rCon3.ReplaceAllString(con, "\n服务器")
	con = rCon4.ReplaceAllString(con, "\n触发函数")
	con = rCon5.ReplaceAllString(con, "\n子标签")
	con = rCon6.ReplaceAllString(con, "\n故障原因")
	con = rCon7m.ReplaceAllString(con, "\n最多通知次数")
	con = rCon7c.ReplaceAllString(con, ", 目前故障次数")
	con = rCon8.ReplaceAllString(con, "\n故障时间")

	//mail := lmstfy.Mail{
	//	Tos:     tos,
	//	Subject: sub,
	//	Content: con,
	//}

	// 在这里加上一个堵塞队列，全部存到redis里，每分钟从redis里取出一次，塞到消息队列里

	// 生成redisKey，防止重复
	redisKey := ""
	for {
		if strings.Contains(sub, "故障") {
			redisKey = "problem-mail-" + GenerateNoWithoutPrefix()
		} else {
			redisKey = "ok-mail-" + GenerateNoWithoutPrefix()
		}

		m := redis.HGetAll(redisKey)
		if len(m) == 0 {
			break
		}
	}

	err := redis.HSet(redisKey, "tos", tos)
	if err != nil {
		log2.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = redis.HSet(redisKey, "sub", sub)
	if err != nil {
		log2.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = redis.HSet(redisKey, "con", con)
	if err != nil {
		log2.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//if strings.Contains(sub, "故障") {
	//	err := lmstfy.ProduceProblemMail(mail)
	//	if err != nil {
	//		log2.Logger.Error("failed to produce.go create new pool account event")
	//		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	//	}
	//} else {
	//	err := lmstfy.ProduceOKMail(mail)
	//	if err != nil {
	//		log2.Logger.Error("failed to produce.go create new pool account event")
	//		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	//	}
	//}

	//fmt.Println(con)
	//req, err := http.NewRequest("GET", "http://127.0.0.1:4000/sender/mail", nil)
	//q := req.URL.Query()
	//q.Add("tos", tos)
	//q.Add("subject", sub)
	//q.Add("content", con)
	//req.URL.RawQuery = q.Encode()
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//defer resp.Body.Close()
	return c.String(http.StatusOK, "success")
}

func parseFlag() {
	flag.BoolVar(&debugMode, "debug", false, "set log level")
	flag.Parse()
}

func GenerateNoWithoutPrefix() string {
	date := time.Now().Format("20060102")
	ts := time.Now().UnixNano() / 1e6 % 1e5
	r := math.Intn(1000)
	no := fmt.Sprintf("%s%07d%03d", date[2:], ts, r)
	return no
}
