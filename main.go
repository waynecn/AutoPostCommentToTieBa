package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Println("main")
	Sleeped = false
	InitRandomContent()
	ScanTieBa()
}

func GetTbs() TBSResult {
	client := &http.Client{Timeout: 5 * time.Second}

	var tbsResult TBSResult
	req, err := http.NewRequest("GET", tbsUrl, nil)
	if err != nil {
		log.Println("创建tbs请求错误:", err)
		tbsResult.Error = err
		return tbsResult
	}

	req.Header.Add("Cookie", cookies)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("请求tbs错误:", err)
		tbsResult.Error = err
		return tbsResult
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("解析tbs请求结果错误:", err)
		tbsResult.Error = err
		return tbsResult
	}

	err = json.Unmarshal(bts, &tbsResult)
	if err != nil {
		log.Println("json结构化tbs请求结果错误:", err)
		tbsResult.Error = err
		return tbsResult
	}
	currentTbs = tbsResult.Tbs
	return tbsResult
}

func GetAllBars() LIKEDBars {
	var likeBars LIKEDBars
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", likedTiebaUrl, nil)
	if err != nil {
		log.Println("创建关注的吧请求错误:", err)
		likeBars.Error2 = err
		return likeBars
	}

	req.Header.Add("Cookie", cookies)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("请求关注的吧错误:", err)
		likeBars.Error2 = err
		return likeBars
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("解析关注的吧请求结果错误:", err)
		likeBars.Error2 = err
		return likeBars
	}

	err = json.Unmarshal(bts, &likeBars)
	if err != nil {
		log.Println("json结构化关注的吧请求结果错误:", err)
		likeBars.Error2 = err
		return likeBars
	}
	return likeBars
}

//获取到所有关注的吧名、forumId
func ScanTieBa() {
	GetTbs()

	tieBarIndex := 0
	exitFlag := false
	for !exitFlag {
		allBars := GetAllBars()
		fmt.Println("贴吧列表:")
		for index, item := range allBars.Data.LikeForum {
			fmt.Println("吧名:", item.ForumName, " 序号:", index+1)
		}

		if tieBarIndex >= len(allBars.Data.LikeForum) {
			tieBarIndex = 0
		}

		kw := strings.Replace(allBars.Data.LikeForum[tieBarIndex].ForumName, "+", "%2B", -1)
		chooseKw = kw
		chooseFid = strconv.Itoa(allBars.Data.LikeForum[tieBarIndex].ForumId)
		tieBarUrl := "https://tieba.baidu.com/f?kw=" + kw

		fmt.Println("进入贴吧:", kw)
		ShowTopics(tieBarUrl)

		tieBarIndex++
		if !Sleeped {
			time.Sleep(20 * time.Second)
			Sleeped = true
		}
	}
}

//找出所有话题
func ShowTopics(strUrl string) {
	client := &http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		fmt.Println("http创建NewRequest发生错误:", err)
		return
	}

	request.Header.Add("Cookie", cookies)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("请求话题列表发生错误:", err)
		return
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("话题列表:", string(bts))
	f, err := os.OpenFile("tiebacontent.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("打开文件发生错误:", err)
		return
	}

	f.Write(bts)
	f.Close()

	f, err = os.OpenFile("tiebacontent.txt", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("打开文件发生错误:", err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var hrefs []string
	var titles []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "href=\"/p") != -1 {
			//该行是标题所在行
			startIndex := strings.Index(line, "href=\"/p")
			line = line[startIndex+6:]
			//fmt.Println("line 2:", line)
			endIndex := strings.Index(line, "\"")
			if -1 == endIndex {
				fmt.Println("Href end index not found")
				break
			}
			href := line[:endIndex]
			//fmt.Println("Href:", href)

			line = line[endIndex+1:]
			startIndex = strings.Index(line, "title=")
			line = line[startIndex+7:]
			endIndex = strings.Index(line, "\"")
			title := line[:endIndex]

			hrefs = append(hrefs, href)
			titles = append(titles, title)

			topicUrl := "https://tieba.baidu.com" + href
			chooseTitle = title
			fmt.Println("自动进入话题链接:", topicUrl, "  话题:", title)
			GoToTopic(topicUrl)

			if !Sleeped {
				time.Sleep(20 * time.Second)
				Sleeped = true
			}
		}
	}
}

//进入话题
func GoToTopic(strUrl string) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		fmt.Println("NewRequest发生错误:", err)
		return
	}

	req.Header.Add("Cookie", cookies)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求话题详情发生错误:", err)
		return
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("解析响应体发生错误:", err)
		return
	}

	f, err := os.OpenFile("topicDetail.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		return
	}
	f.Write(bts)
	f.Close()

	f, err = os.OpenFile("topicDetail.txt", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Count(line, "class=\"p_forbidden_post_content_unfold") <= 1 {
			//表示只有主楼 可以抢占2楼
			AutoReplyTopic(strUrl)

			if !Sleeped {
				time.Sleep(20 * time.Second)
			}
			Sleeped = false
			break
		}

		//这一行找出楼层内回复
		index := strings.Index(line, "class=\"p_forbidden_post_content_unfold")
		if index == -1 {
			continue
		}
		index = strings.Index(line, "style=\"display:;\">")
		if -1 == index {
			continue
		}

		line = line[index+7:]
		//这一行匹配到具体回复内容
		index = strings.Index(line, "style=\"display:;\">")
		if -1 == index {
			continue
		}

		line = line[index+7:]
		index = strings.Index(line, ">")
		if -1 == index {
			continue
		}
		line = line[index+1:]

		index2 := strings.Index(line, "</div>")
		if -1 == index2 {
			continue
		}

		content := line[index:index2]
		fmt.Println("话题链接:", strUrl)
		fmt.Println("话题内容:", content)

		// ReplyTopic(strUrl)

		// breakFlag := false
		// for true {
		// 	fmt.Println("是否继续回复?0退出，其他继续")
		// 	breakStr := ""
		// 	fmt.Scanln(&breakStr)
		// 	if "0" == breakStr {
		// 		breakFlag = true
		// 		break
		// 	}
		// 	break
		// }
		// if breakFlag {
		// 	break
		// }
	}
}

//回复话题
func ReplyTopic(urlStr string) {
	lastSlashIndex := strings.LastIndex(urlStr, "/")
	tid := urlStr[lastSlashIndex+1:]

	var content string
	fmt.Println("输入你要回复的内容:")
	fmt.Scanln(&content)

	payloadStr := GenerateFormDataPayload(tid, content)
	bodyData := bytes.NewBuffer([]byte(payloadStr))

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", addCotentUrl, bodyData)
	if err != nil {
		fmt.Println("回复时创建请求发生错误:", err)
		return
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", "tieba.baidu.com")
	req.Header.Add("Origin", "https://tieba.baidu.com")
	req.Header.Add("Referer", urlStr)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("回复出错:", err)
		return
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("解析响应失败:", err)
		return
	}
	var commentResponse CommentResponse
	err = json.Unmarshal(bts, &commentResponse)
	if err != nil {
		fmt.Println("序列化响应内容失败:", err)
		return
	}
	fmt.Println("响应内容:", string(bts))
	if commentResponse.ErrCode == 0 {
		fmt.Println("回复成功")
	} else {
		fmt.Println("回复失败---，原因:", commentResponse.Error)
	}
}

func AutoReplyTopic(urlStr string) {
	lastSlashIndex := strings.LastIndex(urlStr, "/")
	tid := urlStr[lastSlashIndex+1:]

	content := RandomContent[CurrentLine]
	CurrentLine++
	if len(RandomContent) <= CurrentLine {
		CurrentLine = 0
	}

	fmt.Println("自动回复内容:", content)

	payloadStr := GenerateFormDataPayload(tid, content)
	bodyData := bytes.NewBuffer([]byte(payloadStr))

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", addCotentUrl, bodyData)
	if err != nil {
		fmt.Println("回复时创建请求发生错误:", err)
		return
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", "tieba.baidu.com")
	req.Header.Add("Origin", "https://tieba.baidu.com")
	req.Header.Add("Referer", urlStr)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("回复出错:", err)
		return
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("解析响应失败:", err)
		return
	}
	var commentResponse CommentResponse
	err = json.Unmarshal(bts, &commentResponse)
	if err != nil {
		fmt.Println("序列化响应内容失败:", err)
		return
	}
	fmt.Println("响应内容:", string(bts))
	if commentResponse.ErrCode == 0 {
		fmt.Println("回复成功")
	} else {
		fmt.Println("回复失败---，原因:", commentResponse.Error)
	}
}

func GenerateFormDataPayload(tid string, content string) string {
	u := url.Values{}
	u.Set("ie", "utf-8")
	u.Set("kw", chooseKw)
	u.Set("fid", chooseFid)
	u.Set("tid", tid)
	u.Set("vcode_md5", "")
	u.Set("floor_num", "100000")
	u.Set("rich_text", "1")
	u.Set("tbs", currentTbs)
	u.Set("content", content)
	u.Set("basilisk", "1")
	u.Set("files", "")
	u.Set("mouse_pwd", "29,16,28,4,25,24,24,17,28,33,25,4,24,4,25,4,24,4,25,4,24,4,25,4,24,4,25,4,24,33,27,31,25,17,17,33,25,17,26,24,4,25,24,16,24,16445594424401")
	u.Set("mouse_pwd_t", "1644559442440")
	u.Set("mouse_pwd_isclick", "1")
	u.Set("__type__", "reply")
	u.Set("nick_name", NickName)
	u.Set("ev", "comment")
	u.Set("geetest_success", "0")
	u.Set("_BSK", _BSK)
	u.Set("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	u.Set("biz[po]", BizPo)

	return u.Encode()
}
