package aws

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// NOTE(el): This is based off RFC3339 with some tweaks to make it a valid docker tag
const dockerRFC3339TimeFmt string = "2006-01-02T15-04-05"

// GetUserName will attempt to derive the caller's username
func (b *Backend) GetUserName(ctx context.Context) (string, error) {
	if b.username != nil {
		return *b.username, nil
	}

	var getter func(context.Context) (string, error)
	if util.IsCI(ctx) {
		getter = b.getUsernamefromGitHubActions
	} else {
		getter = b.getUsernameFromAWS
	}

	username, err := getter(ctx)
	username = cleanupUserName(username)

	b.username = &username
	return username, err
}

func (b *Backend) getUsernamefromGitHubActions(ctx context.Context) (string, error) {
	return os.Getenv("GITHUB_ACTOR"), nil
}

func (b *Backend) getUsernameFromAWS(ctx context.Context) (string, error) {
	out, err := b.stsclient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "could not get identity")
	}

	// Role sessions are emails, extract them
	userid := *out.UserId
	fragments := strings.Split(userid, ":")

	log.Debugf("Identity: %s", userid)

	if util.IsLocalstackMode() {
		if len(fragments) == 1 { // Localstack returns identity like AKIAIOSFODNN7EXAMPLE
			return fragments[0], nil
		} else {
			return "", errors.Errorf("unexpected local user identity %s", userid)
		}
	}

	if len(fragments) != 2 {
		return "", errors.Errorf("unexpected user identity %s", userid)
	}

	username := fragments[1]
	return username, nil
}

func (b *Backend) GetAccountID(ctx context.Context) (string, error) {
	if b.awsAccountID != nil {
		return *b.awsAccountID, nil
	}

	out, err := b.stsclient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "could not get identity")
	}

	b.awsAccountID = out.Account

	return *b.awsAccountID, nil
}

func (b *Backend) GenerateTag(ctx context.Context) (string, error) {
	username, err := b.GetUserName(ctx)
	if err != nil {
		return "", err
	}

	username = strings.ReplaceAll(username, "@", "-")

	t := time.Now().UTC().Format(dockerRFC3339TimeFmt)
	tag := fmt.Sprintf("%s-%s", username, t)

	return tag, nil
}
