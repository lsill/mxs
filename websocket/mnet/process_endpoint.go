package mnet

import (
	"bufio"
	"github.com/astaxie/beego/logs"
	"io"
	"syscall"
	"time"
)

type ProcessEndPoint struct {
	process *LaunchedProcess
	closetime time.Duration
	output chan []byte
	bin bool
}

func NewProcessEndpoint(process *LaunchedProcess, bin bool) *ProcessEndPoint {
	return &ProcessEndPoint{
		process:   process,
		output:    make(chan []byte),
		bin:       bin,
	}
}

func (pe *ProcessEndPoint) Terminate() {
	terminated := make(chan struct{})
	go func() {
		pe.process.cmd.Wait()
		terminated <- struct{}{}
	}()
	// 对于某些过程，这足以完成它们...
	pe.process.stdin.Close()
	// 一点冗长，以创建良好的调试线索
	select {
	case <- terminated:
		logs.Debug("process %v terminated after stdin was closed", pe.process.cmd.Process.Pid)
		return // 进程结束
	case <- time.After(100 * time.Millisecond + pe.closetime):
	}
	err := pe.process.cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		// 没有这个错误就代表完成了
		logs.Error("sigint unsuccessful to %v: %v", pe.process.cmd.Process.Pid, err)
	}
	select {
	case <- terminated:
		logs.Debug("process %v terminated after sigint", pe.process.cmd.Process.Pid)
		return // 进行结束
	case <- time.After(250*time.Millisecond + pe.closetime):
	}
	err = pe.process.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		// 没有这个错误就代表完成了
		logs.Error("process SIGTERM unsuccessful to %v: %v", &pe.process.cmd.Process.Pid, err)
	}
	select {
	case <-terminated:
		logs.Debug("process %v terminated after SIGTERM", pe.process.cmd.Process.Pid)
		return
	case <-time.After(500*time.Millisecond + pe.closetime):
	}
	err = pe.process.cmd.Process.Kill()
	if err != nil {
		logs.Error("process SIGKILL unsuccessful to %v: %v", pe.process.cmd.Process.Pid, err)
		return
	}
	select {
	case <-terminated:
		logs.Debug(" Process %v terminated after SIGKILL", pe.process.cmd.Process.Pid)
		return // means process finished
	case <-time.After(1000 * time.Millisecond):
	}
}

func (pe *ProcessEndPoint) Output() chan []byte{
	return pe.output
}

func (pe *ProcessEndPoint) Send(msg []byte) bool {
	pe.process.stdin.Write(msg)
	return true
}

func (pe *ProcessEndPoint) StartReading() {
	go pe.log_stderr()
	if pe.bin {
		go pe.process_binout()
	} else {
		go pe.process_txtout()
	}
}

func (pe *ProcessEndPoint) process_txtout() {
	bufin := bufio.NewReader(pe.process.stdout)
	for {
		buf, err := bufin.ReadBytes('\n')
		if err !=nil{
			if err != io.EOF {
				logs.Error("process Unexpected error while reading STDOUT from process %v", err)
			} else {
				logs.Debug("process stdout closed")
			}
			break
		}
		pe.output <- trimEOL(buf)
	}
	close(pe.output)
}

func (pe *ProcessEndPoint) process_binout() {
	buf := make([]byte, 10 * 1024* 1024)
	for {
		n, err := pe.process.stdout.Read(buf)
		if err != nil {
			if err != io.EOF {
				logs.Error("process Unexpected error while reading STDOUT from process:%v",err)
			}else {
				logs.Debug("process stdout closed")
			}
			break
		}
		pe.output <- append(make([]byte, 0, n), buf[:n]...)
	}
	close(pe.output)
}

func (pe *ProcessEndPoint) log_stderr() {
	bufstderr := bufio.NewReader(pe.process.stderr)
	for {
		buf ,err:= bufstderr.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				logs.Error("process Unexpected error while reading STDERR from process: %s", err)
			}else {
				logs.Debug("process stderr closed")
			}
			break
		}
		logs.Error("stderr %s", string(trimEOL(buf)))
	}
}

// 从字符串中剪切unixy样式\ n和windowsy样式\ r \ n后缀
func trimEOL(b []byte) []byte {
	lns := len(b)
	if lns > 0 && b[lns - 1] == '\n' {
		lns--
		if lns > 0 && b[lns - 1] == '\r' {
			lns--
		}
	}
	return b[:lns]
}
