package misc

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

func GetRandNum(num int) int {
	rand.Seed(time.Now().UnixNano())
	var mu sync.Mutex
	mu.Lock()
	result := rand.Intn(num)
	mu.Unlock()
	return result
}

func Md5_sum(raw []byte) string {
	md5Sum := md5.Sum(raw)
	return hex.EncodeToString(Byte2string(md5Sum))
}

func Byte2string(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}

	return tmp[16:]
}

func IsEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

//取行操作系统的IP
var machineIp = map[string]string{"net": "", "ip": "", "port": ""}

func MachineIp(reqIps ...string) (protocol, ips, ports string) {
	if machineIp["ip"] == "" {
		lip, rip := "0.0.0.0", "119.29.29.29"
		if len(reqIps) == 1 {
			rip = reqIps[0]
		}
		laddr := net.UDPAddr{
			IP: net.ParseIP(lip),
		}
		raddr := net.UDPAddr{
			IP:   net.ParseIP(rip),
			Port: 22,
		}
		udp, _ := net.DialUDP("udp", &laddr, &raddr)
		laInfo := strings.Split(udp.LocalAddr().String(), ":")
		machineIp["net"] = "udp"
		machineIp["ip"] = laInfo[0]
		machineIp["port"] = laInfo[1]
	}
	return machineIp["net"], machineIp["ip"], machineIp["port"]
}
