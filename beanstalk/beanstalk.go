package beanstalk

import (
    "fmt"
    "time"
    Bclient "github.com/beanstalkd/go-beanstalk"
    "github.com/spf13/viper"
)

//var conn *Bclient.Conn

func Put(tubeName string, body []byte) (uint64, error) {
    conn, err := Conn()
    if err != nil {
        fmt.Println(err)
    }
    tube := &Bclient.Tube{conn, tubeName}
    id, err := tube.Put(body, 1, 0, 30*time.Second)
    return id,err
}

func Conn() (*Bclient.Conn, error) {
    return Bclient.Dial("tcp", viper.GetString("beanstalk.host") + ":" + viper.GetString("beanstalk.port"))
}

func init() {
    //conn, err := Bclient.Dial("tcp", "test.yundun.com:11300")
    //if err != nil {
    //    fmt.Println(err)
    //}
}

