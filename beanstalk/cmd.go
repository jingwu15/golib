package beanstalk

import (
	"strconv"
)

//cmd-put 总共执行put指令的次数
//cmd-reserve 总共执行reserve指令的次数
//cmd-use 总共执行use指令的次数
//cmd-watch 总共执行watch指令的次数
//cmd-ignore 总共执行ignore指令的次数
//cmd-release 总共执行release指令的次数
//cmd-bury 总共执行bury指令的次数
//cmd-kick 总共执行kick指令的次数
//cmd-stats 总共执行stats指令的次数
//cmd-stats-job 总共执行stats-job指令的次数
//cmd-stats-tube 总共执行stats-tube指令的次数
//cmd-list-tubes 总共执行list-tubes指令的次数
//cmd-list-tube-used 总共执行list-tube-used指令的次数
//cmd-list-butes-watched 总共执行list-tubes-watched指令的次数
//cmd-pause-tube 总共执行pause-tube指令的次数

func (bs *Beanstalk) Use(tube string) error {
	request, err := bs.cmd("use", tube)
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "USING")
	return err
}

func (bs *Beanstalk) Watch(tube string) (count int64, err error) {
	request, err := bs.cmd("watch", tube)
	if err != nil {
		return 0, err
	}
	err = bs.readResponse(request, "WATCHING %d", &count)
	return count, err
}

func (bs *Beanstalk) Ignore(tube string) (count int64, err error) {
	request, err := bs.cmd("ignore", tube)
	if err != nil {
		return 0, err
	}
	err = bs.readResponse(request, "WATCHING %d", &count)
	return count, err
}

// Put 插入一个job到队列
func (bs *Beanstalk) Put(body []byte, pri uint32, delay, ttr int) (id uint64, err error) {
	request, err := bs.cmdPut("put", body, pri, strconv.Itoa(delay), strconv.Itoa(ttr))
	if err != nil {
		return 0, err
	}
	err = bs.readResponse(request, "INSERTED %d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (bs *Beanstalk) UsePut(tube string, body []byte, pri uint32, delay, ttr int) (id uint64, err error) {
	err = bs.Use(tube)
	if err != nil {
		return 0, err
	}
	return bs.Put(body, pri, delay, ttr)
}

// Delete 删除任务
func (bs *Beanstalk) Delete(id uint64) error {
	request, err := bs.cmd("delete", id)
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "DELETED")
	return err
}

func (bs *Beanstalk) UseDelete(tube string, id uint64) error {
	err := bs.Use(tube)
	if err != nil {
		return err
	}
	return bs.Delete(id)
}

//Release 命令, 执行以下操作：
//将给定作业的优先级设置为pri，将其从由C保留的作业，等待延迟秒，然后将该作业放入就绪队列，使其可供任何客户端预订。
func (bs *Beanstalk) Release(id uint64, pri uint32, delay int) error {
	request, err := bs.cmd("release", id, pri, delay)
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "RELEASED")
	return err
}
func (bs *Beanstalk) UseRelease(tube string, id uint64, pri uint32, delay int) error {
	err := bs.Use(tube)
	if err != nil {
		return err
	}
	return bs.Release(id, pri, delay)
}

// Bury 将任务状态设置为BURIED，并设置优先级。BURIED状态在未改变前，不会被消费
func (bs *Beanstalk) Bury(id uint64, pri uint32) error {
	request, err := bs.cmd("bury", id, pri)
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "BURIED")
	return err
}

func (bs *Beanstalk) UseBury(tube string, id uint64, pri uint32) error {
	err := bs.Use(tube)
	if err != nil {
		return err
	}
	return bs.Bury(id, pri)
}

// Touch 允许worker请求更多的时间执行job，这个对长时间执行的任务有用
func (bs *Beanstalk) Touch(id uint64) error {
	request, err := bs.cmd("touch", id)
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "TOUCHED")
	return err
}

// UseTouch 切换tube, 并允许worker请求更多的时间执行job，这个对长时间执行的任务有用
func (bs *Beanstalk) UseTouch(tube string, id uint64) error {
	err := bs.Use(tube)
	if err != nil {
		return err
	}
	return bs.Touch(id)
}

// Stats 返回整个消息队列系统的整体信息
func (bs *Beanstalk) Stats() (map[string]string, error) {
	request, err := bs.cmd("stats")
	if err != nil {
		return nil, err
	}
	body, err := bs.readBody(request, "OK")
	return parseDict(body), err
}

// StatsJob 统计job的相关信息
func (bs *Beanstalk) StatsJob(id uint64) (map[string]string, error) {
	request, err := bs.cmd("stats-job", id)
	if err != nil {
		return nil, err
	}
	body, err := bs.readBody(request, "OK")
	return parseDict(body), err
}

// UseStatsJob 切换tube, 并统计job状态
func (bs *Beanstalk) UseStatsJob(tube string, id uint64) (map[string]string, error) {
    err := bs.Use(tube)
	if err != nil {
		return nil, err
	}
	return bs.StatsJob(id)
}

// ListTubes 列出所有存在的tube
func (bs *Beanstalk) ListTubes() ([]string, error) {
	request, err := bs.cmd("list-tubes")
	if err != nil {
		return nil, err
	}
	body, err := bs.readBody(request, "OK")
	return parseList(body), err
}

// ListTubeUsed 列出当前在使用的 tube
func (bs *Beanstalk) ListTubeUsed() (tube string, err error) {
	request, err := bs.cmd("list-tube-used")
	if err != nil {
		return "", err
	}
	err = bs.readResponse(request, "USING %s", &tube)
	return tube, err
}

// ListTubes 列出所有存在的tube
func (bs *Beanstalk) ListTubeWatched() ([]string, error) {
	request, err := bs.cmd("list-tubes-watched")
	if err != nil {
		return nil, err
	}
	body, err := bs.readBody(request, "OK")
	return parseList(body), err
}

// Peek 从服务器获取指定作业的副本        peek <id>\r\n
func (bs *Beanstalk) Peek(id uint64) (jobid uint64, body []byte, err error) {
	request, err := bs.cmd("peek", id)
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readBody(request, "FOUND %d", &jobid)
	if err != nil {
		return 0, nil, err
	}
	return jobid, body, nil
}

func (bs *Beanstalk) UsePeek(tube string, id uint64) (jobid uint64, body []byte, err error) {
	err = bs.Use(tube)
	if err != nil {
		return 0, nil, err
	}
	return bs.Peek(id)
}

// PeekReady 返回下一个ready job                peek-ready\r\n
func (bs *Beanstalk) PeekReady() (id uint64, body []byte, err error) {
	request, err := bs.cmd("peek-ready")
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readBody(request, "FOUND %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

func (bs *Beanstalk) UsePeekReady(tube string) (id uint64, body []byte, err error) {
	err = bs.Use(tube)
	if err != nil {
		return 0, nil, err
	}
	return bs.PeekReady()
}

// PeekDelayed 返回下一个延迟剩余时间最短的job           peek-delayed\r\n
func (bs *Beanstalk) PeekDelayed() (id uint64, body []byte, err error) {
	request, err := bs.cmd("peek-delayed")
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readBody(request, "FOUND %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

func (bs *Beanstalk) UsePeekDelayed(tube string) (id uint64, body []byte, err error) {
	err = bs.Use(tube)
	if err != nil {
		return 0, nil, err
	}
	return bs.PeekDelayed()
}

// PeekBuried 返回下一个在buried列表中的job        peek-buried\r\n
func (bs *Beanstalk) PeekBuried() (id uint64, body []byte, err error) {
	request, err := bs.cmd("peek-buried")
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readBody(request, "FOUND %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

func (bs *Beanstalk) UsePeekBuried(tube string) (id uint64, body []byte, err error) {
	err = bs.Use(tube)
	if err != nil {
		return 0, nil, err
	}
	return bs.PeekBuried()
}

// Kick 此指令应用在当前使用的tube中，它将job的状态迁移为ready或者delayed
func (bs *Beanstalk) Kick(bound int) (n int64, err error) {
	request, err := bs.cmd("kick", bound)
	if err != nil {
		return 0, err
	}
	err = bs.readResponse(request, "KICKED %d", &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (bs *Beanstalk) UseKick(tube string, bound int) (n int64, err error) {
	err = bs.Use(tube)
	if err != nil {
		return 0, err
	}
	return bs.Kick(bound)
}

// Stats 统计tube的相关信息
func (bs *Beanstalk) StatsTube(tube string) (map[string]string, error) {
	request, err := bs.cmd("stats-tube", tube)
	if err != nil {
		return nil, err
	}
	body, err := bs.readBody(request, "OK")
	return parseDict(body), err
}

// Pause 暂停任务一段时间
func (bs *Beanstalk) PauseTube(tube string, delay int) error {
	request, err := bs.cmd("pause-tube", tube, delay)
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "PAUSED")
	if err != nil {
		return err
	}
	return nil
}

// Reserve 取出（预订）job，待处理。
// 它将返回一个新预订的job，如果没有job，beanstalkd将直到有job时才发送响应。
// 取出job时，状态迁移为reserved, client被限制在指定的ttr时间内完成，否则超时，状态迁移为ready。
func (bs *Beanstalk) Reserve(timeout int) (id uint64, body []byte, err error) {
	request, err := bs.cmd("reserve-with-timeout", timeout)
	if err != nil {
		return 0, nil, err
	}
	body, err = bs.readBody(request, "RESERVED %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

func (bs *Beanstalk) Quit() error {
	request, err := bs.cmd("quit")
	if err != nil {
		return err
	}
	err = bs.readResponse(request, "quit")
	return err
}
