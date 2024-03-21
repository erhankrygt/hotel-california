package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"hotel-california-backend/internal/localization"
	mysqlstore "hotel-california-backend/internal/store/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"hotel-california-backend/configs/envvars"
	"hotel-california-backend/internal/middleware"
	"hotel-california-backend/internal/service"

	hotelcalifornia "hotel-california-backend"
	authmiddlaware "hotel-california-backend/internal/middleware/auth"
	httptransport "hotel-california-backend/internal/transport/http"
)

func main() {
	var l log.Logger
	{
		l = log.NewLogfmtLogger(os.Stdout)
		l = log.With(l, "time", log.DefaultTimestampUTC)
	}

	var ev *envvars.EnvVars
	var err error
	{
		ev, err = envvars.LoadEnvVars()
		if err != nil {
			_ = l.Log("error", err.Error())
			return
		}
	}

	var ps mysqlstore.Store
	{
		ps, err = mysqlstore.NewStore(mysqlstore.Options{
			UserName:       ev.MySql.UserName,
			Password:       ev.MySql.Password,
			URI:            ev.MySql.URI,
			Database:       ev.MySql.Database,
			Port:           ev.MySql.Port,
			ConnectTimeout: ev.MySql.ConnectTimeout,
		})

		if err != nil {
			_ = l.Log("error", err.Error())
			return
		}
	}

	var s hotelcalifornia.Service
	{
		s = service.NewService(ev.Service.Environment, l, ps, ev.JWTToken)
	}

	var am middleware.Middleware
	{
		am = authmiddlaware.NewAuthMiddleware(l, ev.JWTToken)

		s = am(s)
	}

	var h http.Handler
	{
		h = httptransport.MakeHTTPHandler(log.With(l, "transport", "http"), s)
	}

	var hs *http.Server
	{
		hs = &http.Server{
			Addr:           ev.HTTPServer.Address,
			ReadTimeout:    ev.HTTPServer.ReadTimeout,
			WriteTimeout:   ev.HTTPServer.WriteTimeout,
			IdleTimeout:    ev.HTTPServer.IdleTimeout,
			MaxHeaderBytes: ev.HTTPServer.MaxHeaderBytes,
			Handler:        h,
		}
	}

	err = localization.InitializeBundle(ev.Localization)
	if err != nil {
		_ = l.Log("error", err.Error())
		return
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- errors.New((<-c).String())
	}()

	go func() {
		_ = l.Log("transport", "http", "address", ev.HTTPServer.Address)

		err = hs.ListenAndServe()

		if !errors.Is(http.ErrServerClosed, err) {
			errs <- err
		}
	}()

	err = <-errs
	_ = l.Log("error", err.Error())

	ctx, cf := context.WithTimeout(context.Background(), ev.HTTPServer.ShutdownTimeout)

	defer cf()

	if err = hs.Shutdown(ctx); err != nil {
		_ = l.Log("error", err.Error())
	}

	if err = ps.Close(); err != nil {
		_ = l.Log("error", err.Error())
	}
}
