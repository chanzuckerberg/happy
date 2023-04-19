package stack_mgr

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *StackService) NewStackMeta(stackName string) *StackMeta {
	return &StackMeta{
		StackName: stackName,
	}
}

type StackMeta struct {
	StackName string            `json:"stack_name"`
	Env       string            `json:"env"`
	ImageTag  string            `json:"image_tag"`
	ImageTags map[string]string `json:"image_tags"`
	App       string            `json:"app"`
	Slice     string            `json:"slice"`
	CreatedAt string            `json:"created"`
	UpdatedAt string            `json:"updated"`
	Owner     string            `json:"owner"`
	Priority  int               `json:"priority"`
	Repo      string            `json:"repo"`
	GitSHA    string            `json:"git_sha"`
	GitBranch string            `json:"git_branch"`
}

func (s *StackMeta) Merge(legacy StackMetaLegacy) error {
	if legacy.StackName != "" {
		s.StackName = legacy.StackName
	}
	if legacy.Env != "" {
		s.Env = legacy.Env
	}
	if legacy.ImageTag != "" {
		s.ImageTag = legacy.ImageTag
	}
	if legacy.ImageTags != "" {
		err := json.Unmarshal([]byte(legacy.ImageTags), &s.ImageTags)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal image tags from legacy stack meta")
		}
	}
	if legacy.App != "" {
		s.App = legacy.App
	}
	if legacy.Slice != "" {
		s.Slice = legacy.Slice
	}
	if legacy.CreatedAt != "" {
		s.CreatedAt = legacy.CreatedAt
	}
	if legacy.UpdatedAt != "" {
		s.UpdatedAt = legacy.UpdatedAt
	}
	if legacy.Owner != "" {
		s.Owner = legacy.Owner
	}
	if legacy.Priority != "" {
		var err error
		s.Priority, err = strconv.Atoi(legacy.Priority)
		if err != nil {
			return errors.Wrap(err, "legacy priority is not a valid int")
		}
	}

	return nil
}

// TODO: remove this in a few weeks when folks update their stacks
type StackMetaLegacy struct {
	ParamMap  map[string]string
	StackName string `json:"happy/instance"`
	Env       string `json:"happy/env"`
	ImageTag  string `json:"happy/meta/imagetag"`
	ImageTags string `json:"happy/meta/imagetags"`
	App       string `json:"happy/app"`
	Slice     string `json:"happy/meta/slice"`
	CreatedAt string `json:"happy/meta/created-at"`
	UpdatedAt string `json:"happy/meta/updated-at"`
	Owner     string `json:"happy/meta/owner"`
	Priority  string `json:"happy/meta/priority"`
}

type StackMetaUpdater func(s *StackMeta)

func StackMetaImageTag(tag string) StackMetaUpdater {
	return func(s *StackMeta) {
		s.ImageTag = tag
	}
}

func StackMetaImageTags(tags map[string]string) StackMetaUpdater {
	return func(s *StackMeta) {
		s.ImageTags = tags
	}
}

func StackMetaSliceName(sliceName string) StackMetaUpdater {
	return func(s *StackMeta) {
		s.Slice = sliceName
	}
}

func StackMetaLastUpdated() StackMetaUpdater {
	return func(s *StackMeta) {
		now := time.Now().Format(time.RFC3339)
		if s.CreatedAt == "" {
			s.CreatedAt = now
		}
		s.UpdatedAt = now
	}
}

func StackMetaOwner(owner string) StackMetaUpdater {
	return func(s *StackMeta) {
		s.Owner = owner
	}
}

func StackMetaPriority() StackMetaUpdater {
	return func(s *StackMeta) {
		// pick a random number between 1000 and 5000
		// TODO: new happy apps don't use this feature, phase it out
		s.Priority = random.Random(1000, 5000)
	}
}

func StackMetaRepo(dir string) StackMetaUpdater {
	return func(s *StackMeta) {
		cmd := exec.Command("git", "config", "--get", "remote.origin.url")
		cmd.Dir = dir
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			logrus.Error(err, "error running %s", cmd.String())
			return
		}
		s.Repo = strings.TrimSpace(out.String())
	}
}

func StackMetaAppName(config *config.HappyConfig) StackMetaUpdater {
	return func(s *StackMeta) {
		s.App = config.App()
	}
}

func StackMetaStackName(stackName string) StackMetaUpdater {
	return func(s *StackMeta) {
		s.StackName = stackName
	}
}

func StackMetaEnv(env string) StackMetaUpdater {
	return func(s *StackMeta) {
		s.Env = env
	}
}

func StackMetaGit(dir string) StackMetaUpdater {
	return func(s *StackMeta) {
		cmd := exec.Command("git", "branch", "--show-current")
		cmd.Dir = dir
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			logrus.Error(err, "error running %s", cmd.String())
			return
		}
		s.Repo = strings.TrimSpace(out.String())

		isClean, _, err := util.IsCleanGitTree(dir)
		if err != nil {
			logrus.Errorf(err, "error checking if git tree in %s was clean", dir)
			return
		}
		if !isClean {
			return
		}
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		cmd.Dir = dir
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			logrus.Error(err, "error running %s", cmd.String())
			return
		}
		s.Repo = strings.TrimSpace(out.String())
		s.Repo = strings.TrimSpace(out.String())
	}
}

func (s *StackMeta) update(updaters ...StackMetaUpdater) {
	for _, updater := range updaters {
		updater(s)
	}
}

func (s *StackMeta) UpdateAll(
	tag string,
	tags map[string]string,
	slice string,
	owner string,
	projectRoot string,
	config *config.HappyConfig,
	stackName string,
	env string,
) *StackMeta {
	s.update(
		StackMetaImageTag(tag),
		StackMetaImageTags(tags),
		StackMetaLastUpdated(),
		StackMetaOwner(owner),
		StackMetaSliceName(slice),
		StackMetaPriority(), //TODO: phase this out
		StackMetaRepo(projectRoot),
		StackMetaAppName(config),
		StackMetaStackName(stackName),
		StackMetaEnv(env),
	)
	return s
}
