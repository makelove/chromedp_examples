/**
* @File    :   main.go
* @Time    :   2021/09/04 10:55:50
* @Author  :   GH
* @Desc    :   
如何使用chromedp获取HTTP响应主体？
https://ask.csdn.net/questions/1016423
https://stackoverflow.com/questions/45808799/how-to-get-the-http-response-body-using-chromedp

*/
// TODO 断点调试

package main
import (
    "context"
    "io/ioutil"
    "log"
    "os"
    "time"
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
        chromedp.Flag("window-size", "500,400"),
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
        chromedp.Navigate(`https://space.bilibili.com/180948619/dynamic`),
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
                log.Printf("received headers: %s", resp)
				resp.urlstr // TODO 打印请求网站
				println("------------------------")
				`
				2021/09/04 11:38:18 received headers: &{https://api.bilibili.com/x/emote/user/panel/web?business=reply %!s(int64=200)  map[access-control-allow-credentials:true access-control-allow-headers:Origin,No-Cache,X-Requested-With,If-Modified-Since,Pragma,Last-Modified,Cache-Control,Expires,Content-Type,Access-Control-Allow-Credentials,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Cache-Webcdn,x-bilibili-key-real-ip,x-backend-bili-real-ip access-control-allow-origin:https://space.bilibili.com bili-status-code:0 bili-trace-id:4df453ae086132ea content-encoding:br content-type:application/json; charset=utf-8 date:Sat, 04 Sep 2021 03:38:18 GMT x-cache-webcdn:BYPASS from blzone03]  application/json map[]  %!s(bool=true) %!s(float64=184) 110.43.34.184 %!s(int64=443) %!s(bool=false) %!s(bool=false) %!s(bool=false) %!s(float64=413) %!s(*network.ResourceTiming=&{7289.527079 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 0.401 0.648 0 0 103.324})  %!s(*cdp.TimeSinceEpoch=&{145224192 52912224763 0x1951f60})  h2 secure %!s(*network.SecurityDetails=&{TLS 1.2 ECDHE_RSA P-256 AES_128_GCM  0 *.bilibili.com [*.bilibili.com bilibili.com] GlobalSign RSA OV SSL CA 2018 0xc0002f15c0 0xc0002f15e0 [0xc00039aa80 0xc00039ab00 0xc00039ab80] compliant})}
				`
            }
        }
        // other needed network Event
    })
}