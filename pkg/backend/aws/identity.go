package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

// GetUserName will attempt to derive the caller's username
func (b *awsBackend) GetUserName(ctx context.Context) (string, error) {
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
	return fragments[1], nil
}
