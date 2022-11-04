import { CloudwatchLogGroup } from "@cdktf/provider-aws/lib/cloudwatch-log-group";
import { EcsTaskDefinition } from "@cdktf/provider-aws/lib/ecs-task-definition";
import { Construct } from "constructs";
import { makeName } from "./main";
import { ContainerDefinition, HappyServiceMeta } from "./types";

export class HappyECSTaskDefinition extends Construct {
    readonly taskDef: EcsTaskDefinition
    readonly containerDefs: ContainerDefinition[]

    constructor(
        scope: Construct,
        id: string,
        meta: HappyServiceMeta
    ) {
        super(scope, id)

        const logGroup = new CloudwatchLogGroup(scope, meta.stackName, {
            retentionInDays: 365,
            name: `/happy/${meta.env}/${meta.stackName}/${meta.serviceDef.name}`
        })

        this.containerDefs = [{
            name: meta.serviceDef.name,
            memory: meta.serviceDef.computeLimits.mem,
            image: meta.serviceDef.image,
            essential: true,
            environment: meta.serviceDef.environment,
            portMappings: [{ containerPort: meta.serviceDef.port }],
            logConfiguration: {
                logDriver: 'awslogs',
                options: {
                    "awslogs-group": logGroup.id,
                    "awslogs-region": meta.region,
                    "awslogs-stream-prefix": meta.env,
                }
            },
        }]

        this.taskDef = new EcsTaskDefinition(scope, meta.stackName, {
            family: makeName(meta),
            memory: meta.serviceDef.computeLimits.mem,
            cpu: meta.serviceDef.computeLimits.cpu,
            networkMode: "awsvpc",
            requiresCompatibilities: ["FARGATE"],
            taskRoleArn: "TODO",
            executionRoleArn: "TODO",
            containerDefinitions: JSON.stringify(this.containerDefs),
        })
    }
}
