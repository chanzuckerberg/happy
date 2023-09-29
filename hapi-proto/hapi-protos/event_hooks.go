package hapi_protos

import context "context"

func (e *HappyEventORM) BeforeToORM(ctx context.Context, happyEventORM *HappyEventORM) error {
	return nil
}

var _ HappyEventWithBeforeToORM = &HappyEventORM{}
