package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/qnfnypen/mumori/http/routers"
	_ "github.com/qnfnypen/mumori/public"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	engine := routers.GenerateEngine()

	server := &http.Server{
		Handler:        engine,
		Addr:           viper.GetString("HTTP.Addr"),
		ReadTimeout:    viper.GetDuration("HTTP.ReadTimeout") * time.Second,
		WriteTimeout:   viper.GetDuration("HTTP.WriteTimeout") * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		<-quit

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal().Str("error", err.Error()).Msg("关闭服务器")
		}
	}()

	log.Info().Msg("启动服务器")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Info().Msg("服务器响应请求而关闭")
		} else {
			log.Warn().Str("error", err.Error()).Msg("服务器由于意外而关闭")
		}
	}

}
