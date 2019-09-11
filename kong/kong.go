package kong

import (
    "fmt"
//    "time"
//    "strings"
//    "net/http"
//    "io/ioutil"
    "github.com/jingwu15/golib/json"
    hc "github.com/jingwu15/golib/http/client"
)

var (
    API_ROUTE_ALL       = "%s/routes"
    API_ROUTE_LIST      = "%s/services/%s/routes"
    API_ROUTE_ADD       = "%s/services/%s/routes"
    API_ROUTE_GET       = "%s/routes/%s"
    API_ROUTE_UPDATE    = "%s/routes/%s"
    API_ROUTE_DEL       = "%s/route/%s"

    API_SERVICE_LIST    = "%s/services"
    API_SERVICE_ADD     = "%s/services"
    API_SERVICE_GET     = "%s/services/%s"
    API_SERVICE_UPDATE  = "%s/services/%s"
    API_SERVICE_DEL     = "%s/services/%s"

    API_UPSTREAM_LIST   = "%s/upstreams"
    API_UPSTREAM_ADD    = "%s/upstreams"
    API_UPSTREAM_GET    = "%s/upstreams/%s"
    API_UPSTREAM_UPDATE = "%s/upstreams/%s"
    API_UPSTREAM_DEL    = "%s/upstreams/%s"

    API_TARGET_LIST     = "%s/upstreams/%s/targets"
    API_TARGET_ADD      = "%s/upstreams/%s/targets"
    API_TARGET_DEL      = "%s/upstreams/%s/targets/%s"
)

type Kong struct {
    Protocol    string
    Host        string
    Port        string
    Api         string
    Headers     map[string]string
}

type Service struct {
    ID              string      `json:"id,omitempty"`
    Name            string      `json:"name"`
    Retries         int         `json:"retries"`
    Protocol        string      `json:"protocol"`
    Host            string      `json:"host"`
    Port            int         `json:"port"`
    Path            string      `json:"path,omitempty"`
    ConnectTimeout  int         `json:"connect_timeout,omitempty"`
    WriteTimeout    int         `json:"write_timeout,omitempty"`
    ReadTimeout     int         `json:"read_timeout,omitempty"`
    Tags            []string    `json:"tags,omitempty"`
}

//type Upstream struct {
////id
////tags
////name
////slots
////hash_on
////created_at
////
////hash_fallback
////hash_on_header
////hash_fallback_header
////hash_on_cookie
////hash_on_cookie_path
////
////healthchecks.active.https_verify_certificate
////healthchecks.active.unhealthy.http_statuses
////healthchecks.active.unhealthy.tcp_failures
////healthchecks.active.unhealthy.timeouts
////healthchecks.active.unhealthy.http_failures
////healthchecks.active.unhealthy.interval
////healthchecks.active.http_path
////healthchecks.active.timeout
////healthchecks.active.healthy.http_statuses
////healthchecks.active.healthy.interval
////healthchecks.active.healthy.successes
////healthchecks.active.https_sni
////healthchecks.active.concurrency
////healthchecks.active.type
////
////healthchecks.passive.unhealthy.http_failures
////healthchecks.passive.unhealthy.http_statuses
////healthchecks.passive.unhealthy.tcp_failures
////healthchecks.passive.unhealthy.timeouts
////healthchecks.passive.type
////healthchecks.passive.healthy.successes
////healthchecks.passive.healthy.http_statuses
//}

type Target struct {
    ID          string      `json:"id,omitempty"`
    Upstream    map[string]string      `json:"upstream,omitempty"`
    Target	    string      `json:"target"`
    Weight      int         `json:"weight"`
    Tags        []string    `json:"tags,omitempty"`
    CreatedAt   float64         `json:"create_at,omitempty"`
}

type ServicePage struct {
    Next    string
    Offset  string
    Data    []Route
    Kong    *Kong
}

type Route struct {
    ID              string      `json:"id,omitempty"`
    Name            string      `json:"name"`
    Protocols       []string    `json:"protocols"`
    Methods         []string    `json:"methods"`
    Hosts           []string    `json:"hosts"`
    Paths           []string    `json:"paths"`
    Regex_priority  int         `json:"regex_priority"`
    Strip_path      bool        `json:"strip_path"`
    Preserve_host   bool        `json:"preserve_host,omitempty"`
    Snis            []string    `json:"snis"`
    Sources         []string    `json:"sources"`
    Destinations    []string    `json:"destinations"`
    Tags            []string    `json:"tags"`
}

type RoutePage struct {
    Next    string
    Offset  string
    Data    []Route
    Kong    *Kong
}

func NewKong(api string) Kong {
    return Kong{
        Protocol: "http",
        Host: "172.16.100.188",
        Port: "8001",
        Api: api,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}

func NewRoute() Route {
    return Route {
        Name:           "",
        Protocols:      []string{"http"},
        Methods:        []string{},
        Hosts:          []string{},
        Paths:          []string{},
        Regex_priority: 5,
        Strip_path:     false,
        Preserve_host:  false,
        Snis:           []string{},
        Sources:        []string{},
        Destinations:   []string{},
        Tags:           []string{},
    }
}

func NewService() Service {
    return Service {
        ID:             "",
        Name:           "",
        Retries:        5,
        Protocol:       "http",
        Host:           "",
        Port:           80,
        Path:           "",
        ConnectTimeout: 60000,
        WriteTimeout:   60000,
        ReadTimeout:    60000,
        Tags:           []string{},
    }
}

func (p RoutePage)ReadNext() (RoutePage, error) {
    k := *p.Kong
    page := RoutePage{Kong: &k}
    httpcode, body, _, err := hc.Get(k.Api + p.Next, k.Headers, 10)
    if err != nil { return page, err }
    if httpcode != 200 { return page, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &page)
    return page, err
}

func (p ServicePage)ReadNext() (ServicePage, error) {
    k := *p.Kong
    page := ServicePage{Kong: &k}
    httpcode, body, _, err := hc.Get(k.Api + p.Next, k.Headers, 10)
    if err != nil { return page, err }
    if httpcode != 200 { return page, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &page)
    return page, err
}

func (k Kong)RouteAll() (RoutePage, error) {
    page := RoutePage{Kong: &k}
    url := fmt.Sprintf(API_ROUTE_ALL, k.Api)
    httpcode, body, _, err := hc.Get(url, k.Headers, 10)
    if err != nil { return page, err }
    if httpcode != 200 { return page, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &page)
    return page, err
}

func (k Kong)RouteList(service string) (RoutePage, error) {
    page := RoutePage{Kong: &k}
    url := fmt.Sprintf(API_ROUTE_LIST, k.Api, service)
    httpcode, body, _, err := hc.Get(url, k.Headers, 10)
    if err != nil { return page, err }
    if httpcode != 200 { return page, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &page)
    return page, err
}

func (k Kong)RouteSave(service string, route Route) (Route, error) {
    obj := Route{}
    reqBody, _  := json.Encode(route)
    _, e := k.RouteGet(route.Name)
    if e != nil {        //不存在
        url := fmt.Sprintf(API_ROUTE_ADD, k.Api, service)
        httpcode, body, _, err := hc.Post(url, string(reqBody), k.Headers, 10)

        if err != nil { return obj, err }
        if httpcode != 201 { return obj, fmt.Errorf(body) }

        err = json.Decode([]byte(body), &obj)
        return obj, err
    } else {            //存在
        url := fmt.Sprintf(API_ROUTE_UPDATE, k.Api, route.Name)
        httpcode, body, _, err := hc.Patch(url, string(reqBody), k.Headers, 10)

        if err != nil { return obj, err }
        if httpcode != 200 { return obj, fmt.Errorf(body) }

        err = json.Decode([]byte(body), &obj)
        return obj, err
    }
    return obj, nil
}

func (k Kong)RouteGet(name string) (Route, error) {
    route := Route{}
    url := fmt.Sprintf(API_ROUTE_GET, k.Api, name)
    httpcode, body, _, err := hc.Get(url, k.Headers, 10)
    if err != nil { return route, err }
    if httpcode != 200 { return route, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &route)
    return route, err
}

func (k Kong)RouteDelete(name string) (bool, error) {
    url := fmt.Sprintf(API_ROUTE_DEL, k.Api, name)
    httpcode, body, _, err := hc.Delete(url, k.Headers, 10)
    if err != nil { return false, err }
    if httpcode == 204 { return false, fmt.Errorf(body) }
    return true, nil
}

func (k Kong)ServiceList() (ServicePage, error) {
    page := ServicePage{Kong: &k}
    url := fmt.Sprintf(API_SERVICE_LIST, k.Api)
    httpcode, body, _, err := hc.Get(url, k.Headers, 10)
    if err != nil { return page, err }
    if httpcode != 200 { return page, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &page)
    return page, err
}

func (k Kong)ServiceGet(name string) (Service, error) {
    service := Service{}
    url := fmt.Sprintf(API_SERVICE_GET, k.Api, name)
    httpcode, body, _, err := hc.Get(url, k.Headers, 10)
    if err != nil { return service, err }
    if httpcode != 200 { return service, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &service)
    return service, err
}

func (k Kong)ServiceSave(name string, service Service) (Service, error) {
    obj := Service{}
    reqBody, _  := json.Encode(service)
    _, e := k.ServiceGet(name)
    if e != nil {        //不存在
        url := fmt.Sprintf(API_SERVICE_ADD, k.Api)
        httpcode, body, _, err := hc.Post(url, string(reqBody), k.Headers, 10)

        if err != nil { return obj, err }
        if httpcode != 201 { return obj, fmt.Errorf(body) }

        err = json.Decode([]byte(body), &obj)
        return obj, err
    } else {            //存在
        url := fmt.Sprintf(API_SERVICE_UPDATE, k.Api, name)
        httpcode, body, _, err := hc.Patch(url, string(reqBody), k.Headers, 10)

        if err != nil { return obj, err }
        if httpcode != 200 { return obj, fmt.Errorf(body) }

        err = json.Decode([]byte(body), &obj)
        return obj, err
    }
    return obj, nil
}

func (k Kong)ServiceDelete(name string) (bool, error) {
    url := fmt.Sprintf(API_ROUTE_DEL, k.Api, name)
    httpcode, body, _, err := hc.Delete(url, k.Headers, 10)
    if err != nil { return false, err }
    if httpcode == 204 { return false, fmt.Errorf(body)  }
    return true, nil
}

func (k Kong)UpstreamGet(name string) (d map[string]interface{}, e error) {
    d = map[string]interface{}{}
    //name := fmt.Sprintf("%s:%s", host, port)
    url := fmt.Sprintf(API_UPSTREAM_GET, k.Api, name)
    httpcode, body, _, err := hc.Get(url, k.Headers, 10)
    if err != nil { return d, e }
    if httpcode != 200 { return d, fmt.Errorf(body) }
    err = json.Decode([]byte(body), &d)
    return d, err
}

func (k Kong)UpstreamSave(name string) (e error) {
    obj := map[string]interface{}{}
    reqBody := fmt.Sprintf(`{"name": "%s"}`, name)
    _, e = k.UpstreamGet(name)
    if e != nil {        //不存在
        url := fmt.Sprintf(API_UPSTREAM_ADD, k.Api)
        httpcode, body, _, e := hc.Post(url, reqBody, k.Headers, 10)

        if e != nil { return e }
        if httpcode != 201 { return fmt.Errorf(body) }

        e = json.Decode([]byte(body), &obj)
        return e
    } else {            //存在
        url := fmt.Sprintf(API_UPSTREAM_UPDATE, k.Api, name)
        httpcode, body, _, e := hc.Patch(url, reqBody, k.Headers, 10)

        if e != nil { return e }
        if httpcode != 200 { return fmt.Errorf(body) }

        e = json.Decode([]byte(body), &obj)
        return e
    }
    return nil
}

func (k Kong)TargetSave(upName string, target Target) (d Target, e error) {
    reqBody, _ := json.Encode(target)

    url := fmt.Sprintf(API_TARGET_ADD, k.Api, upName)
    httpcode, body, _, e := hc.Post(url, string(reqBody), k.Headers, 10)

    if e != nil { return d, e }
    if httpcode != 201 { return d, fmt.Errorf(body) }

    e = json.Decode([]byte(body), &d)
    return d, e
}
