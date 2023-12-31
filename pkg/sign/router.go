package sign

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign/transformers"
	"errors"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type RoadApp struct {
	Sites []*SiteConfig `json:"sites"`
}

func (v *RoadApp) Forward(ctx *fiber.Ctx, site *SiteConfig) error {
	if len(site.Upstreams) == 0 {
		return errors.New("invalid configuration")
	}

	// Boot processes
	for _, process := range site.Processes {
		if err := process.BootProcess(); err != nil {
			log.Warn().Err(err).Msgf("An error occurred when booting process (%s) for %s", process.ID, site.ID)
			return fiber.ErrBadGateway
		}
	}

	// Do forward
	idx := rand.Intn(len(site.Upstreams))
	upstream := site.Upstreams[idx]

	switch upstream.GetType() {
	case UpstreamTypeHypertext:
		return makeHypertextResponse(ctx, upstream)
	case UpstreamTypeFile:
		return makeFileResponse(ctx, upstream)
	default:
		return fiber.ErrBadGateway
	}
}

type RequestTransformerConfig = transformers.RequestTransformerConfig

type SiteConfig struct {
	ID           string                      `json:"id"`
	Rules        []*RouterRule               `json:"rules" yaml:"rules"`
	Transformers []*RequestTransformerConfig `json:"transformers" yaml:"transformers"`
	Upstreams    []*UpstreamInstance         `json:"upstreams" yaml:"upstreams"`
	Processes    []*ProcessInstance          `json:"processes" yaml:"processes"`
}

type RouterRule struct {
	Host    []string            `json:"host" yaml:"host"`
	Path    []string            `json:"path" yaml:"path"`
	Queries map[string]string   `json:"queries" yaml:"queries"`
	Headers map[string][]string `json:"headers" yaml:"headers"`
}
