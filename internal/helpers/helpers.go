package helpers

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reverse-proxy/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CreateReverseProxy(targetHost string) gin.HandlerFunc {
	targetURL, err := url.Parse(targetHost)
	if err != nil {
		log.Error().Msgf("Invalid proxy target: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Fix the host header to match the target
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func AddRP(domainName string, port string) bool {

	err := constants.Redis.HSet(context.Background(), "rp", domainName, port)
	if err != nil {
		return false
	}

	constants.RPCtxManager.AddRP(domainName, port)

	return true
}

func RemoveRP(domainName string) bool {
	err := constants.Redis.HDel(context.Background(), "rp", domainName)
	if err != nil {
		return false
	}

	constants.RPCtxManager.RemoveRP(domainName)

	return true
}

func LoadRedisContext() map[string]string {
	if constants.Redis == nil {
		return nil
	}

	res, err := constants.Redis.HGetAll(context.Background(), "rp")
	if err != nil {
		log.Error().Msgf("Error while loading Redis context -> %v", err)
		return nil
	}

	constants.RPCtxManager.LoadContext(res)

	return constants.RPCtxManager.GetContext()
}
