package httpSrv

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
	"encoding/json"
	"bytes"
	"os"
	"errors"
	"sync"
	
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/view"
	"github.com/dronm/gobizap/stat"	
	"github.com/dronm/gobizap/response"
	
)

//HTTP server, OnHandleRequest, must be defined

type CONTENT_DISPOSITION string
const (
	PARAM_TOKEN = "token"
	
	PARAM_CONTROLLER = "c"
	PARAM_METH = "f"
	PARAM_VIEW = "v"
	PARAM_VIEW_TMPL = "t" //view template to send with response, added in http_app
	PARAM_QUERY_ID = "query_id"
	
	CONTROLLER_QUERY_POSF = "_Controller"
	
	DEF_USER_TRANSFORM_CLASS_ID =  "ViewBase"
	DEF_GUEST_TRANSFORM_CLASS_ID =  "Login"
	
	DEF_MULTYPART_MAX_MEM = 256 //32 << 20 - 32mb
	
	CONTENT_DISPOSITION_ATTACHMENT CONTENT_DISPOSITION = "attachment"
	CONTENT_DISPOSITION_INLINE CONTENT_DISPOSITION = "inline"
	
	CHARSET_UTF8 = "charset=utf-8"
)

const HANDLER_TIMEOUT_TEXT = "handler timed out"

type OnBeforeHandleRequestProto func(socket.ClientSocketer)
type OnDefineUserTransformClassIDProto func(*HTTPSocket)

type URLShortcut struct {
	ControllerID string
	MethodID string
	ViewID string
	Params map[string]string
}

type HTTPServer struct {
	srv.BaseServer
	
	// Core http servers. Values may be initialized in an application
	// with finetuned parameters for timeouts. If nil servers are initialized
	// with no timeouts.
	Server *http.Server
	ServerTLS *http.Server
	
	// Defines a timeout for a handler. If not set, no timeout is used.
	TimeoutHandler int64
	
	Statistics stat.SrvStater
	
	OnBeforeHandleRequest OnBeforeHandleRequestProto
	OnDefineUserTransformClassID OnDefineUserTransformClassIDProto
	HTTPDir string	
	AllowedExtensions []string
	Headers map[string]string	
	URLShortcuts map[string]URLShortcut
	viewContentTypes map[string]string
	
	MultypartMaxMemory int64 //bytes
	
	idleConnsClosed chan struct{}
}

//controller ID, method ID, view ID
func (s *HTTPServer) AddURLShortcut(ID, cID, mID, vID string, params map[string]string) {
	if s.URLShortcuts == nil {
		s.URLShortcuts = make(map[string]URLShortcut)
	}
	s.URLShortcuts[ID] = URLShortcut{ControllerID: cID, MethodID: mID, ViewID: vID, Params: params}
}

func (s *HTTPServer) Run() {
	if s.OnHandleRequest == nil {
		s.Logger.Fatal("HTTPServer.OnHandleRequest not defined")
	}
	/*if s.OnHandlePermission != nil && s.OnHandleProhibError == nil {
		s.Logger.Fatal("HTTPServer.OnHandlePermission defined, but OnHandleProhibError not defined")
	}*/
	if s.OnHandleSession != nil && s.OnHandleServerError == nil {
		s.Logger.Fatal("HTTPServer.OnHandleSession defined, but OnHandleServerError not defined")
	}

	//TLS if nedded
	tls_start := (s.TlsAddress != "" && s.TlsCert != "" && s.TlsKey != "")
	ws_start := (s.Address!= "")
	
	s.Statistics = stat.NewSrvStat()
	
	s.idleConnsClosed = make(chan struct{})
	
	//https server: 1 process or gorouting
	if tls_start {
		if s.ServerTLS == nil {
			//defaults
			s.ServerTLS = &http.Server{Addr: s.TlsAddress}
			
		}else if s.Server.Addr == "" {
			s.Server.Addr = s.TlsAddress
		}
		
		if s.TimeoutHandler > 0 {
			s.Server.Handler = http.TimeoutHandler(http.HandlerFunc(s.HandleRequest), time.Duration(s.TimeoutHandler) * time.Millisecond, HANDLER_TIMEOUT_TEXT)
		}else{
			s.Server.Handler = http.HandlerFunc(s.HandleRequest)
		}
		
		s.Logger.Infof("https server: %s", s.Server.Addr)		
		if !ws_start {
			//the only server - main loop
			if err := s.ServerTLS.ListenAndServeTLS(s.TlsCert, s.TlsKey); err != nil && err != http.ErrServerClosed {
				s.Logger.Errorf("HTTPServer ListenAndServeTLS() failed: %v", err)
			}
		}else{
			//2 servers
			go func() {
				if err := s.ServerTLS.ListenAndServeTLS(s.TlsCert, s.TlsKey); err != nil && err != http.ErrServerClosed {
					s.Logger.Errorf("HTTPServer ListenAndServeTLS() failed: %v", err)
				}
			}()
		}
	}
	
	//http server
	if ws_start {		
		if s.Server == nil {
			//defaults
			s.Server = &http.Server{Addr: s.Address}
			
		}else if s.Server.Addr == "" {
			s.Server.Addr = s.Address
		}
		
		if s.TimeoutHandler > 0 {
			s.Server.Handler = http.TimeoutHandler(http.HandlerFunc(s.HandleRequest), time.Duration(s.TimeoutHandler) * time.Millisecond, HANDLER_TIMEOUT_TEXT)
		}else{
			s.Server.Handler = http.HandlerFunc(s.HandleRequest)
		}
		
		s.Logger.Infof("http server: %s", s.Server.Addr)		
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Errorf("HTTPServer ListenAndServe() failed: %v", err)
		}
	}
	
	<-s.idleConnsClosed
}


func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.Server == nil && s.ServerTLS == nil {
		return nil
	}
	
	var wg sync.WaitGroup
	if s.Server != nil {
		wg.Add(1)
		go func(c context.Context, httpS *HTTPServer) {
			defer wg.Done()
			if err := httpS.Server.Shutdown(c); err != nil {
				// Error from closing listeners, or context timeout:
				httpS.Logger.Errorf("HTTPServer Shutdown() failed: %v", err)
			}else{
				httpS.Logger.Warn("http server shutdown gracefully")
			}	
		}(ctx, s)
	}
	if s.ServerTLS != nil {
		wg.Add(1)
		go func(c context.Context, httpS *HTTPServer) {
			defer wg.Done()
			if err := httpS.ServerTLS.Shutdown(c); err != nil {
				// Error from closing listeners, or context timeout:
				httpS.Logger.Errorf("HTTPServer Shutdown() failed: %v", err)
			}else{
				httpS.Logger.Warn("https server shutdown gracefully")
			}	
		}(ctx, s)
	}
	
	wg.Wait()
	
	close(s.idleConnsClosed)	
	
	return nil
}

//parses query params based on query method, queryParams always non-nil map
func (s *HTTPServer) parseQueryParams(r *http.Request, queryParams *url.Values) error {
	if r.Method == http.MethodGet {
		*queryParams = r.URL.Query()
	}else{
		ct := r.Header.Get("Content-Type")
		if strings.Contains(ct, "multipart/form-data") {
			var mem  int64
			if s.MultypartMaxMemory == 0 {
				mem = DEF_MULTYPART_MAX_MEM
			}else{
				mem = s.MultypartMaxMemory
			}
			if err := r.ParseMultipartForm(mem); err != nil {
				return err
			}
			*queryParams = r.MultipartForm.Value
		}else{
			r.ParseForm()
			*queryParams = r.Form
		}
	}
	return nil
}

func (s *HTTPServer) checkExtension(ext string) bool {
	for _,s_ext := range s.AllowedExtensions {
		if ext == s_ext {
			return true
		}
	}
	return false
}

func (s *HTTPServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	sock := NewHTTPSocket(w, r)	
	if r.URL.Path != "/" {
		path_parts := strings.Split(r.URL.Path, "/")

		sh_cut_found := false
		if s.URLShortcuts != nil {
			path := r.URL.Path
			if path[len(path)-1:] == "/" {
				path = path[:len(path)-1]
			}
			if sh_cut, ok := s.URLShortcuts[path]; ok {
				//Shortcuts - predefined paths
				if err := s.parseQueryParams(r, &sock.QueryParams); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					s.Logger.Errorf("HTTPServer parseQueryParams(): %v", err)
					return
				}
							
				sock.QueryParams.Add(PARAM_CONTROLLER, sh_cut.ControllerID)
				sock.QueryParams.Add(PARAM_METH, sh_cut.MethodID)
				sock.QueryParams.Add(PARAM_VIEW, sh_cut.ViewID)
				sh_cut_found = true
			}
		}
		
		if !sh_cut_found {
			file_parts := strings.Split(path_parts[len(path_parts)-1], ".")		
			n := len(file_parts)
			if n > 0 && s.checkExtension(file_parts[n-1]) {
				//file serving
				if view.FileExists(s.HTTPDir + r.URL.Path) {
					http.ServeFile(w, r, s.HTTPDir + r.URL.Path)
				}else{
					s.Logger.Errorf("HTTPServer.OnHandleRequest %s file with extension %s not found", s.HTTPDir + r.URL.Path,file_parts[n-1])
				}
				return
			}
			
			if len(path_parts) >= 2 {		
				//schema: controller/method/view
				if err := s.parseQueryParams(r, &sock.QueryParams); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					s.Logger.Errorf("HTTPServer parseQueryParams: %v", err)
					return
				}
				
				sock.QueryParams.Add(PARAM_CONTROLLER, path_parts[1])
				if len(path_parts) >= 3 {
					sock.QueryParams.Add(PARAM_METH, path_parts[2])
				}	
				if len(path_parts) >= 4 {
					sock.QueryParams.Add(PARAM_VIEW, path_parts[3])
				}
			}else{
				
				//not found
				s.Logger.Errorf("HTTPServer.OnHandleRequest %s file with extension %s not found", s.HTTPDir + r.URL.Path,file_parts[n-1])
				w.WriteHeader(http.StatusNotFound)
				//+ if ViewHTML return NotFound page???
				return
			}
		}
		
	}else if err := s.parseQueryParams(r, &sock.QueryParams); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.Logger.Errorf("HTTPServer parseQueryParams: %v", err)
		return
	}
	
	sock.Token, sock.TokenExpires = extractParam(r, sock.QueryParams, PARAM_TOKEN)
	//check header with the same name
	if sock.Token == "" {
		sock.Token = r.Header.Get("token")
	}
	token_from_query := (sock.Token != "")

	//turn query/body parameters to json payload {"argv": {"par1":val1, "par2":val2}}
	var query_id, view_id string	
	meth_params := []byte(`{"argv":{`) //all other params
	meth_params_exists := false //if at least one parameter is already added
	
	for par_key, par_val:= range sock.QueryParams {
		if par_key == PARAM_CONTROLLER && len(par_val)>0 {
			sock.ControllerID = par_val[0]
			//extract postfix if any
			posf_pos := len(sock.ControllerID)-len(CONTROLLER_QUERY_POSF)
			if posf_pos > 0 && sock.ControllerID[posf_pos:] == CONTROLLER_QUERY_POSF {
				sock.ControllerID = sock.ControllerID[:posf_pos]
			}
		
		}else if par_key == PARAM_METH && len(par_val)>0 {
			sock.MethodID = par_val[0]
			
		}else if par_key == PARAM_VIEW && len(par_val)>0 {
			sock.TransformClassID = par_val[0]
			view_id = par_val[0] 

		}else if par_key == PARAM_VIEW_TMPL && len(par_val)>0 {
			sock.ViewTemplateID = par_val[0]
			
		}else if par_key == PARAM_QUERY_ID && len(par_val)>0 {
			query_id = par_val[0]

		}else if len(par_val)>0 {
			if meth_params_exists {
				meth_params = append(meth_params, ',')
			}else{
				meth_params_exists = true
			}
			var par_val_b []byte
			par_val_len := len(par_val[0])
			if par_val_len >= 2 &&
			( (par_val[0][0:1] == "{" && par_val[0][par_val_len-1:par_val_len] == "}") ||
			(par_val[0][0:1] == "[" && par_val[0][par_val_len-1:par_val_len] == "]") ) {
				//object!!!
				par_val_b = []byte(par_val[0])
			}else{
				//string
				//par_val_s := strings.ReplaceAll(par_val[0], `\n`, `\\n`)
				//par_val_s := strings.ReplaceAll(par_val[0], `"`, `\"`)				
				var err error
				par_val_b, err = json.Marshal(par_val[0])
				if err != nil {
					s.Logger.Errorf("HTTPServer json.Marshal(): %v", err)
					s.OnHandleServerError(s, sock, query_id, view_id)
				}
			}
			//, 
			meth_params = append(meth_params, '"')
			meth_params = append(meth_params, []byte(par_key)...)
			meth_params = append(meth_params, []byte{'"', ':'}...)
			meth_params = append(meth_params, par_val_b...)
		}
	}
	meth_params = append(meth_params, []byte{'}', '}'}...) //close argv
	
	//session
	if s.OnHandleSession != nil {		
		err := s.OnHandleSession(sock)
		if err != nil {
			s.Logger.Errorf("HTTPServer HandleRequest OnHandleSession: %v", err)
			s.OnHandleServerError(s, sock, query_id, view_id)
			return
		}

		if sock.Token == "" {
			//new session started
			sess := sock.GetSession()
			sock.Token = sess.SessionID()			
			s.Statistics.IncHandshakes()
		}
		
		if !token_from_query {			
			//sock.TokenExpires = sess.
			// Make sure the session cookie is not accessable via javascript.
			http.SetCookie(w, &http.Cookie{Name: PARAM_TOKEN,
					Value: sock.Token,
					HttpOnly: true,
					//Expires= sock.TokenExpires,
					//Path:
					//Domain
					//Expires
					//MaxAge
				})
		}		
	}
	
	if sock.TransformClassID == "" && s.OnDefineUserTransformClassID != nil {
		//handler is defined for View absence cases
		s.OnDefineUserTransformClassID(sock)
		
	}else if sock.TransformClassID == "" {
		//defaults
		defineUserTransformClassID(sock)
	}

	view_id = sock.TransformClassID
	if !view.Registered(view_id) {
		view_id = "ViewHTML"
	}	

	if query_id =="" {
		//http always expects result
		query_id = "1"
	}

	//header
	cont_tp := s.GetViewContentType(view_id)
	if cont_tp != "" {
		w.Header().Set("Content-Type", cont_tp)
	}else{
		s.Logger.Warnf("Content type for view %s not defined", view_id)
	}
	
	if s.Headers != nil {
		for key, val := range s.Headers {
			w.Header().Set(key, val)
		}
	}	
	if s.OnBeforeHandleRequest != nil {
		s.OnBeforeHandleRequest(sock)
	}
	s.OnHandleRequest(s, sock, sock.ControllerID, sock.MethodID, query_id, meth_params, view_id)
}

func (s *HTTPServer) SendToClient(sock socket.ClientSocketer, msg []byte) error {	
	if http_sock, ok := sock.(*HTTPSocket); ok {
		//fmt.Fprint(http_sock.Response, string(msg))
		if _, err:= http_sock.Response.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

//empty stub, ClientSockets is not used
func (s *HTTPServer) GetClientSockets() *socket.ClientSocketList{
	return nil
}

func (s *HTTPServer) AddViewContentType(viewID string, mimeType MIME_TYPE, charset string) {
	if s.viewContentTypes == nil {
		s.viewContentTypes = make(map[string]string)
	}
	s.viewContentTypes[viewID] = string(mimeType)
	if charset != "" {
		s.viewContentTypes[viewID] += "; "+charset
	}
}

func (s *HTTPServer) GetViewContentType(viewID string) string {
	if tp, ok := s.viewContentTypes[viewID]; ok {
		return tp
	}
	return ""
}

func (s *HTTPServer) AddFile(viewID string, mimeType MIME_TYPE, charset string) {
	if s.viewContentTypes == nil {
		s.viewContentTypes = make(map[string]string)
	}
	s.viewContentTypes[viewID] = string(mimeType)
	if charset != "" {
		s.viewContentTypes[viewID] += "; "+charset
	}
}

func (s *HTTPServer) GetStatistics() stat.SrvStater {
	return s.Statistics
}

func ServeContent(sock *HTTPSocket, fData *[]byte, fName string, mimetype MIME_TYPE, modTime time.Time, contDisposition CONTENT_DISPOSITION) {
	f_bytes := bytes.NewReader(*fData) // converted to io.ReadSeeker type
	sock.Response.Header().Set("Content-Type", string(mimetype))
	sock.Response.Header().Set("Content-Disposition", string(contDisposition)+";filename="+fName)
	http.ServeContent(sock.Response, sock.Request, fName, modTime, f_bytes)
}

//extractParam Retrieves parameter value and its expiration time, first looks into cookies, then QueryParams
func extractParam(r *http.Request, queryParams url.Values, param string) (string, time.Time) {
	//analyse cookies	
	if v_cookie, err := r.Cookie(param); v_cookie != nil && err == nil && v_cookie.Value != "" {
		return v_cookie.Value, v_cookie.Expires
	}	
	
	//if param is not present in cookies, trying to get it from query params
	if val_par, ok := queryParams[param]; ok && len(val_par)>0 {
		return val_par[0], time.Time{}
	}
	return "", time.Time{}
}


//mimetype default is GetMimeTypeOnFileExt()
//contDisposition (attachment|inline) default is attachment
func DownloadFile(resp *response.Response, sock socket.ClientSocketer, f *os.File, fName string, mimetype MIME_TYPE, contDisposition CONTENT_DISPOSITION) error {
	sock_http, ok := sock.(*HTTPSocket)
	if !ok {
		return errors.New("sock must be *HTTPSocket")
	}

	if mimetype == "" {
		mimetype = GetMimeTypeOnFileExt(fName)
	}
	if contDisposition == "" {
		contDisposition = CONTENT_DISPOSITION_ATTACHMENT
	}
	
	file_info, _ := f.Stat()
	f_size := file_info.Size()
	f_mod := file_info.ModTime()
	
	buffer := make([]byte, f_size)
	f.Read(buffer)
	
	ServeContent(sock_http, &buffer, fName, mimetype, f_mod, contDisposition)
	resp = nil
	
	return nil
}

//Sets default transformation class ID to sock.TransformClassID
//Uses session LOGGED variable to define different classes
func defineUserTransformClassID(sock *HTTPSocket) {
	sess := sock.GetSession()
	//login can actually be absent
	if sess == nil {
		sock.TransformClassID = DEF_GUEST_TRANSFORM_CLASS_ID
	}else if sess.GetBool("LOGGED") {
		sock.TransformClassID = DEF_USER_TRANSFORM_CLASS_ID
	}else{
		sock.TransformClassID = DEF_GUEST_TRANSFORM_CLASS_ID
	}
}

/*
func file_exists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		return true
	}
	return false
}
*/

