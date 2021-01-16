package mnet

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"mxs/client/websocket/iface"
	"mxs/gamex/utils"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ScriptNotFoundError = errors.New("script not found")

// RemoteInfo包含有关远程http客户端的信息
type RemoteInfo struct {
	Addr, Host, Port string
}

// URLInfo-包含有关当前请求的信息的结构，它映射到文件系统
type URLInfo struct {
	ScriptPath string
	PathInfo string
	FilePath string
}


// WebsocketdHandler是单个请求信息和处理结构，它处理守护程序无法处理的所有WS请求（静态，cgi，devconsole）
type WebsocketdHandler struct {
	server *WebsocketServer

	Id string
	*RemoteInfo
	*URLInfo // TODO: I cannot find where it's used except in one single place as URLInfo.FilePath
	Env      []string

	command string
}

// NewWebsocketdHandler构造该结构并解析其中的所有必需内容...
func NewWebsocketdHandler(ws *WebsocketServer, req *http.Request) (wsh *WebsocketdHandler, err error) {
	wsh = &WebsocketdHandler{
		server: ws,
		Id:     generateId(),
	}
	wsh.RemoteInfo, err = GetRemoteInfo(req.RemoteAddr, ws.Config.ReverseLookup)
	if err != nil {
		logs.Error("session Could not understand remote address '%s': %s", req.RemoteAddr, err)
		return nil, err
	}
	logs.Debug("remote %v", wsh.RemoteInfo.Host)

	wsh.URLInfo, err = GetURLInfo(req.URL.Path, ws.Config)
	if err != nil {
		logs.Error("session not fount：%v", err)
		return nil, err
	}
	wsh.command = ws.Config.CommandName
	if ws.Config.UsingScriptDir {
		wsh.command = wsh.URLInfo.FilePath
	}
	logs.Debug("command is %s", wsh.command)
	wsh.Env = createEnv(wsh, req)

	logs.Debug("WebsocketdHandler id is %s", wsh.Id)

	return wsh, nil
}

func (wsh *WebsocketdHandler) accept(ws *websocket.Conn) {
	defer ws.Close()
	logs.Debug("session connection")

	launched ,err := launchCmd(wsh.command, wsh.server.Config.CommandArgs, wsh.Env)
	if err != nil {
		logs.Error("process could not launch process %s %s (%s)",wsh.command, strings.Join(wsh.server.Config.CommandArgs, " "), err)
		return
	}
	logs.Debug("pid", strconv.Itoa(launched.cmd.Process.Pid))

	binary := wsh.server.Config.Binary
	process := NewProcessEndpoint(launched, binary)
	if cms := wsh.server.Config.CloseMs; cms != 0 {
		process.closetime += time.Duration((cms)) * time.Millisecond
	}
	wsEndpoint := NewWebSocketEndPoint(ws, binary)

	iface.PipeEndPoints(process ,wsEndpoint)
}

// (考虑是修改连接id)
func generateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func GetRemoteInfo(remote string, doLookup bool) (*RemoteInfo, error) {
	addr, port, err := net.SplitHostPort(remote)
	if err !=nil {
		return nil, err
	}
	var host string
	if doLookup {
		hosts, err := net.LookupAddr(addr)
		if err != nil || len(hosts) == 0 {
			host = addr
		} else {
			host = hosts[0]
		}
	} else {
		host = addr
	}
	return &RemoteInfo{
		Addr: addr,
		Host: host,
		Port: port,
	}, nil
}

func GetURLInfo(path string, config *utils.WebsocketConfig) (*URLInfo, error){
	if !config.UsingScriptDir {
		return &URLInfo{
			ScriptPath: "/",
			PathInfo:   path,
			FilePath:   "",
		}, nil
	}
	parts := strings.Split(path[1:], "/")
	urlInfo := &URLInfo{}
	for i, part := range parts {
		urlInfo.ScriptPath = strings.Join([]string{urlInfo.ScriptPath, part}, "/")
		urlInfo.FilePath = filepath.Join(config.ScripDir, urlInfo.ScriptPath)
		isLastPart := i == len(parts) - 1
		startInfo, err := os.Stat(urlInfo.FilePath)

		// 不是一个有效的路径
		if err != nil {
			return nil, ScriptNotFoundError
		}

		// 在后面不是一个url而是一个dir
		if isLastPart && startInfo.IsDir() {
			return nil, ScriptNotFoundError
		}

		// 遇到目录，继续寻找
		if startInfo.IsDir() {
			continue
		}

		// 没有额外的参数
		if isLastPart {
			return urlInfo, nil
		}

		urlInfo.PathInfo = "/" + strings.Join(parts[i+1:], "/")
		return urlInfo, nil
	}
	panic(fmt.Sprintf("GetURLInfo cannot parse path %#v", path))
}

func tellHostPort(host string, ssl bool) (server, port string, err error) {
	server, port, err = net.SplitHostPort(host)
	if err != nil {
		if addrerr, ok := err.(*net.AddrError); ok && strings.Contains(addrerr.Err, "missing port") {
			server = host
			if ssl {
				port = "443"
			} else {
				port = "80"
			}
			err = nil
		}
	}
	return server, port, err
}