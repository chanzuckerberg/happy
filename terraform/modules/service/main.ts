import { Construct } from "constructs";
import { App, TerraformOutput, TerraformStack, TerraformVariable } from "cdktf";
import { DataAwsRoute53Zone } from "@cdktf/provider-aws/lib/data-aws-route53-zone";
import { Route53Record } from "@cdktf/provider-aws/lib/route53-record";
import { DataAwsAlb } from "@cdktf/provider-aws/lib/data-aws-alb"
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";
import { EcsTaskDefinition, } from "@cdktf/provider-aws/lib/ecs-task-definition"
import { EcsService, } from "@cdktf/provider-aws/lib/ecs-service"
import { LbTargetGroup } from "@cdktf/provider-aws/lib/lb-target-group"
import { LbListenerRule } from "@cdktf/provider-aws/lib/lb-listener-rule"
import { LbListener } from "@cdktf/provider-aws/lib/lb-listener"
import { CloudwatchLogGroup } from "@cdktf/provider-aws/lib/cloudwatch-log-group"

type ECSComputeLimit = { cpu: ".25 vcpu", mem: "512" } |
{ cpu: ".25 vcpu", mem: "1 GB" } |
{ cpu: ".25 vcpu", mem: "2 GB" } |
{ cpu: ".5 vcpu", mem: "1 GB" } |
{ cpu: ".5 vcpu", mem: "2 GB" } |
{ cpu: ".5 vcpu", mem: "3 GB" } |
{ cpu: ".5 vcpu", mem: "4 GB" } |
{ cpu: "1 vcpu", mem: "2 GB" } |
{ cpu: "1 vcpu", mem: "3 GB" } |
{ cpu: "1 vcpu", mem: "4 GB" } |
{ cpu: "1 vcpu", mem: "5 GB" } |
{ cpu: "1 vcpu", mem: "6 GB" } |
{ cpu: "1 vcpu", mem: "7 GB" } |
{ cpu: "1 vcpu", mem: "8 GB" } |
{ cpu: "2 vcpu", mem: "4 GB" } |
{ cpu: "2 vcpu", mem: "5 GB" } |
{ cpu: "2 vcpu", mem: "6 GB" } |
{ cpu: "2 vcpu", mem: "7 GB" } |
{ cpu: "2 vcpu", mem: "8 GB" } |
{ cpu: "2 vcpu", mem: "9 GB" } |
{ cpu: "2 vcpu", mem: "10 GB" } |
{ cpu: "2 vcpu", mem: "11 GB" } |
{ cpu: "2 vcpu", mem: "12 GB" } |
{ cpu: "2 vcpu", mem: "13 GB" } |
{ cpu: "2 vcpu", mem: "14 GB" } |
{ cpu: "2 vcpu", mem: "15 GB" } |
{ cpu: "2 vcpu", mem: "16 GB" } |
{ cpu: "1 vcpu", mem: "4 GB" }

enum ServiceType {
    Private,
    Internal,
    Public
}

type Environment = "rdev" | "dev" | "staging" | "prod"

type AWSRegion = "us-east-1" | "us-east-2" | "us-west-1" | "us-west-2"

interface ECSServiceDefinition {
    name: string
    // the number of replicas
    desiredCount: number
    // the port to expose
    port: number
    // the image to use in the container
    image: string
    // deployment size (only valid combinations allowed)
    computeLimits: ECSComputeLimit
    // if the service is on the internet, protected by Okta, or only exposed within the cluster
    serviceType: ServiceType
    healthCheckPath?: string
    environment?: EnvironmentVariables
    command?: string
}

interface EnvironmentVariables { name: string, value: string }[]

interface AWSLogConfig {
    logDriver: "awslogs",
    options: {
        "awslogs-stream-prefix": string,
        "awslogs-group": string,
        "awslogs-region": AWSRegion,
    }
}

interface ContainerDefinition {
    name: string,
    image: string,
    memory: string,
    portMappings: {
        containerPort: number,
        hostPort?: number,
        protocol?: string
    }[]
    logConfiguration?: AWSLogConfig
    essential?: boolean,
    environment?: EnvironmentVariables
    command?: string
}


function makeName(env: Environment, stackName: string, serviceDef: ECSServiceDefinition): string {
    return [env, stackName, serviceDef.name].join("-")
}

class HappyLoadBalancer extends Construct {
    readonly targetGroup: LbTargetGroup
    constructor(
        scope: Construct,
        id: string,
        vpcID: string,
        priority: number,
        hostMatch: string,
        serviceDef: ECSServiceDefinition
    ) {
        super(scope, id)

        this.targetGroup = new LbTargetGroup(scope, "lb_target_group", {
            vpcId: vpcID,
            port: serviceDef.port,
            protocol: "HTTP",
            targetType: "ip",
            deregistrationDelay: "10",
            healthCheck: {
                interval: 15,
                path: "/",
                protocol: "HTTP",
                timeout: 5,
                healthyThreshold: 2,
                unhealthyThreshold: 10,
                matcher: "200-299"
            }
        })

        const lbListener = new LbListener(scope, "lb_listener", {})

        new LbListenerRule(scope, "lb_listener_rule", {
            listenerArn: lbListener.arn,
            priority: priority,
            condition: [{
                hostHeader: { values: [hostMatch] },
            }],
            action: [{
                targetGroupArn: this.targetGroup.arn,
                type: "forward"
            }]
        })
    }
}

class HappyECSTaskDefinition extends Construct {
    readonly taskDef: EcsTaskDefinition
    readonly containerDefs: ContainerDefinition[]

    constructor(
        scope: Construct,
        id: string,
        env: Environment,
        stackName: string,
        region: AWSRegion,
        serviceDef: ECSServiceDefinition
    ) {
        super(scope, id)

        const logGroup = new CloudwatchLogGroup(scope, stackName, {
            retentionInDays: 365,
            name: `/happy/${env}/${stackName}/${serviceDef.name}`
        })

        this.containerDefs = [{
            name: serviceDef.name,
            memory: serviceDef.computeLimits.mem,
            image: serviceDef.image,
            essential: true,
            environment: serviceDef.environment,
            portMappings: [{ containerPort: serviceDef.port }],
            logConfiguration: {
                logDriver: 'awslogs',
                options: {
                    "awslogs-group": logGroup.id,
                    "awslogs-region": region,
                    "awslogs-stream-prefix": env,
                }
            },
            command: serviceDef.command,
        }]

        this.taskDef = new EcsTaskDefinition(scope, stackName, {
            family: makeName(env, stackName, serviceDef),
            memory: serviceDef.computeLimits.mem,
            cpu: serviceDef.computeLimits.cpu,
            networkMode: "awsvpc",
            requiresCompatibilities: ["FARGATE"],
            taskRoleArn: "TODO",
            executionRoleArn: "TODO",
            containerDefinitions: JSON.stringify(this.containerDefs),
        })
    }
}

class HappyECSFargateService extends Construct {
    constructor(
        scope: Construct,
        id: string,
        env: Environment,
        stackName: string,
        serviceDef: ECSServiceDefinition,
        ecsClusterARN: string,
        def: HappyECSTaskDefinition,
        lb: HappyLoadBalancer
    ) {
        super(scope, id)

        new EcsService(scope, "ecs_fargate_service", {
            cluster: ecsClusterARN,
            desiredCount: serviceDef.desiredCount,
            taskDefinition: def.taskDef.arn,
            launchType: "FARGATE",
            name: makeName(env, stackName, serviceDef),
            loadBalancer: def.containerDefs.map(c => ({
                containerName: c.name,
                containerPort: c.portMappings[0].containerPort,
                targetGroupArn: lb.targetGroup.arn
            })),
            networkConfiguration: {
                securityGroups: [],
                subnets: [],
                assignPublicIp: false,
            },
            enableExecuteCommand: true,
            waitForSteadyState: true,
        })
    }
}
