package artifact_builder

import (
	"github.com/chanzuckerberg/happy/pkg/config"
	log "github.com/sirupsen/logrus"
)

type artifactBuilderBuildOptions struct {
	slice *config.Slice

	tags []string
}

type ArtifactBuilderBuildOption func(*artifactBuilderBuildOptions)

// BuildSlice sets which slice we will build
func BuildSlice(slice *config.Slice) ArtifactBuilderBuildOption {
	return func(abbo *artifactBuilderBuildOptions) {
		abbo.slice = slice
	}
}

// WithTags sets the tags that will be added to your docker image
// Note that we will also set some default tags in addition
func WithTags(tags ...string) ArtifactBuilderBuildOption {
	log.Error(tags)
	return func(abbo *artifactBuilderBuildOptions) {
		abbo.tags = tags
	}
}
