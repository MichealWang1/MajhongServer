package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/netip"
	"strings"
)

func getIp(ctx *gin.Context) string {
	for _, headerName := range []string{"Cf-Connecting-Ip", "X-Forwarded-For", "X-Real-Ip"} {
		ip, valid := validateIP(ctx.Request.Header.Get(headerName))
		if valid {
			return ip
		}
	}

	rIp := ctx.RemoteIP()
	if rIp == "" {
		return "127.0.0.1"
	}
	return rIp
}

func validateIP(ips string) (clientIp string, valid bool) {
	if ips == "" {
		return "", false
	}

	items := strings.Split(ips, ",")
	for i, value := range items {
		value = strings.TrimSpace(value)
		ip := net.ParseIP(value)
		if ip == nil {
			return "", false
		}

		if i == 0 {
			clientIp = value
			valid = true
		}
	}

	return
}

func RealIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, port, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr))
		if err != nil {
			port = "0"
		}

		clientIp := getIp(ctx)
		addr, err := netip.ParseAddr(clientIp)
		if err != nil {
			return
		}

		if addr.Is6() {
			ctx.Request.RemoteAddr = fmt.Sprintf("[%s]:%s", getIp(ctx), port)
		} else {
			ctx.Request.RemoteAddr = fmt.Sprintf("%s:%s", getIp(ctx), port)
		}
	}
}
