package httpSrv

import(
	"net/http"
	"net"
	"time"
	"io"
	"bytes"
	"net/url"
	
	"github.com/dronm/gobizap/socket"
)

type HTTPSocket struct {
	socket.ClientSocket
	Response http.ResponseWriter
	Request *http.Request
	ControllerID string
	MethodID string
	QueryParams url.Values //all unparsed params. Can be used in views without specifying them in methods, like f_name for ViewExcel
	
	//TransformTemplateID string //transformation template, templ parameter ViewHTML specific parameter
	TransformClassID string //v parameter
	ViewTemplateID string //t parameter
}

func (s *HTTPSocket) GetDescr() string {
	return ""
}

func (s *HTTPSocket) Close() {
}

func (s *HTTPSocket) GetConn() net.Conn {
	return nil
}
/*
func (s *HTTPSocket) GetIP() string{
	if s.Request == nil {
		return ""
	}
	return socket.GetRemoteAddrIP(s.Request.RemoteAddr)
}
*/
func (s *HTTPSocket) GetUploadedFileData(formField string) ([]byte, string, error) {
	//f := http_sock.Request.MultipartForm
	//for k, _ := range f.File {
	//	s.Request.FormFile(k)
	file, file_h, err := s.Request.FormFile(formField)
	if err != nil {
		return nil, "", err
	}
	//file_h.Filename, file_h.Size, file_h.Header)		
	defer file.Close()
	
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, "", err
	}		
	return buf.Bytes(), file_h.Filename, nil 
}

func NewHTTPSocket(w http.ResponseWriter, r *http.Request) *HTTPSocket{	
	return &HTTPSocket{ClientSocket: socket.ClientSocket{LastActivity: time.Now()},
			Response: w,
			Request: r,
	}
}
