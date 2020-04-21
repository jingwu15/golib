package clientold

import (
	"fmt"
	"github.com/jingwu15/golib/time"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	Uuid    string            `json:"uuid"` //请求的唯一标识符
	Ext     map[string]string `json:"ext"`  //扩展数据
}

type Response struct {
	ReqTime  string            `json:"req_time"`
	Httpcode int               `json:"httpcode"`
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Error    error             `json:"error"`
}

type ReqRes struct {
	Req Request  `json:"req"`
	Res Response `json:"res"`
}

type ReqRess struct {
	Req   Request    `json:"req"`
	Times int        `json:"times"`
	Ress  []Response `json:"ress"`
}

type SeqReqOkOne struct {
	ReqResss []ReqRess
	Stats    map[string]int
}

type RRsStats struct {
	Rrs   []ReqRes
	Stats map[string]int
}

type CheckReqOk func(interface{}, ReqRes) (bool, error)

func NewRequest() Request {
	return Request{
		Method:  "",
		Url:     "",
		Body:    "",
		Headers: map[string]string{},
		Uuid:    uuid.NewV4().String(),
		Ext:     map[string]string{},
	}
}

func NewReqRess() ReqRess {
	return ReqRess{
		Req: Request{
			Method:  "",
			Url:     "",
			Body:    "",
			Headers: map[string]string{},
			Uuid:    uuid.NewV4().String(),
			Ext:     map[string]string{},
		},
		Times: 0,
		Ress:  []Response{},
	}
}

//根据url取得host
func GetHostByUrl(raw string) string {
	u, _ := url.Parse(raw)
	return u.Host
}

//根据url取得host, 无端口
func GetHostNoPortByUrl(raw string) string {
	u, _ := url.Parse(raw)
	tmp := strings.Split(u.Host, ":")
	return tmp[0]
}

func HeadersToStr(headers map[string][]string) string {
	lines := []string{}
	for key, rows := range headers {
		lines = append(lines, fmt.Sprintf("%s: %s", key, strings.Join(rows, " ")))
	}
	return strings.Join(lines, "\r\n")
}

func HeadersMapSliceToMapStr(headers map[string][]string) map[string]string {
	headersN := map[string]string{}
	for key, rows := range headers {
		headersN[key] = strings.Join(rows, " ")
	}
	return headersN
}

func DoPost(url, reqBody string, headers map[string]string, timeout int) (httpcode int, body string, respHeaders map[string]string, err error) {
	client := &http.Client{Timeout: time.Keep(timeout)}
    fmt.Println("reqBody", reqBody)
    fmt.Println("headers", headers)

	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	if err != nil {
		return 0, "", map[string]string{}, err
	}

	//使用短连接
	req.Close = true
	req.Header.Set("Connection", "close")
	//host设置比较不同
	if host, ok := headers["Host"]; ok {
		req.Host = host
	}
	if host, ok := headers["host"]; ok {
		req.Host = host
	}

	//设用用户的自定义头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return 0, "", map[string]string{}, err
	}

	bodyRaw, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return response.StatusCode, string(bodyRaw), HeadersMapSliceToMapStr(response.Header), err
	}
	return response.StatusCode, string(bodyRaw), HeadersMapSliceToMapStr(response.Header), nil
}

//使用ReqRess进行请求，用于多次重试并记录每次的请求
func DoPostReqRess(reqRess ReqRess, timeout int) ReqRess {
	res := Response{ReqTime: time.Now().ToStr()}
	res.Httpcode, res.Body, res.Headers, res.Error = DoPost(reqRess.Req.Url, reqRess.Req.Body, reqRess.Req.Headers, timeout)
	reqRess.Times++
	reqRess.Ress = append(reqRess.Ress, res)
	return reqRess
}

//请求，并将结果写入到通道
func DoPostReqRessToChan(reqRess ReqRess, timeout int, ch chan ReqRess) {
	res := Response{ReqTime: time.Now().ToStr()}
	res.Httpcode, res.Body, res.Headers, res.Error = DoPost(reqRess.Req.Url, reqRess.Req.Body, reqRess.Req.Headers, timeout)
	reqRess.Times++
	reqRess.Ress = append(reqRess.Ress, res)
	ch <- reqRess
}

//批量请求，并统计成功与失败个数
func MultiPostReqress(reqRessA1 []ReqRess, timeout int) []ReqRess {
	length := len(reqRessA1)
	ch := make(chan ReqRess, len(reqRessA1))
	for _, reqRess := range reqRessA1 {
		go DoPostReqRessToChan(reqRess, timeout, ch)
	}

	reqRessA1N := []ReqRess{}
	//接收数据
	total := 0
	for {
		select {
		case row, _ := <-ch:
			reqRessA1N = append(reqRessA1N, row)
			total++
			if total == length {
				goto GoRun
			}
		case <-time.After(timeout):
			goto GoRun
			break
		}
	}
GoRun:
	return reqRessA1N
}

//请求，并将结果写入到通道
func DoPostToChan(req Request, timeout int, ch chan ReqRes) {
	rr := ReqRes{Req: req}
	rr.Res.Httpcode, rr.Res.Body, rr.Res.Headers, rr.Res.Error = DoPost(req.Url, req.Body, req.Headers, timeout)
	rr.Res.ReqTime = time.Now().ToStr()
	ch <- rr
}

//批量请求，并统计成功与失败个数
func MultiPost(reqs []Request, timeout int) ([]ReqRes, map[string]int) {
	length := len(reqs)
	ch := make(chan ReqRes, len(reqs))
	for _, req := range reqs {
		go DoPostToChan(req, timeout, ch)
	}

	stats := map[string]int{
		"total":       length,
		"total_ok":    0,
		"total_error": 0,
	}
	httpcodes := []int{}
	rrs := []ReqRes{}
	//接收数据
	total := 0
	for {
		select {
		case row, _ := <-ch:
			rrs = append(rrs, row)
			if row.Res.Error == nil {
				stats["total_ok"]++
				httpcodes = append(httpcodes, row.Res.Httpcode)
			} else {
				stats["total_error"]++
			}
			total++
			if total == length {
				goto GoRun
			}
		case <-time.After(timeout):
			goto GoRun
			break
		}
	}
GoRun:
	for _, httpcode := range httpcodes {
		key := fmt.Sprintf("code_%d", httpcode)
		if _, ok := stats[key]; !ok {
			stats[key] = 0
		}
		stats[key]++
	}
	return rrs, stats
}

//顺序请求，只要成功一个即返回
func PostSeqOkOne(reqs []Request, timeout int, checkok CheckReqOk, okMust interface{}) ([]ReqRes, map[string]int) {
	statsReqOk := 0
	statsReqErr := 0
	rrs := []ReqRes{}
	for _, req := range reqs {
		rr := ReqRes{Req: req}
		rr.Res.Httpcode, rr.Res.Body, rr.Res.Headers, rr.Res.Error = DoPost(req.Url, req.Body, req.Headers, timeout)
		rr.Res.ReqTime = time.Now().ToStr()
		rrs = append(rrs, rr)
		ok, _ := checkok(okMust, rr)
		if ok {
			statsReqOk++
			break
		} else {
			statsReqErr++
		}
	}
	return rrs, map[string]int{"total": len(reqs), "total_req": len(rrs), "total_req_ok": statsReqOk, "total_req_err": statsReqErr}
}

func PostSeqOkOneToChan(reqs []Request, timeout int, checkok CheckReqOk, okMust interface{}, ch chan RRsStats) {
	rrsStats := RRsStats{}
	rrsStats.Rrs, rrsStats.Stats = PostSeqOkOne(reqs, timeout, checkok, okMust)
	ch <- rrsStats
}

func MultiPostSeqOkOne(reqss [][]Request, timeout int, checkok CheckReqOk, okMust interface{}) ([]RRsStats, map[string]int) {
	length := len(reqss)
	ch := make(chan RRsStats, len(reqss))
	for _, reqs := range reqss {
		go PostSeqOkOneToChan(reqs, timeout, checkok, okMust, ch)
	}

	stats := map[string]int{
		"total":       length,
		"total_ok":    0,
		"total_error": 0,
	}
	rrsStatsM := []RRsStats{}
	//接收数据
	total := 0
	for {
		select {
		case row, _ := <-ch:
			rrsStatsM = append(rrsStatsM, row)
			if row.Stats["total_ok"] == 1 {
				stats["total_ok"]++
			} else {
				stats["total_error"]++
			}
			total++
			if total == length {
				goto GoRun
			}
		case <-time.After(timeout):
			//goto GoRun
			//break
		}
	}
GoRun:
	return rrsStatsM, stats
}

//顺序请求，只要成功一个即返回
func ReqRessPostSeqOkOne(reqRessA1 []ReqRess, timeout int, checkok CheckReqOk, okMust interface{}) ([]ReqRess, map[string]int) {
	statsReqOk := 0
	statsReqErr := 0
	reqRessA1N := []ReqRess{}
	for _, reqRess := range reqRessA1 {
		req := reqRess.Req
		reqRessN := DoPostReqRess(reqRess, timeout)
		res := reqRessN.Ress[len(reqRessN.Ress)-1]
		reqRessA1N = append(reqRessA1N, reqRessN)
		rr := ReqRes{Req: req, Res: res}
		ok, _ := checkok(okMust, rr)
		if ok {
			statsReqOk++
			break
		} else {
			statsReqErr++
		}
	}
	return reqRessA1N, map[string]int{"total": len(reqRessA1), "total_req": len(reqRessA1N), "total_req_ok": statsReqOk, "total_req_err": statsReqErr}
}

func ReqRessPostSeqOkOneToChan(reqRessA1 []ReqRess, timeout int, checkok CheckReqOk, okMust interface{}, ch chan SeqReqOkOne) {
	seqReqOkOne := SeqReqOkOne{}
	seqReqOkOne.ReqResss, seqReqOkOne.Stats = ReqRessPostSeqOkOne(reqRessA1, timeout, checkok, okMust)
	ch <- seqReqOkOne
}

func ReqRessMultiPostSeqOkOne(reqRessA2 [][]ReqRess, timeout int, checkok CheckReqOk, okMust interface{}) []SeqReqOkOne {
	length := len(reqRessA2)
	ch := make(chan SeqReqOkOne, len(reqRessA2))
	for _, reqRessA1 := range reqRessA2 {
		go ReqRessPostSeqOkOneToChan(reqRessA1, timeout, checkok, okMust, ch)
	}

	stats := map[string]int{
		"total":       length,
		"total_ok":    0,
		"total_error": 0,
	}
	seqReqOkOneA1 := []SeqReqOkOne{}
	//接收数据
	total := 0
	for {
		select {
		case row, _ := <-ch:
			seqReqOkOneA1 = append(seqReqOkOneA1, row)
			if row.Stats["total_ok"] == 1 {
				stats["total_ok"]++
			} else {
				stats["total_error"]++
			}
			total++
			if total == length {
				goto GoRun
			}
		case <-time.After(timeout):
			goto GoRun
			break
		}
	}
GoRun:
	return seqReqOkOneA1
}
