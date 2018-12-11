// fileload
package httpHelper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/gob"
	_ "flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	//	"golang.org/x/net/proxy"
)

const VERSION = "1.0.3"

var queue, redo, finish chan int
var cor, size, length, timeout int
var hash, dst string
var verify, version, cache bool
var downentity DownEntity

//cnum 线程数，savePath 保存路径 durl 下载路径
//func MultithreadDownloader(cor int, pdst string, durl string) {
func MultithreadDownloader(downbean DownEntity, proxy_str string) {
	downentity = downbean
	cor = downentity.Cor
	dst = downentity.Save_path
	url := downentity.Down_url

	if version || url == "version" {
		fmt.Println("Fileload version:", VERSION)
		return
	}

	if verify {
		file, err := os.Open(url)
		if err != nil {
			log.Println(err)
			return
		}
		if hash == "sha1" {
			h := sha1.New()
			io.Copy(h, file)
			r := h.Sum(nil)
			log.Printf("sha1 of file: %x\n", r)
		} else if hash == "md5" {
			h := md5.New()
			io.Copy(h, file)
			r := h.Sum(nil)
			log.Printf("sha1 of file: %x\n", r)
		}
		return
	}

	if dst == "" {
		_, dst = filepath.Split(url)
	}

	startTime := time.Now()
	//	var client *http.Client

	//	proxy_str := "socks5://127.0.0.1:8889"
	client := GetClient(proxy_str, 0, nil, false)
	request, err := http.NewRequest("GET", url, nil)
	HttpHeaderBuild(request, downentity.Headers)

	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	response.Body.Close()
	num := response.Header.Get("Content-Length")
	length, _ = strconv.Atoi(num)

	if size <= 0 {
		size = int(math.Ceil(float64(length) / float64(cor)))
	}
	fragment := int(math.Ceil(float64(length) / float64(size)))
	queue = make(chan int, cor)
	redo = make(chan int, int(math.Floor(float64(cor)/2)))
	go func() {
		for i := 0; i < fragment; i++ {
			queue <- i
		}
		for {
			j := <-redo
			queue <- j
		}
	}()
	finish = make(chan int, cor)
	for j := 0; j < cor; j++ {
		go Do(request, fragment, j)
	}
	for k := 0; k < fragment; k++ {
		_ = <-finish
	}

	file, err := os.Create(dst)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	var offset int64 = 0
	for x := 0; x < fragment; x++ {
		filename := fmt.Sprintf("%s_%d", dst, x)
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		file.WriteAt(buf, offset)
		offset += int64(len(buf))
		os.Remove(filename)
	}
	log.Println("Written to ", dst)
	//hash
	if hash == "sha1" {
		h := sha1.New()
		io.Copy(h, file)
		r := h.Sum(nil)
		log.Printf("sha1 of file: %x\n", r)
	} else if hash == "md5" {
		h := md5.New()
		io.Copy(h, file)
		r := h.Sum(nil)
		log.Printf("sha1 of file: %x\n", r)
	}

	finishTime := time.Now()
	duration := finishTime.Sub(startTime).Seconds()
	log.Printf("Time:%f Speed:%f Kb/s\n", duration, float64(length)/duration/1024)
}

func Do(request *http.Request, fragment, no int) {
	var req http.Request
	err := DeepCopy(&req, request)
	if err != nil {
		log.Println("ERROR|prepare request:", err)
		log.Panic(err)
		return
	}
	for {
		i := <-queue
		start := i * size
		var end int
		if i < fragment-1 {
			end = start + size - 1
		} else {
			end = length - 1
		}

		filename := fmt.Sprintf("%s_%d", dst, i)
		if cache {
			filesize := int64(end - start + 1)
			file, err := os.Stat(filename)
			if err == nil && file.Size() == filesize {
				log.Printf("[%d][%d]Hint cached %s, size:%d\n", no, i, filename, filesize)
				finish <- i
				continue
			}
		}

		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
		proxy_str := "socks5://127.0.0.1:8889"
		cli := GetClient(proxy_str, timeout, nil, false)
		resp, err := cli.Do(&req)
		if err != nil {
			log.Printf("[%d][%d]ERROR|do request:%s\n", no, i, err.Error())
			redo <- i
			continue
		}

		file, err := os.Create(filename)
		if err != nil {
			log.Printf("[%d][%d]ERROR|create file %s:%s\n", no, i, filename, err.Error())
			file.Close()
			resp.Body.Close()
			redo <- i
			continue
		}
		n, err := io.Copy(file, resp.Body)
		if err != nil {
			log.Printf("[%d][%d]ERROR|write to file %s:%s  -tmp %s\n", no, i, filename, err.Error(), n)
			file.Close()
			resp.Body.Close()
			redo <- i
			continue
		}

		file.Close()
		resp.Body.Close()
		finish <- i
	}
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
