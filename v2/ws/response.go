package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dronm/crudifier"
	"github.com/dronm/gobizap/v2/errs"
	"github.com/dronm/gobizap/v2/logger"
	"github.com/gorilla/websocket"
)

type SrvResponse struct {
	EventID string            `json:"event_id"`
	QueryID string            `json:"query_id"`
	Payload any               `json:"payload"`
	Error   *SrvResponseError `json:"error"`
}

type SrvResponseError struct {
	Code    errs.ErrorCode `json:"code"`
	Message string         `json:"message"`
}

func NewSrvResponseError(httpErr int, fnName string, isProduction bool, err error) *SrvResponseError {
	errText := fmt.Sprintf("%s: %v", fnName, err)

	// log real message here
	logger.Logger.Error(errText)

	resp := SrvResponseError{}

	var pubErr errs.PublicError
	var validErr *crudifier.ValidationError // all validation errors

	if errors.As(err, &pubErr) {
		resp.Message = pubErr.Error()
		resp.Code = pubErr.Code()

	} else if errors.As(err, &validErr) {
		resp.Message = validErr.Error()
		resp.Code = errs.ValidationFailed

	} else {
		switch httpErr {
		case http.StatusInternalServerError:
			resp.Code = errs.InternalError
		case http.StatusBadRequest:
			resp.Code = errs.BadRequest
		default:
			resp.Code = errs.UnknownError
		}
		if !isProduction {
			resp.Message = errText
		} else {
			resp.Message = errs.ErrorDescr(resp.Code)
		}
	}
	return &resp
}

func (s *WSServer) SendMessage(conn *websocket.Conn, resp *SrvResponse) {
	// resp.Status = (resp.Error != nil)

	respData, err := json.Marshal(resp)
	if err != nil {
		logger.Logger.Errorf("WSServer SendMessage json.Marshal(): %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, respData); err != nil {
		logger.Logger.Errorf("WSServer SendMessage conn.WriteMessage(): %v", err)
		return
	}
}
