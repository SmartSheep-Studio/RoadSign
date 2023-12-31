package hypertext

import (
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "RoadSign",
		ServerHeader:          "RoadSign",
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		Prefork:               viper.GetBool("performance.prefork"),
		BodyLimit:             viper.GetInt("hypertext.limitation.max_body_size"),
	})

	if viper.GetBool("performance.request_logging") {
		app.Use(logger.New(logger.Config{
			Output: log.Logger,
			Format: "[Proxies] [${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	if viper.GetInt("hypertext.limitation.max_qps") > 0 {
		app.Use(limiter.New(limiter.Config{
			Max:        viper.GetInt("hypertext.limitation.max_qps"),
			Expiration: 1 * time.Second,
		}))
	}

	UseProxies(app)

	return app
}

func RunServer(app *fiber.App, ports []string, securedPorts []string, pem string, key string) {
	for _, port := range ports {
		port := port
		go func() {
			if viper.GetBool("hypertext.certificate.redirect") {
				redirector := fiber.New(fiber.Config{
					AppName:               "RoadSign",
					ServerHeader:          "RoadSign",
					DisableStartupMessage: true,
					EnableIPValidation:    true,
				})
				redirector.All("/", func(c *fiber.Ctx) error {
					return c.Redirect(strings.ReplaceAll(string(c.Request().URI().FullURI()), "http", "https"))
				})
				if err := redirector.Listen(port); err != nil {
					log.Panic().Err(err).Msg("An error occurred when listening hypertext common ports.")
				}
			} else {
				if err := app.Listen(port); err != nil {
					log.Panic().Err(err).Msg("An error occurred when listening hypertext common ports.")
				}
			}
		}()
	}

	for _, port := range securedPorts {
		port := port
		go func() {
			if err := app.ListenTLS(port, pem, key); err != nil {
				log.Panic().Err(err).Msg("An error occurred when listening hypertext tls ports.")
			}
		}()
	}
}
