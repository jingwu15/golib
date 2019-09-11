package cfg

import (
    "os"
    "fmt"
    "path/filepath"
	"github.com/spf13/viper"
)

type CfgViper struct {
    Ftype   string
    Cfiles  []string
    V       *viper.Viper
}

var cfgVipers = map[string]CfgViper{}

//设置配置文件，支持目录
func SetCfgViper(key, ftype string, cfiles []string, v *viper.Viper) {
    cfgVipers[key] = CfgViper{Ftype: ftype, Cfiles: cfiles, V: v}
}

func ReloadViper(key string, ignoreKeys... string) {
    cfg := cfgVipers[key]
    x := LoadViper(cfg.Ftype, cfg.Cfiles, ignoreKeys...)
    *cfg.V = *x
}

//根据给定的配置文件，加载Viper
func LoadViper(ftype string, cfiles[] string, ignoreKeys... string) *viper.Viper {
    fnames := []string{}
    for _, cfile := range cfiles {
        f, err := os.Stat(cfile)
        if err != nil && os.IsNotExist(err) {
            //文件不存在
            continue
        }
        if f.IsDir() {
            fns, err := filepath.Glob(cfile + "/*." + ftype)
            if err != nil {
                fmt.Println("配置文件查找失败")
                continue
            }
            for _, fn := range fns {
                fnames = append(fnames, fn)
            }
        } else {
            fnames = append(fnames, cfile)
        }
    }
    x := viper.New()
	//x.AddConfigPath(dir)
    x.SetConfigType(ftype)
    for _, fn := range fnames {
        f, _ := os.Open(fn)
	    _ = x.MergeConfig(f)
        f.Close()
    }
    return x
}

//读取文件夹
func ReadDir(ftype string, cfiles[] string) error {
    fnames := []string{}
    for _, cfile := range cfiles {
        f, err := os.Stat(cfile)
        if err != nil && os.IsNotExist(err) {
            //文件不存在
            continue
        }
        if f.IsDir() {
            fns, err := filepath.Glob(cfile + "/*." + ftype)
            if err != nil {
                fmt.Println("配置文件查找失败")
                continue
            }
            for _, fn := range fns {
                fnames = append(fnames, fn)
            }
        } else {
            fnames = append(fnames, cfile)
        }
    }

    viper.SetConfigType(ftype)
    for _, fn := range fnames {
        f, _ := os.Open(fn)
	    _ = viper.MergeConfig(f)
        f.Close()
    }
    return nil
}

