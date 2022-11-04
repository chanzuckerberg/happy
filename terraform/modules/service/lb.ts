
import { Construct } from "constructs";
import { LbListener } from "@cdktf/provider-aws/lib/lb-listener"
import { Lb } from "@cdktf/provider-aws/lib/lb"
import { DataAwsRoute53Zone } from "@cdktf/provider-aws/lib/data-aws-route53-zone";
import { AwsAcmCertificate } from '../.gen/modules/aws-acm-certificate';
import { HappyServiceMeta } from "./types"
import { makeName } from "./main";

export interface HappyLoadBalancer {
    readonly lbName: string
    readonly lbZoneID: string
    readonly recordName: string
    readonly listener: LbListener
    readonly baseZone: DataAwsRoute53Zone
}

export class HappyInternalLoadBalancer extends Construct implements HappyLoadBalancer {
    readonly lbName: string
    readonly lbZoneID: string
    readonly recordName: string
    readonly baseZone: DataAwsRoute53Zone
    readonly listener: LbListener
    constructor(scope: Construct, id: string, meta: HappyServiceMeta, baseZoneName: string) {
        super(scope, id)

        // This requires the *.internal.* as part of the name so that the nginx/oauth proxy
        // can properly route traffic to this load balancer. This zone will already exist
        // from the happy environment's multi-domain-oauth-proxy.
        this.baseZone = new DataAwsRoute53Zone(scope, "data_record", {
            name: `internal.${baseZoneName}`
        })
        this.recordName = `${meta.stackName}-${meta.serviceDef.name}.${this.baseZone.name}`
        const lb = new Lb(scope, "lb", {
            name: makeName(meta),
            internal: true,
            securityGroups: [], // TODO
            subnets: [], //TODO
            idleTimeout: 30,
        })

        this.lbName = lb.name
        this.lbZoneID = lb.zoneId
        this.listener = new LbListener(scope, "lb_listener", {
            loadBalancerArn: lb.arn,
            port: 80,
            protocol: "HTTP",
            defaultAction: [{
                type: "fixed-response",
                fixedResponse: {
                    contentType: "text/plain",
                    messageBody: "Not Found",
                    statusCode: "404",
                }
            }]
        })
    }
}

export class HappyExternalLoadBalancer extends Construct implements HappyLoadBalancer {
    readonly lbName: string
    readonly lbZoneID: string
    readonly baseZone: DataAwsRoute53Zone
    readonly listener: LbListener
    readonly recordName: string

    constructor(scope: Construct, id: string, meta: HappyServiceMeta, baseZoneName: string) {
        super(scope, id)

        const lb = new Lb(scope, "lb", {
            name: makeName(meta),
            internal: false,
            securityGroups: [], // TODO
            subnets: [], //TODO
            idleTimeout: 30,
        })
        this.lbName = lb.name
        this.lbZoneID = lb.zoneId
        this.baseZone = new DataAwsRoute53Zone(scope, "data_record", {
            name: baseZoneName
        })
        this.recordName = `${meta.serviceDef.name}.${this.baseZone.name}`
        const cert = new AwsAcmCertificate(scope, "lb_cert", {
            awsRoute53ZoneId: "TODO",
            certDomainName: this.recordName,
            tags: {},
        })
        this.listener = new LbListener(scope, "lb_listener", {
            loadBalancerArn: lb.arn,
            port: 443,
            protocol: "HTTPS",
            sslPolicy: "ELBSecurityPolicy-TLS-1-2-Ext-2018-06",
            certificateArn: cert.arnOutput,
            defaultAction: [{
                type: "fixed-response",
                fixedResponse: {
                    contentType: "text/plain",
                    messageBody: "Not Found",
                    statusCode: "404",
                }
            }]
        })
    }
}
