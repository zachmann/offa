package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-oidfed/lib"
	"github.com/lestrrat-go/jwx/v3/jwa"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/offa/internal"
	"github.com/go-oidfed/offa/internal/cache"
	"github.com/go-oidfed/offa/internal/config"
	"github.com/go-oidfed/offa/internal/logger"
	"github.com/go-oidfed/offa/internal/server"
)

func main() {
	handleSignals()
	config.MustLoadConfig()
	logger.Init()
	cache.Init()
	internal.InitKeys(internal.FedSigningKeyName, internal.OIDCSigningKeyName)
	for _, c := range config.Get().Federation.TrustMarks {
		if err := c.Verify(
			config.Get().Federation.EntityID, "",
			oidfed.NewTrustMarkSigner(internal.GetKey(internal.FedSigningKeyName), jwa.ES512()),
		); err != nil {
			log.Fatal(err)
		}
	}
	if config.Get().Federation.UseResolveEndpoint {
		oidfed.DefaultMetadataResolver = oidfed.SmartRemoteMetadataResolver{}
	}
	server.Init()
	server.Start()
}

func handleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGUSR1)
	go func() {
		for {
			sig := <-signals
			switch sig {
			case syscall.SIGHUP:
				reload()
			case syscall.SIGUSR1:
				reloadLogFiles()
			}
		}
	}()
}

func reload() {
	log.Info("Reloading config")
	config.MustLoadConfig()
	if config.Get().Federation.UseResolveEndpoint {
		oidfed.DefaultMetadataResolver = oidfed.SmartRemoteMetadataResolver{}
	}
	logger.SetOutput()
	logger.MustUpdateAccessLogger()
}

func reloadLogFiles() {
	log.Debug("Reloading log files")
	logger.SetOutput()
	logger.MustUpdateAccessLogger()
}
