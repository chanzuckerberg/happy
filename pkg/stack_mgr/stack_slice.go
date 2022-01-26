package stack_mgr

import (
	"errors"
	"fmt"
	"strings"

	config "github.com/chanzuckerberg/happy-deploy/pkg/config"
	"github.com/chanzuckerberg/happy-deploy/pkg/util"
)

func BuildSlice(happyConfig config.HappyConfigIface, sliceName string, defaultTag string) (stackTags map[string]string, defaultTag string, err error) {
	stackTags = make(map[string]string)
	defaultTag = ""

	slice, ok := happyConfig.GetSlices()[sliceName]
	if !ok {
		validSlices := joinKeys(happyConfig.GetSlices(), ", ")
		return stackTags, defaultTag, errors.New(fmt.Sprintf("Slice %s is invalid - valid names: %s", sliceName, validSlices))
	}
	buildImages := slice.BuildImages
	sliceTag, err := util.GenerateTag(happyConfig)
	if err != nil {
		return stackTags, defaultTag, err
	}

	err = runPush(PushOptions{tag: sliceTag, images: buildImages})
	if err != nil {
		return stackTags, defaultTag, fmt.Errorf("Failed to push image: %s", err)
	}

	if len(defaultTag) == 0 {
		defaultTag = happyConfig.SliceDefaultTag
	}

	for _, sliceImg := range buildImages {
		stackTags[sliceImg] = sliceTag
	}

	return stackTags, defaultTag, nil
}

func joinKeys(m map[string]config.Slice, separator string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, separator)
}
