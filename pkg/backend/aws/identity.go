package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

// NOTE(el): This is based off RFC3339 with some tweaks to make it a valid docker tag
const dockerRFC3339TimeFmt string = "2006-01-02T15-04-05"

// GetUserName will attempt to derive the caller's username
func (b *Backend) GetUserName(ctx context.Context) (string, error) {
	if b.username != nil {
		return *b.username, nil
	}

	out, err := b.stsclient.GetCallerIdentityWithContext(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "could not get identity")
	}

	// Role sessions are emails, extract them
	userid := *out.UserId
	fragments := strings.Split(userid, ":")
	if len(fragments) != 2 {
		return "", errors.Errorf("unexpected user identity %s", userid)
	}

	username := fragments[1]
	b.username = &username

	return username, nil
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
