package controllers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/dronm/gobizap/v2/api"
	"github.com/dronm/gobizap/v2/database"
	"github.com/gin-gonic/gin"
)

// APIGet executes GET request. Parameters are passed in their original order
// of appearence in the URL.
func APIGet(c *gin.Context) {
	funcName := "APIGet"

	// params
	rawQuery := c.Request.URL.RawQuery
	pairs := strings.Split(rawQuery, "&")
	var params []string

	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		val := ""
		// key, _ = url.QueryUnescape(kv[0])
		if len(kv) > 1 {
			val, _ = url.QueryUnescape(kv[1])
		}
		params = append(params, val)
	}

	callAPIMethod(c, funcName, params)
}

// APIPost handle post requests. It only deals with app/json requests.
func APIPost(c *gin.Context) {
	funcName := "APIPost"

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ServeError(c, http.StatusInternalServerError, funcName+" io.ReadAll()", err)
		return
	}

	params, err := api.UnmarshalParams(bodyBytes)
	if err != nil {
		ServeError(c, http.StatusBadRequest, funcName+" "+err.Error(), err)
		return
	}

	// decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	//
	// tok, err := decoder.Token()
	// if err != nil {
	// 	ServeError(c, http.StatusBadRequest, funcName+" decoder.Token()", err)
	// 	return
	// }
	//
	// delim, ok := tok.(json.Delim)
	// if !ok || delim != '{' {
	// 	ServeError(c, http.StatusBadRequest, funcName+" tok.(json.Delim)", fmt.Errorf("expected json object"))
	// 	return
	// }
	//
	// for decoder.More() {
	// 	//Do not need key here!!!
	// 	// keyTok, err := decoder.Token()
	// 	// if err != nil {
	// 	// 	ServeError(c, http.StatusBadRequest, funcName+" decoder.Token()", fmt.Errorf("error reading key"))
	// 	// 	return
	// 	// }
	//
	// 	// key, ok := keyTok.(string)
	// 	// if !ok {
	// 	// 	ServeError(c, http.StatusBadRequest, funcName+" keyTok.(string)", fmt.Errorf("invalid key type"))
	// 	// 	return
	// 	// }
	// 	//
	// 	// Read value token
	// 	var raw json.RawMessage
	// 	if err := decoder.Decode(&raw); err != nil {
	// 		ServeError(c, http.StatusBadRequest, funcName+" decoder.Decode()", fmt.Errorf("error decoding value"))
	// 		return
	// 	}
	//
	// 	// Optionally: store raw JSON string or decode into string
	// 	// For simplicity, decode into string if possible
	// 	var strVal string
	// 	if err := json.Unmarshal(raw, &strVal); err != nil {
	// 		// If not a string, keep raw JSON text
	// 		strVal = string(raw)
	// 	}
	//
	// 	params = append(params, strVal)
	// }
	//
	// // ensure cloging
	// if _, err := decoder.Token(); err != nil {
	// 	ServeError(c, http.StatusBadRequest, funcName+" decoder.Token()", fmt.Errorf("error reading end of object"))
	// 	return
	// }

	callAPIMethod(c, funcName, params)
}

// callAPIMethod with the given params. funcName is a name of the calling function
// for the log.
func callAPIMethod(c *gin.Context, funcName string, params []string) {
	sess := GetSession(c, funcName)
	if sess == nil {
		return
	}

	service, method, err := extractService(c)
	if err != nil {
		ServeError(c, http.StatusBadRequest, funcName, err)
		return
	}
	srvMeth := fmt.Sprintf("%s.%s()", service, method)

	results, err := api.CallMethod(
		c.Request.Context(),
		service,
		method,
		params,
		&api.ServiceContext{DB: database.DB, Session: sess},
	)
	if err != nil {
		ServeError(c, http.StatusBadRequest, srvMeth, err)
		return
	}

	// last result is always an error
	var htmlResult any
	if len(results) > 0 {
		last := results[len(results)-1]
		errVal := results[len(results)-1]
		if last.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) && !last.IsNil() {
			err := errVal.Interface().(error)
			ServeError(c, http.StatusInternalServerError, srvMeth, err)
			return
		}
		if len(results) > 1 {
			if len(results) == 1 {
				// one model
				htmlResult = results[0].Interface()
			} else {
				// slice of models, minus error result
				res := make([]any, len(results)-1)
				for i := 0; i < len(results)-1; i++ {
					res[i] = results[i].Interface()
				}
				htmlResult = res
			}
		}
	}
	c.JSON(http.StatusOK, htmlResult)
}

// extractService is a helper functin to retrieve
// service and method from http request.
func extractService(c *gin.Context) (service string, method string, err error) {
	service = PascalCaseFromKebabCase(c.Param("service")) // kebab-cased
	if service == "" {
		err = fmt.Errorf("service is undefined")
		return
	}

	method = PascalCaseFromKebabCase(c.Param("method")) // kebab-cased
	if method == "" {
		err = fmt.Errorf("service method is undefined")
		return
	}

	return
}

func PascalCaseFromKebabCase(s string) string {
	if s == "" {
		return ""
	}
	var res strings.Builder
	parts := strings.Split(s, "-")
	for _, w := range parts {
		if len(w) > 0 {
			res.WriteString(strings.ToUpper(w[0:1]))
			if len(w) > 1 {
				res.WriteString(w[1:])
			}
		}
	}
	return res.String()
}

