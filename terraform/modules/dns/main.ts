import { Construct } from "constructs";
import { App, TerraformStack, TerraformVariable } from "cdktf";
import { DataAwsRoute53Zone } from "@cdktf/provider-aws/lib/data-aws-route53-zone";
import { Route53Record } from "@cdktf/provider-aws/lib/route53-record";
import { DataAwsAlb } from "@cdktf/provider-aws/lib/data-aws-alb"
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";

export class HappyDNS extends Construct {
  record: Route53Record

  constructor(scope: Construct, id: string, albName: string, albZoneId: string, recordName: string, zoneId: string) {
    super(scope, id)
    this.record = new Route53Record(this, "happy_stack_record", {
      name: recordName,
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

    const recordName = new TerraformVariable(this, "record_name", {
      type: "string",
      description: "The name of the record to create"
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

    new HappyDNS(this, "happy-dns", alb.name, alb.zoneId, recordName.stringValue, zone.zoneId)
  }
}

const app = new App();
const stack = new HappyDNSModule(app, "happy-dns");
new AwsProvider(stack, "AWS", {
  region: "us-west-2",
});
app.synth();
