package beanstalk

import (
	//"fmt"
	"errors"
	"net/textproto"
)

// beanstalk 协议，名字中允许使用的字符
const NameChars = `\-+/;.$_()0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

var (
	configs    = map[string]map[string]string{}
	space      = []byte{' '}
	crnl       = []byte{'\r', '\n'}
	yamlHead   = []byte{'-', '-', '-', '\n'}
	nl         = []byte{'\n'}
	colonSpace = []byte{':', ' '}
	minusSpace = []byte{'-', ' '}
	//名称格式错误。NameError的err字段包含其中一个
	ErrEmpty      = errors.New("name is empty")
	ErrBadChar    = errors.New("name has bad char") // 有非法字符
	ErrTooLong    = errors.New("name is too long")
	ErrBadFormat  = errors.New("bad command format")
	ErrBuried     = errors.New("buried")
	ErrDeadline   = errors.New("deadline soon")
	ErrDraining   = errors.New("draining")
	ErrInternal   = errors.New("internal error")
	ErrJobTooBig  = errors.New("job too big")
	ErrNoCRLF     = errors.New("expected CR LF")
	ErrNotFound   = errors.New("not found")
	ErrNotIgnored = errors.New("not ignored")
	ErrOOM        = errors.New("server is out of memory")
	ErrTimeout    = errors.New("timeout")
	ErrUnknown    = errors.New("unknown command")
)

var respError = map[string]error{
	"BAD_FORMAT":      ErrBadFormat,
	"BURIED":          ErrBuried,
	"DEADLINE_SOON":   ErrDeadline,
	"DRAINING":        ErrDraining,
	"EXPECTED_CRLF":   ErrNoCRLF,
	"INTERNAL_ERROR":  ErrInternal,
	"JOB_TOO_BIG":     ErrJobTooBig,
	"NOT_FOUND":       ErrNotFound,
	"NOT_IGNORED":     ErrNotIgnored,
	"OUT_OF_MEMORY":   ErrOOM,
	"TIMED_OUT":       ErrTimeout,
	"UNKNOWN_COMMAND": ErrUnknown,
}

type Beanstalk struct {
	Conn       *textproto.Conn
	Addr       string
	Tube       string
	WatchTubes map[string]bool //watched
}

type Request struct {
	id uint
	op string
}

// 名称错误表示名称格式不正确，并且特定错误
type NameError struct {
	Name string
	Err  error
}

type BSError struct {
	Conn *Beanstalk
	Op   string
	Err  error
}

func (e BSError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

type unknownRespError string

func (e unknownRespError) Error() string {
	return "unknown response: " + string(e)
}

func findRespError(s string) error {
	if err := respError[s]; err != nil {
		return err
	}
	return unknownRespError(s)
}
