package server

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zachmann/go-oidfed/pkg"

	"github.com/zachmann/offa/internal"
	"github.com/zachmann/offa/internal/cache"
	"github.com/zachmann/offa/internal/config"
	"github.com/zachmann/offa/internal/model"
)

func isTrustedIP(ipStr string) bool {
	if len(config.Get().Server.TrustedNets) == 0 {
		return true
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, ipNet := range config.Get().Server.TrustedNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func addAuthHandlers(s fiber.Router) {
	path := config.Get().Server.Paths.ForwardAuth
	s.Head(path, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	s.Options(path, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	s.All(
		path, func(c *fiber.Ctx) error {
			clientIP := c.IP()
			if !isTrustedIP(clientIP) {
				log.WithField("ip", clientIP).Info("Blocked untrusted IP")
				return c.Status(fiber.StatusForbidden).SendString("Forbidden: Untrusted Proxy")
			}

			forHost := c.Get(fiber.HeaderXForwardedHost)
			forPath := c.Get("X-Forwarded-Uri")
			if len(config.Get().Auth) == 0 {
				return c.Status(fiber.StatusForbidden).SendString("Forbidden")
			}
			rule := config.Get().Auth.FindRule(forHost, forPath)
			if rule == nil {
				return c.Status(fiber.StatusForbidden).SendString("Forbidden")
			}

			sessionToken := c.Cookies(config.Get().SessionStorage.CookieName)
			if sessionToken == "" {
				next := url.URL{
					Scheme: c.Get(fiber.HeaderXForwardedProto),
					Host:   forHost,
					Path:   forPath,
				}
				return c.Redirect(
					fmt.Sprintf("%s?next=%s", fullLoginPath, next.String()), http.StatusSeeOther,
				)
			}

			userInfos, err := validateSession(sessionToken)
			if err != nil || userInfos == nil {
				if err != nil {
					log.WithError(err).Info("Invalid session")
				}
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid session")
			}

			log.Debugf("auth request Userclaims are: %+v", userInfos)

			if !verifyUser(userInfos, rule.Require) {
				return c.Status(fiber.StatusForbidden).SendString("Forbidden")
			}

			setHeaders(c, rule.ForwardHeaders, userInfos)

			return c.SendStatus(fiber.StatusOK)
		},
	)
}

func verifyUser(
	claims model.UserClaims, require pkg.SliceOrSingleValue[map[model.Claim]pkg.SliceOrSingleValue[string]],
) bool {
	if len(require) == 0 {
		return true
	}
	for _, options := range require {
		var optionFailed bool
		for claim, claimRequires := range options {
			claimValue, ok := claims.GetString(claim)
			if ok {
				// string claim
				if len(claimRequires) != 1 {
					optionFailed = true
					break
				}
				if claimValue == claimRequires[0] {
					continue
				} else {
					optionFailed = true
					break
				}
			}
			claimValues, ok := claims.GetStringSlice(claim)
			if ok {
				if internal.SliceIsSubsetOf(claimRequires, claimValues) {
					continue
				} else {
					optionFailed = true
					break
				}
			}
			optionFailed = true
		}
		if !optionFailed {
			return true
		}
	}
	return false
}

func setHeaders(c *fiber.Ctx, headerClaims map[string]pkg.SliceOrSingleValue[model.Claim], userInfos model.UserClaims) {
	if headerClaims == nil {
		headerClaims = config.DefaultForwardHeaders
	}
	for header, claim := range headerClaims {
		var value string
		var ok bool
		for _, cl := range claim {
			value, ok = userInfos.GetForHeader(cl)
			if ok {
				break
			}
		}
		if value != "" {
			c.Set(header, value)
		}
	}
}

func validateSession(sessionKey string) (claims model.UserClaims, err error) {
	var found bool
	found, err = cache.GetSession(sessionKey, &claims)
	if !found {
		claims = nil
	}
	return
}
