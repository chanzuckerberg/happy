import { App } from "cdktf";
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";
import { RandomProvider } from "@cdktf/provider-random/lib/provider"
import { HappyService, HappyServiceProps } from "./service/main";
import { exit } from "process";

const vpcID = process.env.VPC_ID
const executionRoleArn = process.env.EXECUTION_ROLE_ARN
const clusterARN = process.env.CLUSTER_ARN
if (!vpcID || !executionRoleArn || !clusterARN) {
    exit(-1)
}

// TODO: pull these in using a config file or something
const services: HappyServiceProps[] = [
    {
        env: "rdev",
        serviceImage: "image",
        serviceName: "service1",
        stackName: "myTestStack",
        serviceType: "EXTERNAL",
        baseZoneName: "czi.technology",
        vpcID,
        executionRoleArn,
        clusterARN,
    },
    {
        env: "rdev",
        serviceImage: "image",
        serviceName: "service2",
        stackName: "myTestStack",
        serviceType: "INTERNAL",
        baseZoneName: "czi.technology",
        vpcID,
        executionRoleArn,
        clusterARN,
    },
]

const app = new App();
services.map(s => {
    const stack = new HappyService(app, `happyapp-${s.serviceName}`, s);
    new AwsProvider(stack, "AWS", {
        region: "us-west-2",
    });
    new RandomProvider(stack, "Random", {
    })
})

app.synth();
