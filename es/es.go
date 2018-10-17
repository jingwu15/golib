package es

import (
//	"fmt"
	"strings"
	"strconv"
//	"reflect"
	"encoding/json"
	jparse "github.com/buger/jsonparser"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

var apiMap = map[string]map[string]string{
	"getfp_time" : {
		"uri": "tjkd-app-#date#/_search?filter_path=took,aggregations.fps.buckets",
		"req": `{
			"query": {
				"bool": {
					"filter": {"range":{"Timestamp": {"gte": "#start_time#", "lte": "#end_time#"}}}
				}
			},
			"sort":[],
			"aggs":{
			  "fps":{
			    "terms":{"field": "fingerprint"}
			  }
			}}`,
	},
	"getfp_ip" : {
		"uri": "tjkd-app-#date#/_search?filter_path=took,aggregations.fps.buckets",
		"req": `{
			"query": {
				"bool": {
					"must": {"match": {"remote_addr": "#remote_addr#"}},
					"filter": {"range":{"Timestamp": {"gte": "#start_time#", "lte": "#end_time#"}}}
				}
			},
			"sort":[],
			"aggs":{
			  "fps":{
			    "terms":{"field": "fingerprint"}
			  }
			}}`,
	},
	"count-gethistory_fp" : {
		"uri": "tjkd-app-#date#/_count?filter_path=count",
		"req": `{
          "query": {
            "bool": {
              "must":{
                "match":{"fingerprint": "#fingerprint#"}
              },
              "filter": {
                "range":{"Timestamp": {"gte": "#start_time#", "lte": "#end_time#"}}
              }
            }
          }}`,
	},
	"gethistory_fp" : {
		"uri": "tjkd-app-#date#/_search?filter_path=hits.hits._source",
		"req": `{
            "_source":{
              "includes":[
			  	"continent","country","city","upstream_addr","isp","Timestamp","platform",
				"session_time","manufacturer","tcp_connect_time","accesskey","province",
				"client_connect_time","fingerprint","brand","remote_addr","simulator", "version",
				"server_addr","upstream_session_time","upstream_connect_time","geo_val","device"],
              "excludes":[]
            },
          "query": {
            "bool": {
              "must":{
                "match":{"fingerprint": "#fingerprint#"}
              },
              "filter": {
                "range":{"Timestamp": {"gte": "#start_time#", "lte": "#end_time#"}}
              }
            }
          },
          "from":#from#,
          "size":#size#,
          "sort":{"Timestamp": {"order": "asc"}},
          "aggs":{}}`,
	},
	"count-gethistory_ip" : {
		"uri": "tjkd-app-#date#/_count?filter_path=count",
		"req": `{
          "query": {
            "bool": {
              "must":{
                "match":{"remote_addr": "#remote_addr#"}
              },
              "filter": {
                "range":{"Timestamp": {"gte": "#start_time#", "lte": "#end_time#"}}
              }
            }
          }}`,
	},
	"gethistory_ip" : {
		"uri": "tjkd-app-#date#/_search?filter_path=hits.hits._source",
		"req": `{
          "_source":{
            "includes":[
		    	"continent","country","city","upstream_addr","isp","Timestamp","platform",
		  	"session_time","manufacturer","tcp_connect_time","accesskey","province",
		  	"client_connect_time","fingerprint","brand","remote_addr","simulator", "version",
		  	"server_addr","upstream_session_time","upstream_connect_time","geo_val","device"],
            "excludes":[]
          },
          "query": {
            "bool": {
              "must":{
                "match":{"remote_addr": "#remote_addr#"}
              },
              "filter": {
                "range":{"Timestamp": {"gte": "#start_time#", "lte": "#end_time#"}}
              }
            }
          },
          "from":#from#,
          "size":#size#,
          "sort":{"Timestamp": {"order": "asc"}},
          "aggs":{}}`,
	},
}

func PostES(apikey string, params map[string]string) ([]byte, map[string][]string, error) {
	var err error
	var body []byte
	var response *http.Response
	url := viper.GetString("es_api") + apiMap[apikey]["uri"]
	url = strings.Replace(url, "#date#", params["date"], -1)
	reqdata := apiMap[apikey]["req"]
	for k,v := range params {
		reqdata = strings.Replace(reqdata, "#"+k+"#", v, -1)
	}
	response, err = http.Post(url, "application/json", strings.NewReader(reqdata))
	if err != nil {
		return nil,nil,err
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil,nil,err
	}
	return body, response.Header, nil
}

/**
	var params = map[string]string{
		"date": "2018-04-02",
		"start_time": "2018-04-02 10:20:30",
		"end_time": "2018-04-02 20:20:30",
	}
 */
func GetFpByTime(params map[string]string) (map[string]int64, error) {
	var ok bool
	if _,ok = params["start_time"]; !ok { 
		params["start_time"] = params["date"] + " 00:00:00"
	}
	if _,ok = params["end_time"]; !ok { 
		params["end_time"] = params["date"] + " 23:59:59"
	}

	body, _, err := PostES("getfp_time", params)
	if err != nil {
		return nil, err
	}
	var fps = map[string]int64{}
	jparse.ArrayEach(body, func(value []byte, dataType jparse.ValueType, offset int, err error) {
		if err == nil {
			fp, _ := jparse.GetString(value, "key")
			total, _ := jparse.GetInt(value, "doc_count")
			fps[fp] = total
		}
	}, "aggregations", "fps", "buckets")
	return fps, nil
}

/**
	var params = map[string]string{
		"date": "2018-04-02",
		"remote_addr": "182.133.219.108",
		"start_time": "2018-04-02 10:20:30",
		"end_time": "2018-04-02 20:20:30",
	}
 */
func GetFpByIp(params map[string]string) (map[string]int64, error) {
	var ok bool
	if _,ok = params["start_time"]; !ok {
		params["start_time"] = params["date"] + " 00:00:00"
	}
	if _,ok = params["end_time"]; !ok {
		params["end_time"] = params["date"] + " 23:59:59"
	}

	body, _, err := PostES("getfp_ip", params)
	if err != nil {
		return nil, err
	}
	var fps = map[string]int64{}
	jparse.ArrayEach(body, func(value []byte, dataType jparse.ValueType, offset int, err error) {
		if err == nil {
			fp, _ := jparse.GetString(value, "key")
			total, _ := jparse.GetInt(value, "doc_count")
			fps[fp] = total
		}
	}, "aggregations", "fps", "buckets")
	return fps, nil
}

//记录结构体
type Record struct {
	Continent				string		`json:"continent"`
	Country					string		`json:"country"`
	City					string		`json:"city"`
	Upstream_addr			string		`json:"upstream_addr"`
	Isp						string		`json:"isp"`
	Timestamp				string		`json:"Timestamp"`
	Platform				string		`json:"platform"`
	Session_time			float64		`json:"session_time"`
	Manufacturer			string		`json:"manufacturer"`
	Tcp_connect_time		string		`json:"tcp_connect_time"`
	Accesskey				string		`json:"accesskey"`
	Province				string		`json:"province"`
	Client_connect_time		float64		`json:"client_connect_time"`
	Fingerprint				string		`json:"fingerprint"`
	Brand					string		`json:"brand"`
	Remote_addr				string		`json:"remote_addr"`
	Simulator				int64		`json:"simulator"`
	Version					int64		`json:"version"`
	Server_addr				string		`json:"Server_addr"`
	Upstream_session_time	string		`json:"upstream_session_time"`
	Upstream_session_times	[]float64	`json:"upstream_session_times,omitempty"`
	Upstream_connect_time	string		`json:"upstream_connect_time"`
	Upstream_connect_times	[]float64	`json:"upstream_connect_times,omitempty"`
	Geo_val					[]float64	`json:"geo_val"`
	Device					string		`json:"device"`
}

func StringToFloat64Array(raw string) []float64 {
	tmp := strings.Split(raw, ",")
	var response = []float64{}
	for _,v := range tmp {
		ret , _ := strconv.ParseFloat(v, 64)
		response = append(response, ret)
	}
	return response
}

/**
	var params = map[string]string{
		"date": "2018-04-02",
		"fingerprint": "7B54F84697BFECAF90AE230B500625E7",
		"start_time": "2018-04-02 10:20:30",
		"end_time": "2018-04-02 20:20:30",
		"from": "0",
		"size": "1000",
	}
 */
func Count(apiKey string, params map[string]string) (int64, error) {
	var err error
	body, _, err := PostES("count-"+apiKey, params)
	if err != nil {
		return 0, err
	}
	total , err := jparse.GetInt(body, "count")
	if err != nil {
		return 0, err
	}
	return total,nil
}

/**
	var params = map[string]string{
		"date": "2018-04-02",
		"fingerprint": "7B54F84697BFECAF90AE230B500625E7",
		"start_time": "2018-04-02 10:20:30",
		"end_time": "2018-04-02 20:20:30",
		"from": "0",
		"size": "1000",
	}
 */
func GetHistory(apiKey string, params map[string]string) ([]Record, error) {
	body, _, err := PostES(apiKey, params)
	if err != nil {
		return nil, err
	}
	var records []Record
	jparse.ArrayEach(body, func(value []byte, dataType jparse.ValueType, offset int, err error) {
		if err == nil {
			rd,_,_,err := jparse.Get(value, "_source")
			var record Record
			err = json.Unmarshal(rd, &record)
			if err == nil {
				record.Upstream_session_times = StringToFloat64Array(record.Upstream_session_time)
				record.Upstream_connect_times = StringToFloat64Array(record.Upstream_connect_time)
				records = append(records, record)
			} else {
				//return offset, err
			}
		}
	}, "hits", "hits")
	return records, nil
}

func ParseFp(params map[string]string) (map[string]interface{}, error) {
	var ok bool
	if _,ok = params["start_time"]; !ok {
		params["start_time"] = params["date"] + " 00:00:00"
	}
	if _,ok = params["end_time"]; !ok {
		params["end_time"] = params["date"] + " 23:59:59"
	}
	var apiKey string = "gethistory_fp"
	total, err := Count(apiKey, params)
	if err != nil {
		return nil, err
	}

	var pagesize int64 = 10000
	var pageTotal int64 = (total / pagesize) + 1

	var result map[string]interface{}
	var i int64
	for i = 0; i < pageTotal; i++ {
		params["from"] = strconv.FormatInt(i * pagesize, 10)
		params["size"] = strconv.FormatInt(pagesize, 10)
		rows, _ := GetHistory(apiKey, params)
		result = Parse(rows, "ip", result)
	}
	return result, nil
}

func ParseIp(params map[string]string) (map[string]interface{}, error) {
	var ok bool
	if _,ok = params["start_time"]; !ok {
		params["start_time"] = params["date"] + " 00:00:00"
	}
	if _,ok = params["end_time"]; !ok {
		params["end_time"] = params["date"] + " 23:59:59"
	}
	var apiKey string = "gethistory_ip"
	total, err := Count(apiKey, params)
	if err != nil {
		return nil, err
	}

	var pagesize int64 = 10000
	var pageTotal int64 = (total / pagesize) + 1

	var result map[string]interface{}
	var i int64
	for i = 0; i < pageTotal; i++ {
		params["from"] = strconv.FormatInt(i * pagesize, 10)
		params["size"] = strconv.FormatInt(pagesize, 10)
		rows, _ := GetHistory(apiKey, params)
		result = Parse(rows, "ip", result)
	}
	return result, nil
}

func Parse(records []Record, ptype string, result map[string]interface{}) map[string]interface{} {
	var ok bool
	var response = map[string]interface{}{}

	var fps = []string{}
	var fpMap = map[string]int64{}
	if _,ok = result["fingerprint"]; ok {
		fps = result["fingerprint"].([]string)
		fpMap = result["fingerprintMap"].(map[string]int64)
	}

	var brands = []string{}
	var brandMap = map[string]int64{}
	if _,ok = result["brand"]; ok {
		brands = result["brand"].([]string)
		brandMap = result["brandMap"].(map[string]int64)
	}

	var clientIps = []string{}
	var clientIpMap = map[string]int64{}
	if _,ok = result["clientIp"]; ok {
		clientIps = result["clientIp"].([]string)
		clientIpMap = result["clientIpMap"].(map[string]int64)
	}

	var isps = []string{}
	var ispMap = map[string]int64{}
	if _,ok = result["isp"]; ok {
		isps = result["isp"].([]string)
		ispMap = result["ispMap"].(map[string]int64)
	}

	//var reqRate int64 = 0
	var reqCount int64 = 0
	if _,ok = result["reqCount"]; ok {
		reqCount = result["reqCount"].(int64)
	}

	//地区
	var area string
	var areas  = []string{}
	var areaMap = map[string]int64{}
	if _,ok = result["area"]; ok {
		areas  = result["area"].([]string)
		areaMap = result["areaMap"].(map[string]int64)
	}

	var games = []string{}
	var gameMap = map[string]int64{}
	if _,ok = result["game"]; ok {
		games = result["game"].([]string)
		gameMap = result["gameMap"].(map[string]int64)
	}
	//var sessionTimes = []float64{}
	//var sessionTimeMap = map[float64]int64{}
	//if _,ok = result["sessionTime"]; ok {
	//	sessionTimes = result["sessionTime"].([]float64)
	//	sessionTimeMap = result["sessionTimeMap"].(map[float64]int64)
	//}
	var platforms = []string{}
	var platformMap = map[string]int64{}
	if _,ok = result["platform"]; ok {
		platforms = result["platform"].([]string)
		platformMap = result["platformMap"].(map[string]int64)
	}
	//var proxyIp = []string{}
	//var proxyIpMap = map[string]int64{}
	//if _,ok = result["proxyIp"]; ok {
	//	proxyIp = result["proxyIp"].([]string)
	//	proxyIpMap = result["proxyIpMap"].(map[string]int64)
	//}
	var versions = []int64{}
	var versionMap = map[int64]int64{}
	if _,ok = result["version"]; ok {
		versions = result["version"].([]int64)
		versionMap = result["versionMap"].(map[int64]int64)
	}
	var devices = []string{}
	var deviceMap = map[string]int64{}
	if _,ok = result["device"]; ok {
		devices = result["device"].([]string)
		deviceMap = result["deviceMap"].(map[string]int64)
	}

	for _,v := range records {
		reqCount += 1

		//统计指纹
		if ptype == "ip" {
			if _,ok = fpMap[v.Fingerprint]; !ok {
				fpMap[v.Fingerprint] = 1
				fps = append(fps, v.Fingerprint)
			}
		}

		//统计手机品牌
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = brandMap[v.Brand]; !ok {
				brandMap[v.Brand] = 1
				brands = append(brands, v.Brand)
			}
		}
		//统计客户端IP
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = clientIpMap[v.Remote_addr]; !ok {
				clientIpMap[v.Remote_addr] = 1
				clientIps = append(clientIps, v.Remote_addr)
			}
		}

		//统计ISP
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = ispMap[v.Isp]; !ok {
				ispMap[v.Isp] = 1
				isps = append(isps, v.Isp)
			}
		}
		//统计地域
		area = v.Continent+","+v.Country+","+v.Province+","+v.City
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = areaMap[area]; !ok {
				areaMap[area] = 1
				areas = append(areas, area)
			}
		}
		//统计游戏
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = gameMap[v.Accesskey]; !ok {
				gameMap[v.Accesskey] = 1
				games = append(games, v.Accesskey)
			}
		}
		//会话时长
		//if _,ok = sessionTimeMap[v.Session_time]; !ok {
		//	sessionTimeMap[v.Session_time] = 1
		//	sessionTimes = append(sessionTimes, v.Session_time)
		//}
		//平台
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = platformMap[v.Platform]; !ok {
				platformMap[v.Platform] = 1
				platforms = append(platforms, v.Platform)
			}
		}
		//代理IP
		//SDK版本
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = versionMap[v.Version]; !ok {
				versionMap[v.Version] = 1
				versions = append(versions, v.Version)
			}
		}
		//设备
		if ptype == "fingerprint" || ptype == "ip" {
			if _,ok = deviceMap[v.Device]; !ok {
				deviceMap[v.Device] = 1
				devices = append(devices, v.Device)
			}
		}
	}
	if ptype == "ip" {
		response["fingerprint"] = fps
		response["fingerprintMap"] = fpMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["brand"] = brands
		response["brandMap"] = brandMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["clientIp"] = clientIps
		response["clientIpMap"] = clientIpMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["isp"] = isps
		response["ispMap"] = ispMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["area"] = areas
		response["areaMap"] = areaMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["game"] = games
		response["gameMap"] = gameMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["platform"] = platforms
		response["platformMap"] = platformMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["version"] = versions
		response["versionMap"] = versionMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["device"] = devices
		response["deviceMap"] = deviceMap
	}
	if ptype == "fingerprint" || ptype == "ip" {
		response["reqCount"] = reqCount
	}
	return response
}
