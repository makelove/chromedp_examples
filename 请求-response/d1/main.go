/**
* @File    :   main.go
* @Time    :   2021/09/04 10:55:50
* @Author  :   GH
* @Desc    :    基本完成
如何使用chromedp获取HTTP响应主体？
https://ask.csdn.net/questions/1016423
https://stackoverflow.com/questions/45808799/how-to-get-the-http-response-body-using-chromedp

参考
https://github.com/chromedp/chromedp/issues/409
*/
// TODO 断点调试

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.WindowSize(1920, 1080),
		chromedp.UserDataDir(dir),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 10*time.Second)
	defer cancel()
	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}
	// listen network event
	listenForNetworkEvent(taskCtx)
	chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(`https://space.bilibili.com/180948619/video`),
		chromedp.WaitVisible(`body`, chromedp.BySearch),
	)
}
func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				// log.Printf("received headers: %s", resp.Headers)
				// log.Printf("received headers: %s", resp.URL)
				// // log.Printf("received headers: %s", resp)
				// // TODO 打印请求网站
				// println("------------------------")

				// 2021/09/04 11:38:18 received headers: &{https://api.bilibili.com/x/emote/user/panel/web?business=reply %!s(int64=200)  map[access-control-allow-credentials:true access-control-allow-headers:Origin,No-Cache,X-Requested-With,If-Modified-Since,Pragma,Last-Modified,Cache-Control,Expires,Content-Type,Access-Control-Allow-Credentials,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Cache-Webcdn,x-bilibili-key-real-ip,x-backend-bili-real-ip access-control-allow-origin:https://space.bilibili.com bili-status-code:0 bili-trace-id:4df453ae086132ea content-encoding:br content-type:application/json; charset=utf-8 date:Sat, 04 Sep 2021 03:38:18 GMT x-cache-webcdn:BYPASS from blzone03]  application/json map[]  %!s(bool=true) %!s(float64=184) 110.43.34.184 %!s(int64=443) %!s(bool=false) %!s(bool=false) %!s(bool=false) %!s(float64=413) %!s(*network.ResourceTiming=&{7289.527079 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 0.401 0.648 0 0 103.324})  %!s(*cdp.TimeSinceEpoch=&{145224192 52912224763 0x1951f60})  h2 secure %!s(*network.SecurityDetails=&{TLS 1.2 ECDHE_RSA P-256 AES_128_GCM  0 *.bilibili.com [*.bilibili.com bilibili.com] GlobalSign RSA OV SSL CA 2018 0xc0002f15c0 0xc0002f15e0 [0xc00039aa80 0xc00039ab00 0xc00039ab80] compliant})}

				if strings.Contains(resp.URL, "space/arc/search") {
					log.Printf("received : %s , %s", resp.URL, ev.RequestID)

					// how to modify the response Header/Body ?
					go func() {
						// print response body
						c := chromedp.FromContext(ctx)
						rbp := network.GetResponseBody(ev.RequestID)
						body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
						if err != nil {
							fmt.Println(err)
						}
						//写入json文件
						if err = ioutil.WriteFile(ev.RequestID.String()+".json", body, 0644); err != nil {
							log.Fatal(err)
						}

						if err == nil {
							fmt.Printf("%s\n", body)
						}
					}()
					/**
					2021/09/05 19:44:00 received : https://api.bilibili.com/x/space/arc/search?mid=180948619&ps=30&tid=0&pn=1&keyword=&order=pubdate&jsonp=jsonp , 4391.69

					{"code":0,"message":"0","ttl":1,"data":{"list":{"tlist":{"1":{"tid":1,"count":5,"name":"动画"},"160":{"tid":160,"count":68,"name":"生活"},"181":{"tid":181,"count":4,"name":"影视"},"188":{"tid":188,"count":144,"name":"科技"},"211":{"tid":211,"count":2,"name":"美食"},"217":{"tid":217,"count":6,"name":"动物圈"},"223":{"tid":223,"count":12,"name":"汽车"},"234":{"tid":234,"count":8,"name":"运动"},"36":{"tid":36,"count":327,"name":"知识"},"4":{"tid":4,"count":1,"name":"游戏"},"5":{"tid":5,"count":1,"name":"娱乐"}},"vlist":[{"comment":0,"typeid":231,"play":24,"pic":"http://i1.hdslb.com/bfs/archive/ccc3a03c8ea938a782f73d96eb20ce9ae46f7527.jpg","subtitle":"","description":"https://www.youtube.com/watch?v=Qf07DxKBuKU\n在页面注入js，文本到语音tts\n有点意思","copyright":"2","title":"【编程】无头浏览器自动化Headless Chrome and browser automation  with Eric Bidelman","review":0,"author":"程序员怎样挣钱","mid":180948619,"created":1630832450,"length":"46:37","video_review":0,"aid":975266147,"bvid":"BV1444y187mb","hide_click":false,"is_pay":0,"is_union_video":0,"is_steins_gate":0,"is_live_playback":0}
					*/
				}
			}
		}
		// other needed network Event
	})
}
