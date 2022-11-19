import { Construct } from "constructs";
import { LbTargetGroup } from "@cdktf/provider-aws/lib/lb-target-group"
import { LbListenerRule } from "@cdktf/provider-aws/lib/lb-listener-rule"
import { HappyDNS } from "../dns/main";
import { HappyServiceMeta } from "./types"
import { HappyExternalLoadBalancer, HappyInternalLoadBalancer, HappyLoadBalancer } from "./lb";

interface HappyNetworkingProps {
    vpcID: string,
    baseZoneName: string,
    healthCheckPath: string,
    meta: HappyServiceMeta,
    tags?: { [key: string]: string }
}
export class HappyNetworking extends Construct {
    readonly lb: HappyLoadBalancer
    readonly targetGroup: LbTargetGroup
    constructor(
        scope: Construct,
        id: string,
        config: HappyNetworkingProps,
    ) {
        super(scope, id)

        if (config.meta.serviceDef.serviceType === "EXTERNAL") {
            this.lb = new HappyExternalLoadBalancer(scope, "ext_lb", { meta: config.meta, baseZoneName: config.baseZoneName })
        } else {
            this.lb = new HappyInternalLoadBalancer(scope, "int_lb", { meta: config.meta, baseZoneName: config.baseZoneName })
        }

        // the DNS record that will be used to hit this service. It will point to the happy load balancer above
        const dns = new HappyDNS(
            scope,
            "record",
            this.lb.lbName,
            this.lb.lbZoneID,
            this.lb.recordName,
            this.lb.baseZone.name
        )

        // health checking and port routing to the specific container
        this.targetGroup = new LbTargetGroup(scope, "lb_target_group", {
            vpcId: config.vpcID,
            port: config.meta.serviceDef.port,
            protocol: "HTTP",
            targetType: "ip",
            deregistrationDelay: "10",
            healthCheck: {
                interval: 15,
                path: config.healthCheckPath,
                protocol: "HTTP",
                timeout: 5,
                healthyThreshold: 2,
                unhealthyThreshold: 10,
                matcher: "200-299"
            }
        })

        // if the host matches the record that was created, take the request to the container
        new LbListenerRule(scope, "lb_listener_rule", {
            listenerArn: this.lb.listener.arn,
            condition: [{
                hostHeader: { values: [dns.record.name] },
            }],
            action: [{
                targetGroupArn: this.targetGroup.arn,
                type: "forward"
            }]
        })
    }
}
