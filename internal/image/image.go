package image

import v1 "github.com/google/go-containerregistry/pkg/v1"

type Reference struct {
	ImageURI    string
	ImageFSPath string
	ImageInfo   v1.Image
}
