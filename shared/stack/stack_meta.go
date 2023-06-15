package stack

import (
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/util"
	log "github.com/sirupsen/logrus"
)

func (s *StackService) NewStackMeta(stackName string) *StackMeta {
	return &StackMeta{
		StackName: stackName,
	}
}

type StackMeta struct {
	StackName         string            `json:"stack_name"`
	Env               string            `json:"env"`
	ImageTag          string            `json:"image_tag"`
	ImageTags         map[string]string `json:"image_tags"`
	App               string            `json:"app"`
	Slice             string            `json:"slice"`
	CreatedAt         string            `json:"created"`
	UpdatedAt         string            `json:"updated"`
	Owner             string            `json:"owner"`
	Priority          int               `json:"priority"` //TODO: DEPRECATED
	Repo              string            `json:"repo"`
	GitSHA            string            `json:"git_sha"`
	GitBranch         string            `json:"git_branch"`
	HappyConfigSecret string            `json:"happy_config_secret"`
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

func StackMetaLastUpdatedCreated() StackMetaUpdater {
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

// TODO: DEPRECATED
func StackMetaPriority() StackMetaUpdater {
	return func(s *StackMeta) {
		// pick a random number between 1000 and 5000
		// TODO: new happy apps don't use this feature, phase it out
		s.Priority = rand.Intn(5000-1000) + 1000
	}
}

func StackMetaRepo(dir string) StackMetaUpdater {
	return func(s *StackMeta) {
		path, err := exec.LookPath("git")
		if err != nil {
			log.Error("git not found in path")
			return
		}
		cmd := exec.Command(path, "config", "--get", "remote.origin.url")
		cmd.Dir = dir
		var out strings.Builder
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Error(err, "error running %s", cmd.String())
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

func StackMetaGitBranch(dir string) StackMetaUpdater {
	return func(s *StackMeta) {
		path, err := exec.LookPath("git")
		if err != nil {
			log.Error("git not found in path")
			return
		}
		cmd := exec.Command(path, "branch", "--show-current")
		cmd.Dir = dir
		var out strings.Builder
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Error(err, "error running %s", cmd.String())
			return
		}
		s.GitBranch = strings.TrimSpace(out.String())
	}
}

func StackMetaGitHash(dir string) StackMetaUpdater {
	return func(s *StackMeta) {
		isClean, _, err := util.IsCleanGitTree(dir)
		if err != nil {
			log.Errorf("error checking if git tree in %s was clean: %s", dir, err)
			return
		}
		if !isClean {
			s.GitSHA = "dirty git tree (PLEASE COMMIT YOUR CHANGES)"
			return
		}
		path, err := exec.LookPath("git")
		if err != nil {
			log.Error("git not found in path")
			return
		}
		cmd := exec.Command(path, "rev-parse", "--short", "HEAD")
		cmd.Dir = dir
		var out strings.Builder
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Error(err, "error running %s", cmd.String())
			return
		}
		s.GitSHA = strings.TrimSpace(out.String())
	}
}

func StackHappyConfigVersion(config *config.HappyConfig) StackMetaUpdater {
	return func(s *StackMeta) {
		s.HappyConfigSecret = config.GetSecretId()
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
		StackMetaImageTags(tags), // TODO: change the name of this, its confusing
		StackMetaLastUpdatedCreated(),
		StackMetaOwner(owner),
		StackMetaSliceName(slice),
		StackMetaPriority(), //TODO: DEPRECATED
		StackMetaRepo(projectRoot),
		StackMetaGitBranch(projectRoot),
		StackMetaGitHash(projectRoot),
		StackMetaAppName(config),
		StackMetaStackName(stackName),
		StackMetaEnv(env),
		StackHappyConfigVersion(config),
	)
	return s
}
