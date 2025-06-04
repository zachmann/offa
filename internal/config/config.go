package config

import (
	"net"
	"net/url"
	"os"
	"regexp"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/zachmann/go-oidfed/pkg"
	"gopkg.in/yaml.v3"

	"github.com/zachmann/offa/internal/model"
)

var conf *Config

// Get returns the Config
func Get() *Config {
	return conf
}

// Config holds the configuration for this application
type Config struct {
	Server         serverConf     `yaml:"server"`
	Logging        loggingConf    `yaml:"logging"`
	Federation     federationConf `yaml:"federation"`
	Auth           authConf       `yaml:"auth"`
	SessionStorage sessionConf    `yaml:"sessions"`
	DebugAuth      bool           `yaml:"debug_auth"`
}

type federationConf struct {
	EntityID                    string                                    `yaml:"entity_id"`
	ClientName                  string                                    `yaml:"client_name"`
	LogoURI                     string                                    `yaml:"logo_uri"`
	Scopes                      []string                                  `yaml:"scopes"`
	TrustAnchors                pkg.TrustAnchors                          `yaml:"trust_anchors"`
	AuthorityHints              []string                                  `yaml:"authority_hints"`
	OrganizationName            string                                    `yaml:"organization_name"`
	KeyStorage                  string                                    `yaml:"key_storage"`
	OnlyAutomaticOPs            bool                                      `yaml:"filter_to_automatic_ops"`
	TrustMarks                  []*pkg.EntityConfigurationTrustMarkConfig `yaml:"trust_marks"`
	UseResolveEndpoint          bool                                      `yaml:"use_resolve_endpoint"`
	UseEntityCollectionEndpoint bool                                      `yaml:"use_entity_collection_endpoint"`
}

type sessionConf struct {
	TTL             int                                            `yaml:"ttl"`
	RedisAddr       string                                         `yaml:"redis_addr"`
	MemCachedAddr   string                                         `yaml:"memcached_addr"`
	MemCachedClaims map[string]pkg.SliceOrSingleValue[model.Claim] `yaml:"memcached_claims"`
	CookieName      string                                         `yaml:"cookie_name"`
	CookieDomain    string                                         `yaml:"cookie_domain"`
}

func (c sessionConf) validate() error {
	if c.MemCachedClaims != nil {
		if _, set := c.MemCachedClaims["UserName"]; !set {
			return errors.New("sessions.memcached_claims is set, but no claim for 'UserName' is set")
		}
		if _, set := c.MemCachedClaims["Groups"]; !set {
			return errors.New("sessions.memcached_claims is set, but no claim for 'Groups' is set")
		}
	}
	return nil
}

type authConf []*authRule

type authRule struct {
	Domain             string                                                                 `yaml:"domain"`
	DomainRegex        string                                                                 `yaml:"domain_regex"`
	DomainPattern      *regexp.Regexp                                                         `yaml:"-"`
	Path               string                                                                 `yaml:"path"`
	PathRegex          string                                                                 `yaml:"path_regex"`
	PathPattern        *regexp.Regexp                                                         `yaml:"-"`
	Require            pkg.SliceOrSingleValue[map[model.Claim]pkg.SliceOrSingleValue[string]] `yaml:"require"`
	ForwardHeaders     map[string]pkg.SliceOrSingleValue[model.Claim]                         `yaml:"forward_headers"`
	RedirectStatusCode int                                                                    `yaml:"redirect_status"`
}

var DefaultForwardHeaders = map[string]pkg.SliceOrSingleValue[model.Claim]{
	"X-Forwarded-User": {
		"preferred_username",
		"sub",
	},
	"X-Forwarded-Email":    {"email"},
	"X-Forwarded-Provider": {"iss"},
	"X-Forwarded-Subject":  {"sub"},
	"X-Forwarded-Groups": {
		"entitlements",
		"groups",
	},
	"X-Forwarded-Name": {"name"},
}
var DefaultMemCachedClaims = map[string]pkg.SliceOrSingleValue[model.Claim]{
	"UserName": {
		"preferred_username",
		"sub",
	},
	"Groups":    {"groups"},
	"Email":     {"email"},
	"Name":      {"name"},
	"GivenName": {"given_name"},
	"Provider":  {"iss"},
	"Subject":   {"sub"},
}

func (r *authRule) validate() error {
	if r.Domain != "" {
		r.DomainRegex = regexp.QuoteMeta(r.Domain)
	}
	if r.Path != "" {
		r.PathRegex = regexp.QuoteMeta(r.Path)
	}
	if r.Domain == "" {
		return errors.New("domain or domain_regex is required")
	}
	r.DomainPattern = regexp.MustCompile(r.DomainRegex)
	if r.PathRegex != "" {
		r.PathPattern = regexp.MustCompile(r.PathRegex)
	}
	return nil
}

func (c *authConf) validate() error {
	for i, rule := range *c {
		if err := rule.validate(); err != nil {
			return err
		}
		(*c)[i] = rule
	}
	return nil
}

func (c authConf) FindRule(host, path string) *authRule {
	for _, rule := range c {
		if rule.DomainPattern.MatchString(host) {
			if rule.PathPattern == nil {
				return rule
			}
			if rule.PathPattern.MatchString(path) {
				return rule
			}
		}
	}
	return nil
}

type serverConf struct {
	Port           int          `yaml:"port"`
	TLS            tlsConf      `yaml:"tls"`
	TrustedProxies []string     `yaml:"trusted_proxies"`
	TrustedNets    []*net.IPNet `yaml:"-"`
	Paths          pathConf     `yaml:"paths"`
	Secure         bool         `yaml:"-"`
	Basepath       string       `yaml:"-"`
}

type pathConf struct {
	Login       string `yaml:"login"`
	ForwardAuth string `yaml:"forward_auth"`
}

type tlsConf struct {
	Enabled      bool   `yaml:"enabled"`
	RedirectHTTP bool   `yaml:"redirect_http"`
	Cert         string `yaml:"cert"`
	Key          string `yaml:"key"`
}

type loggingConf struct {
	Access   LoggerConf         `yaml:"access"`
	Internal internalLoggerConf `yaml:"internal"`
}

type internalLoggerConf struct {
	LoggerConf `yaml:",inline"`
	Smart      smartLoggerConf `yaml:"smart"`
}

// LoggerConf holds configuration related to logging
type LoggerConf struct {
	Dir    string `yaml:"dir"`
	StdErr bool   `yaml:"stderr"`
	Level  string `yaml:"level"`
}

type smartLoggerConf struct {
	Enabled bool   `yaml:"enabled"`
	Dir     string `yaml:"dir"`
}

func checkLoggingDirExists(dir string) error {
	if dir != "" && !fileExists(dir) {
		return errors.Errorf("logging directory '%s' does not exist", dir)
	}
	return nil
}

func (c *serverConf) validate() error {
	for _, cidr := range c.TrustedProxies {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return errors.Wrapf(err, "invalid trusted proxy CIDR '%s'", cidr)
		}
		c.TrustedNets = append(c.TrustedNets, ipnet)
	}
	return nil
}

func (log *loggingConf) validate() error {
	if err := checkLoggingDirExists(log.Access.Dir); err != nil {
		return err
	}
	if err := checkLoggingDirExists(log.Internal.Dir); err != nil {
		return err
	}
	if log.Internal.Smart.Enabled {
		if log.Internal.Smart.Dir == "" {
			log.Internal.Smart.Dir = log.Internal.Dir
		}
		if err := checkLoggingDirExists(log.Internal.Smart.Dir); err != nil {
			return err
		}
	}
	return nil
}

var possibleConfigLocations = []string{
	".",
	"config",
	"/config",
	"/offa/config",
	"/offa",
	"/data/config",
	"/data",
	"/etc/offa",
}

func validate() error {
	if conf == nil {
		return errors.New("config not set")
	}
	if err := conf.Logging.validate(); err != nil {
		return err
	}
	if err := conf.Server.validate(); err != nil {
		return err
	}
	if err := conf.Auth.validate(); err != nil {
		return err
	}
	if err := conf.SessionStorage.validate(); err != nil {
		return err
	}
	u, err := url.Parse(conf.Federation.EntityID)
	if err != nil {
		return err
	}
	conf.Server.Secure = u.Scheme == "https"
	conf.Server.Basepath = u.Path
	if conf.Server.Basepath != "" {
		if conf.Server.Basepath[len(conf.Server.Basepath)-1] == '/' {
			conf.Server.Basepath = conf.Server.Basepath[:len(conf.Server.Basepath)-2]
		}
		if conf.Server.Basepath[0] != '/' {
			conf.Server.Basepath = "/" + conf.Server.Basepath
		}
	}
	return nil
}

func MustLoadConfig() {
	data, _ := mustReadConfigFile("config.yaml", possibleConfigLocations)
	conf = &Config{
		Server: serverConf{
			Port: 15661,
			Paths: pathConf{
				Login:       "/login",
				ForwardAuth: "/auth",
			},
		},
		SessionStorage: sessionConf{
			TTL:        3600,
			CookieName: "offa-session",
		},
		Federation: federationConf{ClientName: "OFFA - Openid Federation Forward Auth"},
	}
	if err := yaml.Unmarshal(data, conf); err != nil {
		log.Fatal(err)
	}
	if conf.Federation.KeyStorage == "" {
		log.Fatal("key_storage must be given")
	}
	if conf.Federation.LogoURI == "" {
		conf.Federation.LogoURI = conf.Federation.EntityID + "/static/img/offa-text.svg"
	}
	d, err := os.Stat(conf.Federation.KeyStorage)
	if err != nil {
		log.Fatal(err)
	}
	if !d.IsDir() {
		log.Fatalf("key_storage '%s' must be a directory", conf.Federation.KeyStorage)
	}
	if err = validate(); err != nil {
		log.Fatalf("%s", err)
	}
}
