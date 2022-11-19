import { TerraformStack, TerraformVariable } from "cdktf";
import { Construct } from "constructs";
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
    executionRoleArn: string,
    vpcID: string,
    baseZoneName: string,
    clusterARN: string,

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
        const meta: HappyServiceMeta = {
            env: config.env,
            stackName: config.stackName,
            region: config.region || "us-west-2",
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

        const taskDef = new HappyECSTaskDefinition(this, "task_def", {
            executionRoleArn: config.executionRoleArn,
            meta,
            tags
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

export class HappyServiceModule extends TerraformStack {
    constructor(scope: Construct, id: string) {
        super(scope, id)

        const env = new TerraformVariable(this, "env", {
            type: "string",
            description: "The environment this service is being deployed in (rdev, dev, staging, prod, etc...)"
        })
        const stackName = new TerraformVariable(this, "stack_name", {
            type: "string",
            description: "The name of the stack"
        })
        const serviceName = new TerraformVariable(this, "service_name", {
            type: "string",
            description: "The name of the service"
        })
        const servicePort = new TerraformVariable(this, "service_port", {
            type: "number",
            description: "The port on the container to expose"
        })
        const serviceImage = new TerraformVariable(this, "service_image", {
            type: "string",
            description: "The Docker image to deploy"
        })
        const serviceType = new TerraformVariable(this, "service_type", {
            type: "string",
            description: "The type of service to deploy: EXTERNAL, INTERNAL, PRIVATE"
        })
        const clusterARN = new TerraformVariable(this, "cluster_arn", {
            type: "string",
            description: "The ARN of the cluster to run this service on"
        })
        const vpc = new TerraformVariable(this, "vpc_id", {
            type: "string",
            description: "The VPC ID the happy service is operating in"
        })
        const baseZoneName = new TerraformVariable(this, "base_zone_name", {
            type: "string",
            description: "The name of the zone to attach load balancers and DNS records for this service (ex: example.com would create appName-stackName.example.com URLs for rdev stacks)"
        })


        const replicas = new TerraformVariable(this, "replicas", {
            type: "number",
            default: 2,
        })
        const cpu = new TerraformVariable(this, "cpu", {
            type: "string",
            default: ".25 vcpu",
        })
        const mem = new TerraformVariable(this, "mem", {
            type: "string",
            default: "512",
        })
        const healthCheckPath = new TerraformVariable(this, "health_check_path", {
            type: "string",
            default: "/healthcheck"
        })
        const envVars = new TerraformVariable(this, "env_vars", {
            type: "list(object({name:string, value:string}))",
            default: [],
        })
        const { stringValue: executionRoleArn } = new TerraformVariable(this, "task_execution_role", {
            type: "string",
            description: "The task execution role created for the ECS cluster to use on the services."
        })

        const meta: HappyServiceMeta = {
            env: env.stringValue as Environment,
            stackName: stackName.stringValue,
            region: "us-west-2",
            serviceDef: {
                name: serviceName.stringValue,
                desiredCount: replicas.numberValue,
                port: servicePort.numberValue,
                image: serviceImage.stringValue,
                computeLimits: { cpu: cpu.stringValue, mem: mem.stringValue } as ECSComputeLimit,
                serviceType: serviceType.value as ServiceType,
                healthCheckPath: healthCheckPath.stringValue,
                environment: envVars.value as EnvironmentVariables
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

        const taskDef = new HappyECSTaskDefinition(this, "task_def", {
            executionRoleArn,
            meta,
            tags
        })

        const networking = new HappyNetworking(this, "networking", {
            vpcID: vpc.stringValue,
            baseZoneName: baseZoneName.stringValue,
            healthCheckPath: healthCheckPath.stringValue,
            meta,
            tags,
        })

        new HappyECSFargateService(this, "farget_service", {
            taskDef,
            networking,
            ecsClusterARN: clusterARN.stringValue,
            meta,
            tags,
        })
    }
}
