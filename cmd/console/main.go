package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	dist string
	log  *zap.Logger

	port int
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("DIST", "/app")
	viper.SetDefault("LOG_LEVEL", 0)
	viper.SetDefault("PORT", 80)

	port = viper.GetInt("PORT")
	level := viper.GetInt("LOG_LEVEL")

	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.Level(level))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	log = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	dist = viper.GetString("DIST")
}

var ETAG = uuid.New().String()

func StaticHandler(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))

	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("If-Modified-Since")
		w.Header().Set("Etag", ETAG)

		fs.ServeHTTP(w, r)
	}
}

func main() {
	defer func() {
		_ = log.Sync()
	}()
	log.Debug("Debug enabled")

	mux := http.NewServeMux()
	mux.Handle("/", StaticHandler(dist))

	log.Info("Starting HTTP Server", zap.Int("port", port), zap.String("etag", ETAG))
	log.Fatal("Failed to serve", zap.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour),
		Debug:            false,
	}).Handler(mux))))
}
