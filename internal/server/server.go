package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v3/jwa"
	log "github.com/sirupsen/logrus"
	"github.com/zachmann/go-oidfed/pkg"

	"github.com/zachmann/offa/internal"
	"github.com/zachmann/offa/internal/config"
)

var server *fiber.App

var serverConfig = fiber.Config{
	ReadTimeout:    3 * time.Second,
	WriteTimeout:   3 * time.Second,
	IdleTimeout:    150 * time.Second,
	ReadBufferSize: 8192,
	// WriteBufferSize: 4096,
	ErrorHandler: handleError,
	Network:      "tcp",
}

var federationLeafEntity *pkg.FederationLeaf
var requestObjectProducer *pkg.RequestObjectProducer
var scopes string
var redirectURI string
var fullLoginPath string

// Init initializes the server
func Init() {
	initHtmls()
	initFederationEntity()
	server = fiber.New(serverConfig)
	addMiddlewares(server)
	addFederationEndpoints(server)
	addAuthHandlers(server)
	addLoginHandlers(server)
}

func initFederationEntity() {
	fedConfig := config.Get().Federation
	if fedConfig.EntityID[len(fedConfig.EntityID)-1] == '/' {
		redirectURI = fedConfig.EntityID + "redirect"
	} else {
		redirectURI = fedConfig.EntityID + "/redirect"
	}
	fullLoginPath = fedConfig.EntityID + getFullPath(config.Get().Server.Paths.Login)
	scopes = strings.Join(fedConfig.Scopes, " ")
	if scopes == "" {
		scopes = "openid profile email"
	}
	requestObjectProducer = pkg.NewRequestObjectProducer(
		fedConfig.EntityID, internal.GetKey(internal.OIDCSigningKeyName), jwa.ES512(), 60,
	)

	metadata := &pkg.Metadata{
		RelyingParty: &pkg.OpenIDRelyingPartyMetadata{
			Scope:                   scopes,
			RedirectURIS:            []string{redirectURI},
			ResponseTypes:           []string{"code"},
			GrantTypes:              []string{"authorization_code"},
			ApplicationType:         "web",
			ClientName:              fedConfig.ClientName,
			LogoURI:                 fedConfig.LogoURI,
			JWKS:                    internal.GetJWKS(internal.OIDCSigningKeyName),
			OrganizationName:        fedConfig.OrganisationName,
			ClientRegistrationTypes: []string{"automatic"},
		},
		FederationEntity: &pkg.FederationEntityMetadata{
			OrganizationName: fedConfig.OrganisationName,
			LogoURI:          fedConfig.LogoURI,
		},
	}
	var err error
	federationLeafEntity, err = pkg.NewFederationLeaf(
		fedConfig.EntityID, fedConfig.AuthorityHints, fedConfig.TrustAnchors, metadata,
		pkg.NewEntityStatementSigner(
			internal.GetKey(internal.FedSigningKeyName),
			jwa.ES512(),
		), 86400, internal.GetKey(internal.OIDCSigningKeyName), jwa.ES512(),
	)
	if err != nil {
		log.Fatal(err)
	}
	federationLeafEntity.TrustMarks = fedConfig.TrustMarks
}

func start(s *fiber.App) {
	if !config.Get().Server.TLS.Enabled {
		log.WithField("port", config.Get().Server.Port).Info("TLS is disabled starting http server")
		log.WithError(s.Listen(fmt.Sprintf(":%d", config.Get().Server.Port))).Fatal()
	}
	// TLS enabled
	if config.Get().Server.TLS.RedirectHTTP {
		httpServer := fiber.New(serverConfig)
		httpServer.All(
			"*", func(ctx *fiber.Ctx) error {
				//goland:noinspection HttpUrlsUsage
				return ctx.Redirect(
					strings.Replace(ctx.Request().URI().String(), "http://", "https://", 1),
					fiber.StatusPermanentRedirect,
				)
			},
		)
		log.Info("TLS and http redirect enabled, starting redirect server on port 80")
		go func() {
			log.WithError(httpServer.Listen(":80")).Fatal()
		}()
	}
	time.Sleep(time.Millisecond) // This is just for a more pretty output with the tls header printed after the http one
	log.Info("TLS enabled, starting https server on port 443")
	log.WithError(s.ListenTLS(":443", config.Get().Server.TLS.Cert, config.Get().Server.TLS.Key)).Fatal()
}

// Start starts the server
func Start() {
	start(server)
}

func getFullPath(path string) string {
	if len(path) == 0 {
		return config.Get().Server.Basepath
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return config.Get().Server.Basepath + path
}
