package stack_mgr

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type StackMeta struct {
	StackName string
	Priority  int

	DataMap  map[string]string
	TagMap   map[string]string
	ParamMap map[string]string
}

// Update the image tag with the given newTag, and set the priority randomly.
// To not collide, setting priority requires knowing the priority of all other
// stacks from the StackMgr
func (s *StackMeta) Update(
	ctx context.Context,
	newTag string,
	stackTags map[string]string,
	sliceName string,
	stackSvc *StackService,
) error {
	s.DataMap["imagetag"] = newTag

	var imageTagsJson []byte
	imageTagsJson, err := json.Marshal(stackTags)
	if err != nil {
		return errors.Wrap(err, "unable to convert image tags to json")
	}
	s.DataMap["imagetags"] = string(imageTagsJson)

	if sliceName != "" {
		s.DataMap["slice"] = sliceName
	}

	now := time.Now().Format(time.RFC3339)
	if createdAt, ok := s.DataMap["created"]; !ok || createdAt == "" {
		s.DataMap["created"] = now
	}
	s.DataMap["updated"] = now

	ownerVal, ok := s.DataMap["owner"]
	if !ok || ownerVal == "" || ownerVal == unknown {
		username, err := stackSvc.backend.GetUserName(ctx)
		if err != nil {
			return err
		}
		s.DataMap["owner"] = username
	}

	err = s.setPriority(ctx, stackSvc)
	return errors.Wrap(err, "failed to update")
}

func (s *StackMeta) setPriority(ctx context.Context, stackMgr *StackService) error {
	if s.Priority != 0 {
		return nil
	}

	existingPrioirty := map[int]bool{}
	stacks, err := stackMgr.GetStacks(ctx)
	if err != nil {
		return err
	}

	for _, stack := range stacks {
		stackMeta, err := stack.Meta(ctx)
		if err != nil {
			if errors.Is(err, tfe.ErrResourceNotFound) {
				// Some stacks may not have a meta, so we ignore those errors. These stacks could be in dry run mode and will be cleaned up.
				log.Warnf("failed to get stack meta for %s: %s, skipping.", stack.Name, err.Error())
				continue
			}
			return errors.Wrap(err, "failed to retrieve stack meta")
		}
		existingPrioirty[stackMeta.Priority] = true
	}
	// pick a random number between 1000 and 5000 that's not in use right now.
	priority := 0
	for priority == 0 {
		randInt := random.Random(1000, 5000)
		if _, ok := existingPrioirty[randInt]; !ok {
			priority = randInt
		}
	}
	s.DataMap["priority"] = strconv.Itoa(priority)
	log.WithField("priority", priority).Debug("Priority for workspace set")
	return nil
}

func (s *StackMeta) GetTags() map[string]string {
	out := map[string]string{}
	for k, v := range s.TagMap {
		out[v] = s.DataMap[k]
	}
	return out
}

func (s *StackMeta) GetParameters() map[string]string {
	out := map[string]string{}
	for k, v := range s.ParamMap {
		out[v] = s.DataMap[k]
	}
	return out
}

func (s *StackMeta) Load(existingTags map[string]string) {
	for shortTag, tagName := range s.TagMap {
		if _, ok := existingTags[tagName]; ok {
			s.DataMap[shortTag] = existingTags[tagName]
		} else if _, ok := s.DataMap[shortTag]; !ok {
			s.DataMap[shortTag] = ""
		}
	}
}
