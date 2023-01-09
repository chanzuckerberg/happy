package client

import (
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
)

// type HappyStacklistAPI interface {
// 	GetStacklist(appName, environment string) ([]string, error)
// 	AddToStacklist(appName, environment, stackName string) error
// 	DeleteFromStacklist(appName, environment, stackName string) error
// }

func (c *HappyClient) GetStacklist(appName, environment, taskLaunchType string) ([]string, error) {
	body := model.MakeAppStackPayload(appName, environment, "", taskLaunchType)
	result := model.WrappedAppStacksWithCount{}
	err := c.GetParsed("/v1/stacklistItems", body, &result)
	if err != nil {
		return nil, err
	}

	stacklist := []string{}
	for _, record := range result.Records {
		stacklist = append(stacklist, record.Stack)
	}

	return stacklist, nil
}

func (c *HappyClient) AddToStacklist(appName, environment, stackName, taskLaunchType string) error {
	body := model.MakeAppStackPayload(appName, environment, stackName, taskLaunchType)
	result := model.WrappedAppStack{}
	err := c.PostParsed("/v1/stacklistItems", body, &result)
	if err != nil {
		return err
	}

	if result.Record == nil {
		return errors.New("tried to create a record that already existed")
	}

	return nil
}

func (c *HappyClient) DeleteFromStacklist(appName, environment, stackName, taskLaunchType string) error {
	body := model.MakeAppStackPayload(appName, environment, stackName, taskLaunchType)
	result := model.WrappedAppStack{}
	err := c.DeleteParsed("/v1/stacklistItems", body, &result)
	if err != nil {
		return err
	}

	if result.Record == nil {
		return errors.New("tried to delete a record that did not exist")
	}

	return nil
}
