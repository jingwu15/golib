package beanstalk

import (
	"fmt"
	"io"
	"net"
	"net/textproto"
)

func New(addr string) (*Beanstalk, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return &Beanstalk{Addr: addr}, err
	}
	return &Beanstalk{Conn: textproto.NewConn(conn), Addr: addr}, nil
}

func (bs *Beanstalk) ReConn() error {
	err := bs.Close()
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", bs.Addr)
	if err != nil {
		return err
	}
	bs.Conn = textproto.NewConn(conn)
	return nil
}

// 关闭网络连接
func (bs *Beanstalk) Close() error {
	return bs.Conn.Close()
}

func (bs *Beanstalk) cmd(op string, args ...interface{}) (Request, error) {
	request := Request{bs.Conn.Next(), op}
	bs.Conn.StartRequest(request.id)

	fmt.Fprint(bs.Conn.W, op)
	for _, a := range args {
		bs.Conn.W.Write(space)
		fmt.Fprint(bs.Conn.W, a)
	}
	bs.Conn.W.Write(crnl)
	err := bs.Conn.W.Flush()
	if err != nil {
		return Request{}, BSError{bs, op, err}
	}
	bs.Conn.EndRequest(request.id)
	return request, nil
}

func (bs *Beanstalk) cmdPut(op string, body []byte, args ...interface{}) (Request, error) {
	request := Request{bs.Conn.Next(), op}
	bs.Conn.StartRequest(request.id)

	fmt.Fprint(bs.Conn.W, op)
	args = append(args, len(body))
	for _, a := range args {
		bs.Conn.W.Write(space)
		fmt.Fprint(bs.Conn.W, a)
	}
	bs.Conn.W.Write(crnl)
	bs.Conn.W.Write(body)
	bs.Conn.W.Write(crnl)
	err := bs.Conn.W.Flush()
	if err != nil {
		return Request{}, BSError{bs, op, err}
	}
	bs.Conn.EndRequest(request.id)
	return request, nil
}

func (bs *Beanstalk) readResponse(request Request, f string, a ...interface{}) error {
	bs.Conn.StartResponse(request.id)
	defer bs.Conn.EndResponse(request.id)
	line, err := bs.Conn.ReadLine()
	if err != nil {
		return BSError{bs, request.op, err}
	}
	toScan := line
	err = scan(toScan, f, a...)
	if err != nil {
		return BSError{bs, request.op, err}
	}
	return nil
}

func (bs *Beanstalk) readBody(request Request, f string, a ...interface{}) (body []byte, err error) {
	bs.Conn.StartResponse(request.id)
	defer bs.Conn.EndResponse(request.id)
	line, err := bs.Conn.ReadLine()
	if err != nil {
		return nil, BSError{bs, request.op, err}
	}
	toScan := line

	var size int
	toScan, size, err = parseSize(toScan)
	if err != nil {
		return nil, BSError{bs, request.op, err}
	}
	body = make([]byte, size+2) // 包括 CR NL
	_, err = io.ReadFull(bs.Conn.R, body)
	if err != nil {
		return nil, BSError{bs, request.op, err}
	}
	body = body[:size] // 不包括 CR NL

	err = scan(toScan, f, a...)
	if err != nil {
		return nil, BSError{bs, request.op, err}
	}
	return body, nil
}
