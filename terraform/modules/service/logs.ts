import { Construct } from "constructs";
import { CloudwatchLogGroup } from "@cdktf/provider-aws/lib/cloudwatch-log-group";
import { AWSLogConfig, HappyServiceMeta } from "./types";

interface HappyLogsProps {
    meta: HappyServiceMeta,
    tags?: { [key: string]: string };
}

export class HappyLogs extends Construct {
    readonly logGroup: CloudwatchLogGroup
    readonly logConfig: AWSLogConfig
    constructor(scope: Construct, id: string, config: HappyLogsProps) {
        super(scope, id)

        this.logGroup = new CloudwatchLogGroup(scope, `cwg`, {
            retentionInDays: 365,
            name: `/happy/${config.meta.env}/${config.meta.stackName}/${config.meta.serviceDef.name}`,
            tags: config.tags,
        })

        this.logConfig = {
            logDriver: 'awslogs',
            options: {
                "awslogs-group": this.logGroup.id,
                "awslogs-region": config.meta.region,
                "awslogs-stream-prefix": config.meta.env,
            }
        }
    }
}
