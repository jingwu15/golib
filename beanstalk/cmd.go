package beanstalk

import (
	"strconv"
	"time"
)

func (bs *Beanstalk) Use(tube string) error {
	request, err := bs.cmd("use", tube)
	if err != nil {
		return err
	}
	_, err = bs.readResponse(request, true, "USING")
	return err
}

// Delete 删除任务
func (bs *Beanstalk) Delete(id uint64) error {
	request, err := bs.cmd("delete", id)
	if err != nil {
		return err
	}
	_, err = bs.readResponse(request, false, "DELETED")
	return err
}

//Release 命令, 执行以下操作：
//将给定作业的优先级设置为pri，将其从由C保留的作业，等待延迟秒，然后将该作业放入就绪队列，使其可供任何客户端预订。
func (bs *Beanstalk) Release(id uint64, pri uint32, delay int) error {
	request, err := bs.cmd("release", id, pri, delay)
	if err != nil {
		return err
	}
	_, err = bs.readResponse(request, false, "RELEASED")
	return err
}

// Bury 将任务状态设置为BURIED，并设置优先级。BURIED状态在未改变前，不会被消费
func (bs *Beanstalk) Bury(id uint64, pri uint32) error {
	request, err := bs.cmd("bury", id, pri)
	if err != nil {
		return err
	}
	_, err = bs.readResponse(request, false, "BURIED")
	return err
}

// Touch 允许worker请求更多的时间执行job，这个对长时间执行的任务有用
func (bs *Beanstalk) Touch(id uint64) error {
	request, err := bs.cmd("touch", id)
	if err != nil {
		return err
	}
	_, err = bs.readResponse(request, false, "TOUCHED")
	return err
}

// Peek 从服务器获取指定作业的副本。
// peek <id>\r\n  返回id对应的job
// peek-ready\r\n 返回下一个ready job
// peek-delayed\r\n 返回下一个延迟剩余时间最短的job
// peek-buried\r\n 返回下一个在buried列表中的job
func (bs *Beanstalk) Peek(id uint64) (body []byte, err error) {
	request, err := bs.cmd("peek", id)
	if err != nil {
		return nil, err
	}
	return bs.readResponse(request, true, "FOUND %d", &id)
}

// Stats 返回整个消息队列系统的整体信息
func (bs *Beanstalk) Stats() (map[string]string, error) {
	request, err := bs.cmd("stats")
	if err != nil {
		return nil, err
	}
	body, err := bs.readResponse(request, true, "OK")
	return parseDict(body), err
}

// StatsJob 统计job的相关信息
func (bs *Beanstalk) StatsJob(id uint64) (map[string]string, error) {
	request, err := bs.cmd("stats-job", id)
	if err != nil {
		return nil, err
	}
	body, err := bs.readResponse(request, true, "OK")
	return parseDict(body), err
}

// ListTubes 列出所有存在的tube
func (bs *Beanstalk) ListTubes() ([]string, error) {
	request, err := bs.cmd("list-tubes")
	if err != nil {
		return nil, err
	}
	body, err := bs.readResponse(request, true, "OK")
	return parseList(body), err
}

// Put 插入一个job到队列
func (bs *Beanstalk) Put(body []byte, pri uint32, delay, ttr int) (id uint64, err error) {
	request, err := bs.cmdPut("put", body, pri, strconv.Itoa(delay), strconv.Itoa(ttr))
	if err != nil {
		return 0, err
	}
	_, err = bs.readResponse(request, false, "INSERTED %d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// PeekReady 返回下一个ready job
func (bs *Beanstalk) PeekReady() (id uint64, body []byte, err error) {
	request, err := bs.cmd("peek-ready")
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readResponse(request, true, "FOUND %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

// PeekDelayed 返回下一个延迟剩余时间最短的job
func (bs *Beanstalk) PeekDelayed() (id uint64, body []byte, err error) {
	request, err := bs.cmd("peek-delayed")
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readResponse(request, true, "FOUND %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

// PeekBuried 返回下一个在buried列表中的job
func (bs *Beanstalk) PeekBuried() (id uint64, body []byte, err error) {
	request, err := bs.cmd("peek-buried")
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readResponse(request, true, "FOUND %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

// Kick 此指令应用在当前使用的tube中，它将job的状态迁移为ready或者delayed
func (bs *Beanstalk) Kick(bound int) (n int, err error) {
	request, err := bs.cmd("kick", bound)
	if err != nil {
		return 0, err
	}
	_, err = bs.readResponse(request, false, "KICKED %d", &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Stats 统计tube的相关信息
func (bs *Beanstalk) StatsTube(tube string) (map[string]string, error) {
	request, err := bs.cmd("stats-tube", tube)
	if err != nil {
		return nil, err
	}
	body, err := bs.readResponse(request, true, "OK")
	return parseDict(body), err
}

// Pause 暂停任务一段时间
func (bs *Beanstalk) PauseTube(tube string, d time.Duration) error {
	request, err := bs.cmd("pause-tube", tube, time.Duration(d))
	if err != nil {
		return err
	}
	_, err = bs.readResponse(request, false, "PAUSED")
	if err != nil {
		return err
	}
	return nil
}

// Reserve 取出（预订）job，待处理。
// 它将返回一个新预订的job，如果没有job，beanstalkd将直到有job时才发送响应。
// 取出job时，状态迁移为reserved, client被限制在指定的ttr时间内完成，否则超时，状态迁移为ready。
func (bs *Beanstalk) Reserve(timeout time.Duration) (id uint64, body []byte, err error) {
	request, err := bs.cmd("reserve-with-timeout", time.Duration(timeout))
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readResponse(request, true, "RESERVED %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}
