# httpHelper
下载工具，方便配置代理，headers,简单的获取网页html内容，具体可以看demo.go

    

###   获取网站html内容
    func GetHtmlDemo() {
      dst_url := "https://www.google.com"
      proxy_str := "socks5://127.0.0.1:1081"
      defaultCookieJar, _ := cookiejar.New(nil)
      var headerMap map[string]string = map[string]string{
        "User-Agent": GetRandomUserAgent(),
        "Referer":    "http://www.baidu.com",
      }
      //httpGet(dst_url string, proxy_str string, headerMap map[string]string, defaultCookieJar *http.CookieJar)

      _, response, _, err := HttpGet(dst_url, proxy_str, headerMap, defaultCookieJar)
      if err != nil {
        log.Fatal(err)
      }
      defer response.Body.Close()
      b, err := ioutil.ReadAll(response.Body)
      if err != nil {
        log.Fatal(err)
      }
      log.Println(string(b))

    }

###  下载文件
    func DownDemo() {
      proxy_str := "socks5://127.0.0.1:1081"
      var downBean DownEntity
      downBean.Down_url = "http://doshantianfang1.zgpingshu.com/%E5%8D%95%E7%94%B0%E8%8A%B3%E8%AF%84%E4%B9%A6_%E9%9A%8B%E5%94%90%E6%BC%94%E4%B9%89216%E5%9B%9E%E7%89%88216%E5%9B%9E1.13GB_32k/300043529F.mp3"
      downBean.Save_Root = "D:/down/go/zgps/"
      downBean.Save_name = "a.mp3"
      //	downBean.Save_path = downBean.Save_Root + downBean.Save_name
      downBean.Save_path = "a.mp3"
      downBean.Cor = 10
      downBean.Headers = map[string]string{
        "User-Agent":      "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:48.0) Gecko/20100101 Firefox/48.0",
        "Referer":         "http://www.zgpingshu.com/down/1040/",
        "Accept-Encoding": "gzip",
      }
      MultithreadDownloader(downBean, proxy_str)
    }


## 下载模块使用了@WindGreen 兄弟写的多线程下载模块，https://github.com/WindGreen/fileload 感谢额
