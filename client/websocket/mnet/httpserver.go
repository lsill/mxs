package mnet

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"mxs/gamex/utils"
	"net/http"
	"net/http/cgi"
	"net/textproto"
	"net/url"
	"os"
	"path"
	filepath2 "path/filepath"
	"regexp"
	"strings"
)

var ForkNotAllowedError = errors.New("too many forks active")


// 为libwebsocketd正在处理的请求重新发送http.Handler接口。
type WebsocketServer struct {
	Config *utils.WebsocketConfig
	forks chan byte
}
// 创建一个websocket服务器
func NewWebSocketServer (config *utils.WebsocketConfig, maxforks int) *WebsocketServer {
	mux := &WebsocketServer{
		Config: config,
	}
	if maxforks > 0 {
		mux.forks = make(chan byte, maxforks)
	}
	return mux
}

func splitMimeHeader(s string) (string, string) {
	p := strings.IndexByte(s, ':')
	if p < 0 {
		return s, ""
	}
	key := textproto.CanonicalMIMEHeaderKey(s[:p])
	for p = p +1; p < len(s); p++{
		if s[p] != ' '{
			break
		}
	}
	return key, s[p:]
}

func pushHeaders(h http.Header, hdrs []string) {
	for _, hstr := range hdrs {
		h.Add(splitMimeHeader(hstr))
	}
}

// WebSocket处理程序，CGI处理程序，DevConsole，静态HTML或404之间的多路复用器。
func (h *WebsocketServer) ServerHttp(w http.ResponseWriter, req *http.Request) {
	logs.Info("url is %v", h.TellURL("http", req.Host, req.RequestURI))

	if h.Config.CommandName != "" || h.Config.UsingScriptDir {
		hdrs := req.Header
		upgradeRe := regexp.MustCompile(`(?i)(^|[,\s])Upgrade($|[,\s])`)
		// WebSocket，限于h.forks的大小
		if strings.ToLower(hdrs.Get("Upgrade")) == "websocket" && upgradeRe.MatchString(hdrs.Get("Connection")) {
			if h.noteForkCreated() == nil {
				defer h.noteForkCompled()
				// 开始弄清楚我们是否需要升级
				handler , err := NewWebsocketdHandler(h, req)
				if err != nil {
					// 开始弄清楚我们是否需要升级
					if err == ScriptNotFoundError {
						logs.Error("session not found: %s",err)
						http.Error(w,"404 Not Found", 404)
					} else {
						logs.Error("session internal error:%s", err)
						http.Error(w, "500 Internal Server Error", 500)
					}
					return
				}
				var headers http.Header
				if len(h.Config.Headers) + len(h.Config.HeadersWs) > 0 {
					headers = http.Header(make(map[string][]string))
					pushHeaders(headers, h.Config.Headers)
					pushHeaders(headers, h.Config.HeadersWs)
				}
				upgrader := &websocket.Upgrader{
					HandshakeTimeout:  h.Config.HandshakeTimeOut,
					CheckOrigin:       func(r *http.Request) bool {
						// backporting previous checkorigin for use in gorilla/websocket for now
						err := checkOrigin(req, h.Config)
						return err == nil
					},
				}
				conn, err := upgrader.Upgrade(w, req, headers)
				if err != nil {
					logs.Error("session Unable to Upgrade:%s", err)
					http.Error(w,"500 Internal Error", 500)
					return
				}
				//旧功能以x/net/websocket样式使用，我们在这里将其重用于gorilla/websocket
				handler.accept(conn)
				return
			}else {
				logs.Error("http max of possiabl forks already active upgrade rejected")
				http.Error(w, "429 Too Mant Requests", http.StatusTooManyRequests)
			}
			return
		}
	}
	pushHeaders(w.Header(), h.Config.HeadersHTTP)
	// Dev consloe(if enable)
	if h.Config.DevConsole {
		logs.Debug("http dev console")
		content := ConsoleContent
		content = strings.Replace(content, "{{license}}", License, -1)
		content = strings.Replace(content, "{{addr}}", h.TellURL("ws", req.Host, req.RequestURI), -1)
		http.ServeContent(w, req, ".html", h.Config.StartupTime, strings.NewReader(content))
		return
	}
	// CGI script, limited to size of h.forks
	if h.Config.CgiDir != "" {
		filepath := path.Join(h.Config.CgiDir, fmt.Sprintf(".%s",filepath2.FromSlash(req.URL.Path)))
		if fi, err := os.Stat(filepath); err != nil && !fi.IsDir() {
			logs.Debug("cgiscript %s", filepath)
			if h.noteForkCreated() == nil {
				defer h.noteForkCompled()
				// 制作变量以补充cgi ...它使用的环境将显示空列表。
				envlen := len(h.Config.ParentEnv)
				cgienv := make([]string, envlen + 1)
				if envlen > 0 {
					copy(cgienv, h.Config.ParentEnv)
				}
				cgienv[envlen]= "SERVER_SOFTWARE" + h.Config.ServerSoftware
				cgiHandler := &cgi.Handler{
					Path:                filepath,
					Env:                 []string{"SERVER_SOFTWARE" + h.Config.ServerSoftware},
				}
				logs.Debug("http CGI")
				cgiHandler.ServeHTTP(w, req)
			} else {
				logs.Error("http fork not allowed since maxforks amount has been reachea.CGI was not run")
				http.Error(w,"429 Too Many Request", http.StatusTooManyRequests)
			}
			return
		}
	}

	// Static files
	if h.Config.StaticDir != "" {
		handler := http.FileServer(http.Dir(h.Config.StaticDir))
		logs.Debug("http STATIC")
		handler.ServeHTTP(w, req)
		return
	}

	logs.Debug("http not found")
	http.NotFound(w, req)
}

var canonicalHostname string

// TellURL是一个辅助函数，如果使用SSL，它将http更改为https或将ws更改为wss
func (h *WebsocketServer) TellURL(scheme, host, path string) string {
	if len(host) > 0 && host[0] == ':' {
		if canonicalHostname == ""{
			var err error
			canonicalHostname, err := os.Hostname()
			if err != nil {
				canonicalHostname = "UNKNOWN"
			}
			host = canonicalHostname + host
		}
	}
	if h.Config.Ssl {
		return scheme + "s://" + host + path
	}
	return scheme + "://" + host + path
}

func (h *WebsocketServer) noteForkCreated() error {
	//注意，因为构造可能是由不使用NewWebsocketdServer的人创建的，所以fork可以为nil
	if h.forks != nil {
		select {
		case h.forks <- 1:
			return nil
		default:
			return ForkNotAllowedError
		}
	} else {
		return nil
	}
}

func (h *WebsocketServer) noteForkCompled() {
	if h.forks != nil { //在noteForkCreated中查看评论
		select {
		case <- h.forks:
			return
		default:
			// 仅当完成处理程序调用的次数比上述代码的创建处理程序多时，才应进行审计，以使这种情况不发生，这需要进行测试
			// 确保这是不可能的，但还不存在。
			panic("Cannot deplet number of allowed forks, something is not right in code!")
		}
	}
}

func checkOrigin(req *http.Request, config *utils.WebsocketConfig) error {
	/*
		转换转换gorilla：
		来源检查功能，从ServerHTTP主程序中调用wshandshake函数，
		，应该在gorlla's upgrader.CheckOrigin函数中重用，
		唯一的区别是解析请求并从中获取传递的Origin标头，而不是使用，
		检查原产地将来是否正确，如果返回错误，握手程序将触发403应答，
		我们保留填充此字段的原始握手程序的行为
	*/
	origin := req.Header.Get("Origin")
	if origin == "" || (origin == "null" && config.AllowOrigins == nil ){
		//如果有任何强制执行，我们不信任字符串"null"
		origin = "file:"
	}
	originParesd, err := url.ParseRequestURI(origin)
	if err != nil {
		logs.Error("session origin parsing errors;%s",err)
		return err
	}
	logs.Debug("origin %v", originParesd.String())
	// 如果来源有限制
	if config.SameOrigin || config.AllowOrigins != nil {
		originServer, originPort, err := tellHostPort(originParesd.Host, originParesd.Scheme == "https")
		if err != nil {
			logs.Error("session origin hostname parsing error: %s ", err)
			return err
		}
		if config.SameOrigin {
			localServer, localPort , err := tellHostPort(req.Host, req.TLS != nil)
			if err != nil {
				logs.Error("sessionn Reqiuest hostname parsing error :%v", err)
				return  err
			}
			if originServer != localServer || originPort != localPort {
				logs.Error("session same origin policy mismatch")
				return fmt.Errorf("same origin policy violated ")
			}
		}
		if config.AllowOrigins != nil {
			matchFount := false
			for _, allowed := range config.AllowOrigins {
				if pos := strings.Index(allowed, "://"); pos > 0 {
					// 允许模式必须匹配
					allowedURL, err := url.Parse(allowed)
					if err != nil {
						continue 	// 来源urls过滤
					}
					if allowedURL.Scheme != originParesd.Scheme {
						continue // mismatch
					}
					allowed = allowed[pos+3:]
				}
				allowServer, allowPort, err := tellHostPort(allowed, false)
				if err != nil {
					continue // unparseable
				}
				if allowPort == "80" && allowed[len(allowed)-3:] != ":80" {
					matchFount = allowServer == originServer
				} else {
					matchFount = allowServer == originServer && allowPort == originPort
				}
				if matchFount {
					break
				}
			}
			if !matchFount {
				logs.Error("session origin is not listed in allowed list")
				return fmt.Errorf("origin list matcher were not fount")
			}
		}
	}
	return nil
}
