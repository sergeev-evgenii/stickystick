package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const viewerCookieName = "viewer_id"
const viewerCookieMaxAge = 365 * 24 * 3600 // 1 год

// getViewerKey возвращает ключ зрителя: по userID если авторизован, иначе по cookie viewer_id.
// Для анонимов при отсутствии cookie выставляет cookie и возвращает новый ключ.
func getViewerKey(c *gin.Context) string {
	if uid, ok := c.Get("userID"); ok {
		if id, ok := uid.(uint); ok {
			return "u:" + strconv.FormatUint(uint64(id), 10)
		}
	}
	val, err := c.Cookie(viewerCookieName)
	if err != nil || val == "" {
		b := make([]byte, 16)
		if _, _ = rand.Read(b); true {
			val = hex.EncodeToString(b)
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     viewerCookieName,
			Value:    val,
			Path:     "/",
			MaxAge:   viewerCookieMaxAge,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
		})
	}
	return "a:" + val
}
