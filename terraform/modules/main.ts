import { Construct } from "constructs";
import { App, TerraformOutput, TerraformStack, TerraformVariable } from "cdktf";
import { DataAwsRoute53Zone } from "@cdktf/provider-aws/lib/data-aws-route53-zone";
import { Route53Record } from "@cdktf/provider-aws/lib/route53-record";
import { DataAwsAlb } from "@cdktf/provider-aws/lib/data-aws-alb"

export class HappyDNS extends Construct {
  prefix: string
  record: Route53Record

  constructor(scope: Construct, id: string, albName: string, albZoneId: string, zoneName: string, zoneId: string, stackName: string, appName: string) {
    super(scope, id)
    console.log(albName, albZoneId)
    this.prefix = `${stackName}-${appName}`
    this.record = new Route53Record(this, "happy_stack_record", {
      name: `${this.prefix}.${zoneName}`,
      zoneId: zoneId,
      type: "A",
      alias: [{
        name: albName,
        zoneId: albZoneId,
        evaluateTargetHealth: false,
      }]
    })
  }
}

export class HappyDNSModule extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

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
      arn: albARN.stringValue,
    })

    const zone = new DataAwsRoute53Zone(this, "zone", {
      name: zoneName.stringValue
    })

    const dns = new HappyDNS(this, "happy-dns", alb.name, alb.zoneId, zone.name, zone.zoneId, stackName.stringValue, appName.stringValue)
    new TerraformOutput(this, "dns_prefix", {
      value: dns.prefix,
    })
  }
}

const app = new App();
new HappyDNSModule(app, "happy-dns");

app.synth();
