// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import "cdktf/lib/testing/adapters/jest"; // Load types for expect matchers
import { TerraformStack, Testing } from "cdktf";
import { HappyECSTaskDefinition } from "../service/task_def"
import { HappyServiceMeta, ServiceType, ECSComputeLimit } from "../service/types"
import { CloudwatchLogGroup } from "@cdktf/provider-aws/lib/cloudwatch-log-group";
import { EcsTaskDefinition } from "@cdktf/provider-aws/lib/ecs-task-definition";
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";

describe("happy service", () => {
  const meta: HappyServiceMeta = {
    env: "rdev",
    stackName: "myStackName",
    region: "us-west-2",
    serviceDef: {
      name: "myServiceName",
      desiredCount: 1,
      port: 8080,
      image: "dockerImage",
      computeLimits: { cpu: ".25 vcpu", mem: "1 GB" } as ECSComputeLimit,
      serviceType: "private" as ServiceType,
      healthCheckPath: "/healthcheck"
    }
  }

  it("creates a AWS Route53 record", () => {
    const construct = Testing.synthScope((scope) => {
      new HappyECSTaskDefinition(scope, "myTaskDefinition", meta);
    })

    expect(construct).toHaveResourceWithProperties(CloudwatchLogGroup, {
      retention_in_days: 365,
      name: "/happy/rdev/myStackName/myServiceName",
    });

    expect(construct).toHaveResourceWithProperties(EcsTaskDefinition, {
      family: "rdev-myStackName-myServiceName",
      memory: "1 GB",
      cpu: ".25 vcpu",
      network_mode: "awsvpc",
      requires_compatibilities: ["FARGATE"],
      execution_role_arn: "TODO",
      task_role_arn: "TODO",
    });
  });
  describe("can synthesize so", () => {
    const app = Testing.app();
    const stack = new TerraformStack(app, "testHappyService");
    new HappyECSTaskDefinition(stack, "myTaskDefinition", meta);
    new AwsProvider(stack, "AWS", {
      region: "us-west-2",
    });
    const synth = Testing.synth(stack)
    const fullSynth = Testing.fullSynth(stack)

    it("it looks like it did before", () => {
      expect(synth).toMatchSnapshot();
    });

    it("is valid terraform", () => {
      expect(fullSynth).toBeValidTerraform();
    });

    it("plans without error", () => {
      expect(fullSynth).toPlanSuccessfully()
    });
  })
});
