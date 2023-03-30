package main

import (
	"fmt"
	"github.com/bitleak/lmstfy/client"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
)

type QueueLmstfyClient struct {
	Client *client.LmstfyClient
}

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func main() {
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

	fmt.Println(con)
	req, err := http.NewRequest("GET", "http://127.0.0.1:4000/sender/mail", nil)
	q := req.URL.Query()
	q.Add("tos", tos)
	q.Add("subject", sub)
	q.Add("content", con)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	return c.String(http.StatusOK, "success")
}