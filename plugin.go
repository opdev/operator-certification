package plugin

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/opdev/operator-certification/internal/check"
	"github.com/opdev/operator-certification/internal/flags"
	"github.com/opdev/operator-certification/internal/image"
	"github.com/opdev/operator-certification/internal/pyxis"
	"golang.org/x/exp/slog"

	"github.com/Masterminds/semver/v3"
	plugin "github.com/opdev/knex/plugin/v0"
	"github.com/opdev/knex/types"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Assert that we implement the Plugin interface.
var _ plugin.Plugin = NewPlugin()

var vers = semver.MustParse("0.0.1")

var (
	DefaultCertImageFilename    = "cert-image.json"
	DefaultRPMManifestFilename  = "rpm-manifest.json"
	DefaultTestResultsFilename  = "results.json"
	DefaultArtifactsTarFileName = "artifacts.tar"
	DefaultPyxisHost            = "catalog.redhat.com/api/containers"
	DefaultPyxisEnv             = "prod"
	SystemdDir                  = "/etc/systemd/system"
)

func init() {
	plugin.Register("check-operator", NewPlugin())
}

type plug struct {
	image string
}

func NewPlugin() *plug {
	p := plug{}
	// plugin-related things may happen here.
	return &p
}

func (p *plug) Register() error {
	return nil
}

func (p *plug) Name() string {
	return "Operator Certification"
}

func (p *plug) Init(ctx context.Context, cfg *viper.Viper, args []string) error {
	slog.Info("Initializing Operator Certification")
	if len(args) != 1 {
		return errors.New("a single argument is required (the container image to test)")
	}

	return nil
}

func (p *plug) Flags() *pflag.FlagSet {
	return flags.FlagSet()
}

func (p *plug) Version() semver.Version {
	return *vers
}

func (p *plug) ExecuteChecks(ctx context.Context) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("Execute Called")
	pyxisClient := pyxis.NewPyxisClient(DefaultPyxisHost, "", "", &http.Client{Timeout: 60 * time.Second})
	certifiedImagesCheck := check.NewCertifiedImagesCheck(pyxisClient)

	// TODO: Get the image
	imageRef := image.Reference{
		ImageURI: p.image,
	}
	_, err := certifiedImagesCheck.Validate(ctx, imageRef)
	if err != nil {
		return err
	}
	return nil
}

func (p *plug) Results(ctx context.Context) types.Results {
	return types.Results{}
}

func (p *plug) Submit(ctx context.Context) error {
	slog.Info("Submit called")
	return nil
}
