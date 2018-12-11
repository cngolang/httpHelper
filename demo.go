package httpHelper

import (
	"io/ioutil"
	"log"
	"net/http/cookiejar"
)

func GetHtmlDemo() {
	dst_url := "https://www.zhihu.com"
	proxy_str := "http://127.0.0.1:8888"
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

/*

package main

import (
	"log"
	"net/http/cookiejar"

	"io/ioutil"

	"./httpHelper"
)

func main() {
	//	dst_url := "http://appstore.kenxinda.com"
	dst_url := "https://www.google.com"
	proxy_str := "socks5://127.0.0.1:1081"
	defaultCookieJar, _ := cookiejar.New(nil)
	var headerMap map[string]string = map[string]string{
		"User-Agent": httpHelper.GetRandomUserAgent(),
		"Referer":    "http://www.baidu.com",
	}
	//httpGet(dst_url string, proxy_str string, headerMap map[string]string, defaultCookieJar *http.CookieJar)

	_, response, _, err := httpHelper.HttpGet(dst_url, proxy_str, headerMap, defaultCookieJar)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(b))

	//	var downBean httpHelper.DownEntity
	//	downBean.Down_url = "http://doshantianfang1.zgpingshu.com/%E5%8D%95%E7%94%B0%E8%8A%B3%E8%AF%84%E4%B9%A6_%E9%9A%8B%E5%94%90%E6%BC%94%E4%B9%89216%E5%9B%9E%E7%89%88216%E5%9B%9E1.13GB_32k/300043529F.mp3"
	//	downBean.Save_Root = "D:/down/go/zgps/"
	//	downBean.Save_name = "a.mp3"
	//	downBean.Save_path = downBean.Save_Root + downBean.Save_name
	//	downBean.Cor = 2
	//	downBean.Headers = map[string]string{
	//		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:48.0) Gecko/20100101 Firefox/48.0",
	//		"Referer":         "http://www.zgpingshu.com/down/1040/",
	//		"Accept-Encoding": "gzip",
	//	}

	//	httpHelper.MultithreadDownloader(downBean, proxy_str)

}

*/
