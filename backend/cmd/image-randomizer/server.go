package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrCache "github.com/rl404/fairy/monitoring/newrelic/cache"
	"github.com/rl404/image-randomizer/internal/delivery/rest/api"
	"github.com/rl404/image-randomizer/internal/delivery/rest/ping"
	"github.com/rl404/image-randomizer/internal/delivery/rest/swagger"
	imageRepository "github.com/rl404/image-randomizer/internal/domain/image/repository"
	imageCache "github.com/rl404/image-randomizer/internal/domain/image/repository/cache"
	imageDB "github.com/rl404/image-randomizer/internal/domain/image/repository/db"
	tokenRepository "github.com/rl404/image-randomizer/internal/domain/token/repository"
	tokenCache "github.com/rl404/image-randomizer/internal/domain/token/repository/cache"
	userRepository "github.com/rl404/image-randomizer/internal/domain/user/repository"
	userCache "github.com/rl404/image-randomizer/internal/domain/user/repository/cache"
	userDB "github.com/rl404/image-randomizer/internal/domain/user/repository/db"
	"github.com/rl404/image-randomizer/internal/service"
	"github.com/rl404/image-randomizer/internal/utils"
	"github.com/rl404/image-randomizer/pkg/cache"
	"github.com/rl404/image-randomizer/pkg/http"
)

func server() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init newrelic.
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Newrelic.Name),
		newrelic.ConfigLicense(cfg.Newrelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		utils.Error(err.Error())
	} else {
		defer nrApp.Shutdown(10 * time.Second)
		utils.AddLog(_nr.NewFromNewrelicApp(nrApp, _nr.ErrorLevel))
		utils.Info("newrelic initialized")
	}

	// Init cache.
	c, err := cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
	if err != nil {
		return err
	}
	c = nrCache.New(cfg.Cache.Dialect, cfg.Cache.Address, c)
	utils.Info("cache initialized")
	defer c.Close()

	// Init in-memory.
	im, err := cache.New(cache.InMemory, "", "", time.Minute)
	if err != nil {
		return err
	}
	im = nrCache.New("inmemory", "inmemory", im)
	utils.Info("in-memory initialized")
	defer im.Close()

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	tmp, _ := db.DB()
	defer tmp.Close()

	// Init user.
	var user userRepository.Repository
	user = userDB.New(db)
	user = userCache.New(c, user)
	user = userCache.New(im, user)
	utils.Info("repository user initialized")

	// Init image.
	var image imageRepository.Repository
	image = imageDB.New(db)
	image = imageCache.New(c, image)
	image = imageCache.New(im, image)
	utils.Info("repository image initialized")

	// Init token.
	var token tokenRepository.Repository = tokenCache.New(c,
		cfg.JWT.AccessSecret,
		cfg.JWT.AccessExpired,
		cfg.JWT.RefreshSecret,
		cfg.JWT.RefreshExpired,
	)
	utils.Info("repository token initialized")

	// Init service.
	service := service.New(user, image, token)
	utils.Info("service initialized")

	// Init web server.
	httpServer := http.New(http.Config{
		Port:            cfg.App.Port,
		ReadTimeout:     cfg.App.ReadTimeout,
		WriteTimeout:    cfg.App.WriteTimeout,
		GracefulTimeout: cfg.App.GracefulTimeout,
	})
	utils.Info("http server initialized")

	r := httpServer.Router()
	r.Use(middleware.RealIP)
	r.Use(utils.Recoverer)
	utils.Info("http server middleware initialized")

	// Register ping route.
	ping.New().Register(r)
	utils.Info("http route ping initialized")

	// Register swagger route.
	swagger.New().Register(r)
	utils.Info("http route swagger initialized")

	// Register api route.
	api.New(service, cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret).Register(r, nrApp)
	utils.Info("http route api initialized")

	// Run web server.
	httpServerChan := httpServer.Run()
	utils.Info("http server listening at :%s", cfg.App.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-httpServerChan:
		if err != nil {
			return err
		}
	case <-sigChan:
	}

	return nil
}
