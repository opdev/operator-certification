package flags

import "github.com/spf13/pflag"

const (
	KeyDockerConfig = "docker-config"
	KeyKubeconfig   = "kubeconfig"
	KeyIndexImage   = "index-image"
)

func dockerConfigFilePathFlag(fs *pflag.FlagSet) *string {
	return fs.StringP(
		KeyDockerConfig,
		"d",
		"",
		"Path to docker config.json file. This value is optional for publicly accessible images.\n"+
			"However, it is strongly encouraged for public Docker Hub images,\n"+
			"due to the rate limit imposed for unauthenticated requests. (env: PFLT_DOCKERCONFIG)",
	)
}

func kubeconfigFilePath(fs *pflag.FlagSet) *string {
	return fs.String(
		KeyKubeconfig,
		"",
		"Path to kubeconfig file.",
	)
}

func indexImageUri(fs *pflag.FlagSet) *string {
	return fs.String(
		KeyIndexImage,
		"",
		"Index image URI",
	)
}

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("operator-certification", pflag.ContinueOnError)

	dockerConfigFilePathFlag(fs)
	kubeconfigFilePath(fs)
	indexImageUri(fs)

	return fs
}
