import { Construct } from "constructs";
import { App, TerraformStack, TerraformVariable } from "cdktf";
import { AwsProvider } from "../../.gen/providers/aws"
import { Route53Record, DataAwsRoute53Zone } from "../../.gen/providers/aws/route53"

class HappyService extends TerraformStack {
    constructor(scope: Construct, name: string) {
        super(scope, name);
        //TODO: assume role
        new AwsProvider(this, "default", {

        })

        const prefix = new TerraformVariable(this, "prefix", {
            type: "string",
            description: "The name of the zone record",
        });

        const zoneName = new TerraformVariable(this, "zone_name", {
            type: "string",
            description: "The zone ID of the record to attach to.",
        });

        const zone = new DataAwsRoute53Zone(this, "dns_record", {
            name: zoneName.value
        })

        new Route53Record(this, "happy_stack_record", {
            name: `${prefix.value}.${zone.name}`,
            zoneId: zone.zoneId,
            type: "A"
        })
    }
}

const app = new App();
new HappyService(app, "happy-service");
app.synth();
