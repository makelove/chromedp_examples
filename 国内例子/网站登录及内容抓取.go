/**
* @File    :   网站登录及内容抓取.js
* @Time    :   2021/08/08 20:51:03
* @Author  :   GH
* @Desc    :

go语言下用chromedp框架编写爬虫程序实现网站登录及内容抓取
https://blog.csdn.net/peihexian/article/details/104436496
*/

package main

import(
    "context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main()  {
    var buf []byte

    // create chrome instance
    ctx, cancel := chromedp.NewContext(
        context.Background(),
        chromedp.WithDebugf(log.Printf),
    )
    defer cancel()

    // create a timeout
    ctx, cancel = context.WithTimeout(ctx, 50 * time.Second)
    defer cancel()

    // run task list
    var res string
    var err error

    width, height := 1920, 1080
    err = chromedp.Run(ctx, chromedp.Tasks{
        emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false),
    })

    loginUrl:= `http://a.b.c.d/login`

    var executed * runtime.RemoteObject
    username:= "yourusername"
    password:= "yourpassword"
    err = chromedp.Run(ctx, chromedp.Tasks{
        chromedp.Navigate(loginUrl),
        chromedp.Sleep(5 * time.Second),
        chromedp.Evaluate(`var jq = document.createElement('script'); jq.src = "https://cdn.bootcss.com/jquery/1.4.2/jquery.js"; document.getElementsByTagName('head')[0].appendChild(jq);`,& executed),
        chromedp.Sleep(5 * time.Second),
        chromedp.WaitVisible(`#mypassword`, chromedp.ByID),
        chromedp.SendKeys(`input[name="username"]`, username, chromedp.NodeVisible),
        chromedp.SendKeys(`#mypassword`, password, chromedp.ByID),
        chromedp.Sleep(2 * time.Second),
        chromedp.Click(`#login_btn`, chromedp.ByID),
        chromedp.Sleep(5 * time.Second),
        chromedp.CaptureScreenshot(& buf),
    })

    if err != nil {
        log.Fatal(err)
    }
    if err := ioutil.WriteFile("1.png", buf, 0644); err != nil {
        log.Fatal(err)
    }

    indexPageUrl:= `http://a.b.c.d/somepage`

    err = chromedp.Run(ctx, chromedp.Tasks{
        chromedp.Navigate(indexPageUrl),
        chromedp.Sleep(5 * time.Second),
        chromedp.CaptureScreenshot(& buf),
        chromedp.WaitVisible(`#somehtmlid`, chromedp.ByID),
    })
    if err != nil {
        log.Fatal(err)
    }
    if err := ioutil.WriteFile("2.png", buf, 0644); err != nil {
        log.Fatal(err)
    }

    log.Printf("got: `%s`", strings.TrimSpace(res))

}
