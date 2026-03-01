package interceptor

import (
	"strconv"
	"strings"
	"time"

	"github.com/NARUBROWN/spine/core"
	"github.com/NARUBROWN/spine/pkg/httpx"
)

type CookieInterceptor struct {
}

func NewCookieInterceptor() *CookieInterceptor {
	return &CookieInterceptor{}
}

func (i *CookieInterceptor) PreHandle(ctx core.ExecutionContext, _ core.HandlerMeta) error {
	cookie := strings.TrimSpace(ctx.Header("Cookie"))

	if cookie == "" {
		return nil
	}

	cookies := map[string]*httpx.Cookie{}
	isFirst := true

	for len(cookie) > 0 {
		var part string

		// 세미콜론 찾기
		if idx := strings.IndexByte(cookie, ';'); idx >= 0 {
			part = cookie[:idx]
			cookie = cookie[idx+1:]
		} else {
			part = cookie
			cookie = ""
		}

		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		key, value := part, ""
		if idxEqual := strings.IndexByte(part, '='); idxEqual >= 0 {
			key, value = part[:idxEqual], part[idxEqual+1:]
		}

		// 쿠키는 큰따옴표가 있을 수 있음 (RFC 6265)
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		c := httpx.Cookie{}

		if isFirst {
			c.Name = key
			c.Value = value
			isFirst = false
			continue
		}

		switch strings.ToLower(key) {
		case "path":
			c.Path = value
		case "domain":
			c.Domain = value
		case "max-age":
			if secs, err := strconv.Atoi(value); err == nil {
				c.MaxAge = secs
			}
		case "expires":
			if t, err := time.Parse(time.RFC1123, value); err == nil {
				c.Expires = &t
			} else if t, err := time.Parse("Mon, 02-Jan-2006 15:04:05 MST", value); err == nil {
				c.Expires = &t
			}
		case "httponly":
			c.HttpOnly = true
		case "secure":
			c.Secure = true
		case "samesite":
			switch strings.ToLower(value) {
			case "lax":
				c.SameSite = httpx.SameSiteLax
			case "strict":
				c.SameSite = httpx.SameSiteStrict
			case "none":
				c.SameSite = httpx.SameSiteNone
			}
		case "priority":
			c.Priority = value
		}

		cookies[key] = &c
	}

	ctx.Set("cookies", cookies)

	return nil
}

func (i *CookieInterceptor) PostHandle(core.ExecutionContext, core.HandlerMeta)             {}
func (i *CookieInterceptor) AfterCompletion(core.ExecutionContext, core.HandlerMeta, error) {}
