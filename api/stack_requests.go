package api

import "github.com/pkg/errors"

func (c *HappyClient) GetStacklist() ([]string, error) {
	body := NewAppMetadata(c.happyConfig.App(), c.happyConfig.GetEnv(), "")
	resp, err := c.Get("/v1/stacklistItems", body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed with")
	}

	result := WrappedAppStacksWithCount{}
	err = ParseResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	stacklist := []string{}
	for _, record := range result.Records {
		stacklist = append(stacklist, record.Stack)
	}

	return stacklist, nil
}

func (c *HappyClient) AddToStacklist(stackName string) error {
	body := NewAppMetadata(c.happyConfig.App(), c.happyConfig.GetEnv(), stackName)
	resp, err := c.Post("/v1/stacklistItems", body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	result := WrappedAppStack{}
	err = ParseResponse(resp, &result)
	if err != nil {
		return err
	}

	if result.Record == nil {
		return errors.New("tried to create a record that already existed")
	}

	return nil
}

func (c *HappyClient) DeleteFromStacklist(stackName string) error {
	body := NewAppMetadata(c.happyConfig.App(), c.happyConfig.GetEnv(), stackName)
	resp, err := c.Delete("/v1/stacklistItems", body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	result := WrappedAppStack{}
	err = ParseResponse(resp, &result)
	if err != nil {
		return err
	}

	if result.Record == nil {
		return errors.New("tried to delete a record that did not exist")
	}

	return nil
}
