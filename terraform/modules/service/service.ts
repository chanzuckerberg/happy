import { EcsService } from "@cdktf/provider-aws/lib/ecs-service";
import { Construct } from "constructs";
import { makeName } from "./main";
import { HappyNetworking } from "./networking";
import { HappyECSTaskDefinition } from "./task_def";
import { HappyServiceMeta } from "./types";
import { Pet } from "@cdktf/provider-random/lib/pet";

interface HappyECSFargateServiceProps {
    ecsClusterARN: string,
    meta: HappyServiceMeta,
    taskDef: HappyECSTaskDefinition,
    networking: HappyNetworking,
    tags?: { [key: string]: string };
}
export class HappyECSFargateService extends Construct {
    constructor(
        scope: Construct,
        id: string,
        config: HappyECSFargateServiceProps,
    ) {
        super(scope, id)

        const pet = new Pet(scope, "petecsservice", { prefix: makeName(config.meta) })
        new EcsService(scope, "ecs_fargate_service", {
            cluster: config.ecsClusterARN,
            desiredCount: config.meta.serviceDef.desiredCount,
            taskDefinition: config.taskDef.taskDef.arn,
            launchType: "FARGATE",
            name: pet.id,
            loadBalancer: config.taskDef.containerDefs.map(c => ({
                containerName: c.name,
                containerPort: c.portMappings[0].containerPort,
                targetGroupArn: config.networking.targetGroup.arn
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
