import { TerraformStack, TerraformVariable } from "cdktf";
import { Construct } from "constructs";
import { HappyIntegrationSecret } from "./integration_secret";
import { HappyNetworking } from "./networking";
import { HappyECSFargateService } from "./service";
import { HappyECSTaskDefinition } from "./task_def";
import { AWSRegion, ECSComputeLimit, Environment, EnvironmentVariables, HappyServiceMeta, ServiceType } from "./types";
export function makeName(meta: HappyServiceMeta): string {
    return [meta.env, meta.stackName, meta.serviceDef.name].join("-")
}

export interface HappyServiceProps {
    env: Environment,
    stackName: string,
    serviceName: string,
    serviceImage: string,
    serviceType: ServiceType,

    servicePort?: number,
    computeLimits?: ECSComputeLimit,
    environment?: EnvironmentVariables,
    healthCheckPath?: string,
    replicas?: number,
    region?: AWSRegion
}

export class HappyService extends TerraformStack {
    constructor(scope: Construct, id: string, config: HappyServiceProps) {
        super(scope, id)

        const serviceType = process.env.serviceType

        switch serviceType {
            case "PRIVATE":
            case "INTERAL":
            case "PUBLIC":
            default:
            // synth
        }

        const meta: HappyServiceMeta = {
            env: config.env,
            stackName: config.stackName,
            region: config.region || AWSRegion.WEST2,
            serviceDef: {
                name: config.serviceName,
                desiredCount: config.replicas || 2,
                port: config.servicePort || 8080,
                image: config.serviceImage,
                computeLimits: config.computeLimits || { cpu: ".25 vcpu", mem: "512" },
                serviceType: config.serviceType,
                healthCheckPath: config.healthCheckPath || "/healthcheck",
                environment: config.environment || [],
            }
        }

        const tags = {
            env: meta.env,
            stackName: meta.stackName,
            serviceName: meta.serviceDef.name,
            region: meta.region,
            image: meta.serviceDef.image,
            serviceType: meta.serviceDef.serviceType
        }

        const secret = new HappyIntegrationSecret(this, "int_secret", {
            secretID: "blah"
        })
        secret.secret.to

        const taskDef = new HappyECSTaskDefinition(this, "task_def", {
            executionRoleArn: config.executionRoleArn,
            meta,
            tags,
        })

        const networking = new HappyNetworking(this, "networking", {
            vpcID: config.vpcID,
            baseZoneName: config.baseZoneName,
            healthCheckPath: meta.serviceDef.healthCheckPath,
            meta,
            tags,
        })

        new HappyECSFargateService(this, "farget_service", {
            taskDef,
            networking,
            ecsClusterARN: config.clusterARN,
            meta,
            tags,
        })
    }
}


interface MyAppConfig {

}
class MyApp extends Construct {
    constructor(scope: Construct, id: string, config: MyAppConfig) {
        super(scope, id)
        // terraform resource 1

        // terraform resource 2

        // terraform resource 2
    }
}

new MyApp(scope, id, {
    devSettings...
})

new MyApp(scope, id, {
    stagingSettings...
})

new MyApp(scope, id, {
    prodSettings...
})
