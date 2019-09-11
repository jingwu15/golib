package apidoc

import (
    "fmt"
    "strings"
)

type ApiDoc struct {
    Type        string
    Url         string
    Title       string
    Version     string
    Group       string
    Name        string
    Header      map[string]map[string][]map[string]interface{}
    Parameter   map[string]interface{}
    Success     map[string]interface{}
    Error       map[string][]map[string]string
    Filename    string
    GroupTitle  string
}

// api doc 中提取出来的路由信息
type ApiRoute struct {
    Methods []string
    Hosts   []string
    Name    string
    Paths   []string
    //    "methods": row["type"].upper().replace('''/''', "|").split("|"),
    //    "hosts": ["yundunapiv4.test.nodevops.cn", "kong.yundunapiv4.test.nodevops.cn"],
    //    "name": row["url"].replace(".", "-").replace("{", "").replace("}", ""),
    //    "paths": ["/V4/" + row["url"]],
    //item["name"] = item["name"] + "-" + item["methods"][0]
}

//转换api为路由信息
func ApiToRoute(apiDocs []ApiDoc) []ApiRoute {
    m := map[string]int{
        "DELETE": 1,
        "GET": 1,
        "POST": 1,
        "PUT": 1,
    }
    apiRoutes := []ApiRoute{}
    for _, apidoc := range apiDocs {
        typeRaw := strings.ToUpper(apidoc.Type)
        typeRaw = strings.Replace(typeRaw, "/", "|", -1)
        methods := strings.Split(typeRaw, "|")
        for _, method := range methods {
            if _, ok := m[method]; !ok {
                fmt.Println(apidoc, "-------err----------", method, " not support")
                continue
            }
        }
        nameRaw := strings.Replace(apidoc.Name, ".", "-", -1)
        nameRaw = strings.Replace(nameRaw, "{", "", -1)
        nameRaw = strings.Replace(apidoc.Name, "}", "", -1)
        apiRoute := ApiRoute{
            Methods: methods,
            Hosts: []string{"yundunapiv4.test.nodevops.cn", "kong.yundunapiv4.test.nodevops.cn"},
            Name: nameRaw,
            Paths: []string{fmt.Sprintf("/V4/%s", apidoc.Url)},
        }
        //fmt.Println(methods, nameRaw, apidoc.Url)
        apiRoutes = append(apiRoutes, apiRoute)
    }
    return apiRoutes
}
