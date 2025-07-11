package server

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-oidfed/lib"
	"github.com/go-oidfed/lib/apimodel"
	"github.com/go-oidfed/lib/oidfedconst"
	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v3/jws"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/offa/internal"
	"github.com/go-oidfed/offa/internal/cache"
	"github.com/go-oidfed/offa/internal/config"
	"github.com/go-oidfed/offa/internal/model"
	"github.com/go-oidfed/offa/internal/pkce"
)

const browserStateCookieName = "_offa_auth_state"

func addLoginHandlers(s fiber.Router) {
	path := config.Get().Server.Paths.Login
	s.Get(
		path, func(c *fiber.Ctx) error {
			opID := internal.FirstNonEmptyQueryParameter(c, "op", "entity_id", "entity", "iss", "issuer")
			if opID != "" {
				return doLogin(c, opID)
			}
			return showLoginPage(c)
		},
	)
	s.Get("/redirect", codeExchange)
}

var opOptions string

func scheduleBuildOPOptions() {
	ticker := time.NewTicker(time.Duration(config.Get().Federation.EntityCollectionInterval) * time.Minute) // Replace 5 with your desired interval

	go buildOPOptions()

	go func() {
		for range ticker.C {
			buildOPOptions()
		}
	}()
}

func buildOPOptions() {
	const opOptionFmt = `<option value="%s">%s</option>`
	var options string
	filters := []oidfed.EntityCollectionFilter{}
	allOPs := make(map[string]*oidfed.CollectedEntity)
	for _, ta := range config.Get().Federation.TrustAnchors {
		var collector oidfed.EntityCollector
		if config.Get().Federation.UseEntityCollectionEndpoint {
			collector = oidfed.SmartRemoteEntityCollector{TrustAnchors: config.Get().Federation.TrustAnchors.EntityIDs()}
		} else {
			collector = &oidfed.SimpleEntityCollector{}
		}
		ops := oidfed.FilterableVerifiedChainsEntityCollector{
			Collector: collector,
			Filters:   filters,
		}.CollectEntities(
			apimodel.EntityCollectionRequest{
				TrustAnchor: ta.EntityID,
				EntityTypes: []string{oidfedconst.EntityTypeOpenIDProvider},
			},
		)
		for _, op := range ops {
			allOPs[op.EntityID] = op
		}
	}
	for _, op := range allOPs {
		options += fmt.Sprintf(
			opOptionFmt, op.EntityID, getDisplayNameFromEntityInfo(op),
		)
	}
	opOptions = options
}

func getDisplayNameFromEntityInfo(entity *oidfed.CollectedEntity) string {
	if entity == nil {
		return ""
	}
	if entity.UIInfos == nil {
		return entity.EntityID
	}
	op, ok := entity.UIInfos[oidfedconst.EntityTypeOpenIDProvider]
	if ok && op.DisplayName != "" {
		return op.DisplayName
	}
	fed, ok := entity.UIInfos[oidfedconst.EntityTypeFederationEntity]
	if ok && fed.DisplayName != "" {
		return fed.DisplayName
	}
	return entity.EntityID
}

func showLoginPage(c *fiber.Ctx) error {
	var img string
	if config.Get().Federation.LogoURI != "" {
		img = fmt.Sprintf(`<img src="%s" alt="%s" class="logo"/>`, config.Get().Federation.LogoURI, "Logo")
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.SendString(fmt.Sprintf(loginHtml, config.Get().Federation.ClientName, img, opOptions, c.Query("next")))
}

type stateData struct {
	CodeChallenge pkce.PKCE
	Issuer        string
	BrowserState  string
	Next          string
}

func doLogin(c *fiber.Ctx, opID string) error {
	r, err := internal.RandomString(256)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}
	state := r[:64]
	browserState := r[64:128]
	pkceVerifier := r[128:192]
	nonce := r[192:224]

	pkceChallenge := pkce.NewS256PKCE(pkceVerifier)
	if err = cache.Set(
		cache.KeyStateData, state, stateData{
			CodeChallenge: *pkceChallenge,
			Issuer:        opID,
			BrowserState:  browserState,
			Next:          c.Query("next"),
		}, 5*time.Minute,
	); err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}
	challenge, err := pkceChallenge.Challenge()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}

	params := url.Values{}
	params.Set("nonce", nonce)
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", pkceChallenge.Method().String())
	params.Set("prompt", "consent")

	authURL, err := federationLeafEntity.GetAuthorizationURL(opID, redirectURI, state, scopes, params)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}
	c.Cookie(
		&fiber.Cookie{
			Name:     browserStateCookieName,
			Value:    browserState,
			Path:     getFullPath("/redirect"),
			MaxAge:   300,
			HTTPOnly: true,
			Secure:   config.Get().Server.Secure,
		},
	)
	return c.Redirect(authURL, fiber.StatusSeeOther)
}

func codeExchange(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	e := c.Query("error")
	errorDescription := c.Query("error_description")
	if e != "" {
		c.Status(444)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage(e, errorDescription))
	}
	var stateInfo stateData
	found, err := cache.Get(cache.KeyStateData, state, &stateInfo)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}
	if !found {
		c.Status(444)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("state mismatch", ""))
	}

	if stateInfo.BrowserState != c.Cookies(browserStateCookieName) {
		c.Status(444)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("state mismatch", ""))
	}

	params := url.Values{}
	params.Set("code_verifier", stateInfo.CodeChallenge.Verifier())
	log.WithField("code_verifier", stateInfo.CodeChallenge.Verifier()).Info("Code exchange with code verifier")

	tokenRes, errRes, err := federationLeafEntity.CodeExchange(stateInfo.Issuer, code, redirectURI, params)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}
	if errRes != nil {
		c.Status(444)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage(errRes.Error, errRes.ErrorDescription))
	}

	msg, err := jws.ParseString(tokenRes.IDToken)
	if err != nil {
		c.Status(444)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("error parsing id token", err.Error()))
	}
	c.ClearCookie(browserStateCookieName)
	if err = cache.Set(cache.KeyStateData, state, nil, time.Nanosecond); err != nil {
		log.WithError(err).Error("failed to clear state cache")
	}
	var idTokenData model.UserClaims
	err = json.Unmarshal(msg.Payload(), &idTokenData)
	if err != nil {
		c.Status(444)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("error decoding id token", err.Error()))
	}
	log.Debugf("Userclaims are: %+v", idTokenData)
	//TODO userinfo endpoint

	sessionID, err := internal.RandomString(128)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}
	if err = cache.SetSession(sessionID, idTokenData); err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(errorPage("internal server error", err.Error()))
	}

	c.Cookie(
		&fiber.Cookie{
			Name:     config.Get().SessionStorage.CookieName,
			Value:    sessionID,
			Domain:   config.Get().SessionStorage.CookieDomain,
			MaxAge:   config.Get().SessionStorage.TTL,
			HTTPOnly: true,
			Secure:   config.Get().Server.Secure,
			SameSite: "none",
		},
	)
	if stateInfo.Next == "" {
		stateInfo.Next = "/"
	}
	return c.Redirect(stateInfo.Next)
}
