module github.com/chanzuckerberg/happy

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.68
	github.com/blang/semver v3.5.1+incompatible
	github.com/chanzuckerberg/go-misc v0.0.0-20210623155112-450db80199ff
	github.com/golang/mock v1.5.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gruntwork-io/terratest v0.34.2
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.8 // indirect
	github.com/hashicorp/go-slug v0.7.0 // indirect
	github.com/hashicorp/go-tfe v0.14.0
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// replace github.com/chanzuckerberg/go-misc => ../go-misc
