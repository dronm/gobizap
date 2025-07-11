package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	sess "github.com/dronm/session"

	"github.com/dronm/gobizap/v2/logger"
	"github.com/dronm/gobizap/v2/session"
)

var IsProduction bool
var MaxLifeTime int64

func extractCookieSessionID(c *gin.Context) string {
	// cookie value if exists
	if vCookie, err := c.Cookie(session.SESS_COOKIE_KEY); err == nil && vCookie != "" {
		logger.Logger.Debug("got cookie ", vCookie)
		return vCookie
	}
	return ""
}

func extractHeaderSessionID(c *gin.Context) string {
	// cookie valuer
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			c.Next()
			return
		}
		sessID := extractCookieSessionID(c)
		if sessID == "" {
			sessID = extractHeaderSessionID(c) // for API calls without cookies
		}

		var sessionAlreadyExists bool
		if sessID != "" {
			sessionAlreadyExists = true
		}

		c.Set("session_loader", func() (sess.Session, error) {
			fmt.Println("SessionMiddleware session_loader, starting for sessID:", sessID)
			sess, err := session.SessManager.SessionStart(sessID)
			if err != nil {
				return nil, err
			}

			if !sessionAlreadyExists {
				now := time.Now()
				_ = sess.Put("time_created", now)
			}

			fmt.Println("setting cookie to a session ID:", sess.SessionID())
			// set cookie
			c.SetCookie(
				session.SESS_COOKIE_KEY,
				sess.SessionID(),
				int(MaxLifeTime),
				"/",
				"",
				IsProduction,
				true,
			)
			c.Set("session", sess)
			return sess, nil
		})

		c.Next()
	}
}
