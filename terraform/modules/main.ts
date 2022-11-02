import { Construct } from "constructs";
import { App, TerraformOutput, TerraformStack, TerraformVariable } from "cdktf";
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";
import { DataAwsRoute53Zone } from "@cdktf/provider-aws/lib/data-aws-route53-zone";
import { Route53Record } from "@cdktf/provider-aws/lib/route53-record";
import {DataAwsAlb} from "@cdktf/provider-aws/lib/data-aws-alb"


class HappyDNS extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new AwsProvider(this, "AWS", {
      region: "us-west-2",
    });

    const appName = new TerraformVariable(this, "app_name", {
      type: "string",
      description: "The name of the ECS service",
    });
    const stackName = new TerraformVariable(this, "stack_name", {
      type: "string",
      description: "The name of the Happy stack"
    })
    const zoneName = new TerraformVariable(this, "zone_name", {
      type: "string",
      description: "The zone ID of the record to attach to",
    });
    const albARN = new TerraformVariable(this, "alb_arn", {
      type: "string",
      description: "The ARN of the ALB to use for this happy stack",
    });

    const alb = new DataAwsAlb(this, "alb", {
      arn: albARN.toString(),
    })

    const zone = new DataAwsRoute53Zone(this, "zone", {
      name: zoneName.toString()
    })

    const prefix = `${stackName.toString()}-${appName.toString()}`
    new Route53Record(this, "happy_stack_record", {
      name: `${prefix}.${zone.name}`,
      zoneId: zone.zoneId,
      type: "A",
      alias: [{
        name: alb.name,
        zoneId: alb.zoneId,
        evaluateTargetHealth: false,
      }]
    })

    new TerraformOutput(this, "dns_prefix", {
      value: prefix,
    })
  }
}

const app = new App();
new HappyDNS(app, "happy-dns");
app.synth();
