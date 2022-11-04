import { EcsService } from "@cdktf/provider-aws/lib/ecs-service";
import { Construct } from "constructs";
import { makeName } from "./main";
import { HappyNetworking } from "./networking";
import { HappyECSTaskDefinition } from "./task_def";
import { HappyServiceMeta } from "./types";


export class HappyECSFargateService extends Construct {
    constructor(
        scope: Construct,
        id: string,
        ecsClusterARN: string,
        meta: HappyServiceMeta,
        def: HappyECSTaskDefinition,
        networking: HappyNetworking
    ) {
        super(scope, id)

        new EcsService(scope, "ecs_fargate_service", {
            cluster: ecsClusterARN,
            desiredCount: meta.serviceDef.desiredCount,
            taskDefinition: def.taskDef.arn,
            launchType: "FARGATE",
            name: makeName(meta),
            loadBalancer: def.containerDefs.map(c => ({
                containerName: c.name,
                containerPort: c.portMappings[0].containerPort,
                targetGroupArn: networking.targetGroup.arn
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
