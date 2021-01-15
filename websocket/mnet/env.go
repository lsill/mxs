package mnet

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strings"
)
const (
	gatewayInterface = "websocketd-CGI/0.1"
)

var headerNewlineToSpace = strings.NewReplacer("\n", " ", "\r", " ")
var headerDashToUnderscore = strings.NewReplacer("-","_")

func createEnv(handler *WebsocketdHandler, req *http.Request) []string {
	headers := req.Header
	url := req.URL
	serverName, serverPort, err := tellHostPort(req.Host, handler.server.Config.Ssl)
	if err != nil {
		//这确实意味着我们无法检测到来自Host：标头的端口...仅继续使用“”，猜测是错误的。
		logs.Warn("env host port detection errors:%s", err)
		serverPort = ""
	}
	standardEnvCount := 20
	if handler.server.Config.Ssl {
		standardEnvCount += 1
	}
	parentLen := len(handler.server.Config.ParentEnv)
	env := make([]string, 0, len(headers) + standardEnvCount + parentLen + len(handler.server.Config.Env))

	// 该变量可以从外部重写
	env = append(env, "SERVER_SOFTWARE", handler.server.Config.ServerSoftware)

	parentStarts := len(env)
	env = append(env, handler.server.Config.ParentEnv...)

	// 重要--->添加标题？确保standardEnvCount（以上）是最新的。

	// 标准CGI规范标头。
	// 如 http://tools.ietf.org/html/rfc3875 中所定义
	env = appendEnv(env, "REMOTE_ADDR",handler.RemoteInfo.Addr)
	env = appendEnv(env, "REMOTE_HOST", handler.RemoteInfo.Host)
	env = appendEnv(env, "SERVER_NAME", serverName)
	env = appendEnv(env, "SERVER_PORT", serverPort)
	env = appendEnv(env, "SERVER_PROTOCOL", req.Proto)
	env = appendEnv(env, "GATEWAY_INTERFACE", gatewayInterface)
	env = appendEnv(env, "REQUEST_METHOD", req.Method)
	env = appendEnv(env, "SCRIPT_NAME", handler.URLInfo.ScriptPath)
	env = appendEnv(env, "PATH_INFO", handler.URLInfo.PathInfo)
	env = appendEnv(env, "PATH_TRANSLATED", url.Path)
	env = appendEnv(env, "QUERY_STRING", url.RawQuery)

	// 不支持，但是我们明确清除了它们，因此我们不会从父环境中泄漏。
	env = appendEnv(env, "AUTH_TYPE", "")
	env = appendEnv(env, "CONTENT_LENGTH", "")
	env = appendEnv(env, "CONTENT_TYPE", "")
	env = appendEnv(env, "REMOTE_IDENT", "")
	env = appendEnv(env, "REMOTE_USER", "")

	// 非标准但常用的标头。
	env = appendEnv(env, "UNIQUE_ID", handler.Id) // Based on Apache mod_unique_id.
	env = appendEnv(env, "REMOTE_PORT", handler.RemoteInfo.Port)
	env = appendEnv(env, "REQUEST_URI", url.RequestURI()) // e.g. /foo/blah?a=b

	//以下变量是CGI规范的一部分，但是可选的而不是由websocketd设置：
	// AUTH_TYPE，REMOTE_USER，REMOTE_IDENT -身份验证留给基础程序。
	// CONTENT_LENGTH，CONTENT_TYPE -对于WebSocket连接没有意义。
	// SSL_ * -不支持SSL变量，为使用--ssl运行的websocketd添加了HTTPS = on

	if handler.server.Config.Ssl {
		env = appendEnv(env, "HTTPS", "on")
	}

	for i, v := range env {
		if i >= parentStarts && i < parentLen+parentStarts {
			logs.Debug("env Parent envvar: %v", v)
		} else {
			logs.Debug("env Std. variable: %v", v)
		}
	}


	for k, hdrs := range headers {
		header := fmt.Sprintf("HTTP_%s", headerDashToUnderscore.Replace(k))
		env = appendEnv(env, header, hdrs...)
		logs.Debug("env external variable:%s", env[len(env)- 1])
	}

	for _, v := range handler.server.Config.Env {
		env = appendEnv(env, v)
		logs.Debug("env external variable: %s", v)
	}

	return env
}

// 改编自 net/http/header.go
func appendEnv(env []string, k string, v ...string) []string {
	if len(v) == 0{
		return env
	}
	vCleaned := make([]string, 0, len(v))
	for _, val := range v {
		vCleaned = append(vCleaned, strings.TrimSpace(val))
	}
	return append(env, fmt.Sprintf("%s=%s",
		strings.ToUpper(k),
		strings.Join(vCleaned,", ")))
}