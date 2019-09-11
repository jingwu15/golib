package misc

import (
	"os"
    "net"
	"sync"
	"time"
	"reflect"
    "strings"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

func GetRandNum(num int) int {
	rand.Seed(time.Now().UnixNano())
	var mu sync.Mutex
	mu.Lock()
	result := rand.Intn(num)
	mu.Unlock()
	return result
}

func KeysMapString(dataMap map[string]string) []string {
	keys := make([]string, len(dataMap))

	i := 0
	for key := range dataMap {
		keys[i] = key
		i++
	}

	return keys
}

func KeysMapInt64(dataMap map[string]int64) []string {
	keys := make([]string, len(dataMap))

	i := 0
	for key := range dataMap {
		keys[i] = key
		i++
	}

	return keys
}

func IntsToMap(ints []int) map[int]int {
	var mapInt = map[int]int{}
	for _, i := range ints {
		mapInt[i] = 1
	}
	return mapInt
}

func MapKeyExist(mapInt map[int]int, key int) bool {
	_, ok := mapInt[key]
	return ok
}

func SelectIp(ips [][]byte, limit int) []string {
	ipsSelect := map[string]string{}
	lenIps := len(ips)
	if lenIps < limit {
		for _, ip := range ips {
			ipsSelect[string(ip)] = "1"
		}
	} else {
		i := 0
		for {
			//限制取到的数据量
			if len(ipsSelect) >= limit {
				break
			}
			//避免死循环
			if i > 100 {
				break
			}
			index := GetRandNum(lenIps)
			ip := ips[index]
			ipsSelect[string(ip)] = "1"
			i++
		}
	}
	return KeysMapString(ipsSelect)
}

func Timeout(tag, detailed string, start time.Time, timeLimit float64) {
	//应用示例
	//defer Timeout("SaveAppLogMain", "Total", time.Now(), float64(3))
	dis := time.Now().Sub(start).Seconds()
	if dis > timeLimit {
		log.Warning("timeout", tag, " detailed:", detailed, "TimeoutWarning using", dis, "s")
	}
}

func FormatPathSuffix(pathname string) string {
	length := len(pathname)
	suffix := pathname[length-1 : length]
	if suffix != "/" {
		pathname = pathname + "/"
	}
	return pathname
}

func Byte2string(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}

	return tmp[16:]
}

func Md5_sum(raw []byte) string {
	md5Sum := md5.Sum(raw)
	return hex.EncodeToString(Byte2string(md5Sum))
}

func File_exists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
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
            IP: net.ParseIP(rip),
            Port: 22,
        }
        udp, _ := net.DialUDP("udp", &laddr, &raddr)
        laInfo := strings.Split(udp.LocalAddr().String(), ":")
        machineIp["net"]    = "udp"
        machineIp["ip"]     = laInfo[0]
        machineIp["port"]   = laInfo[1]
    }
    return machineIp["net"], machineIp["ip"], machineIp["port"]
}

/*
 * 动态设置结构体中字段的值
 * 示例：
 * var record = Record{City: []byte(`ss`)}
 * _ = SetValue(reflect.ValueOf(&record), "City", []byte(`hello world`))
 * fmt.Println(record)
 */
//func SetValue(rVal reflect.Value, field string, val interface{}) error {
//	//判断是否可以设置值
//	if rVal.Kind() == reflect.Ptr && !rVal.Elem().CanSet() {
//		return errors.New("connot set")
//	}
//	fVal := rVal.Elem().FieldByName(field)
//	if !fVal.IsValid() {
//		return errors.New("value is zero")
//	}
//	switch fVal.Kind() {
//	case reflect.Bool:
//		fVal.SetBool(val.(bool))
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//		fVal.SetInt(val.(int64))
//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//		fVal.SetUint(val.(uint64))
//	case reflect.Float32, reflect.Float64:
//		fVal.SetFloat(val.(float64))
//	case reflect.String:
//		fVal.SetString(val.(string))
//	case reflect.UnsafePointer, reflect.Uintptr, reflect.Ptr:
//		fVal.SetPointer(val.(unsafe.Pointer))
//	case reflect.Complex64, reflect.Complex128:
//		fVal.SetComplex(val.(complex128))
//	//case reflect.Map:		//@todo
//		//switch val.(type) {
//		//case map[string]string:
//		//	fVal.Set(val.(map[string]string))
//		//default:
//		//}
//	case reflect.Slice:
//		switch val.(type) {
//		//switch vtype:=val.(type) {
//		case []byte:
//			fVal.SetBytes(val.([]byte))
//		default:
//			//return errors.New("not supported the type: " + vtype)
//			return errors.New("not supported the type: ")
//		}
//	//case reflect.Struct:	//@todo
//	default:
//		return errors.New("not supported the type: " + fVal.Kind().String())
//    //Array
//    //Chan
//    //Func
//    //Interface
//	}
//	return nil
//}
