import { App } from "cdktf";
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";
import { RandomProvider } from "@cdktf/provider-random/lib/provider"
import { HappyService } from "./service/main";
import { Environment, ServiceType } from "./service/types";

const envs: Environment[] = [Environment.RDEV, Environment.STAGING, Environment.PROD]

// the below code represents what used to be everything in the "stack"
// module. For each stack, this will be used to make the stack.
// A happy service is the module that builds the networking, compute,
// and container definitions.
const app = new App();
envs.map(env => {
    // TODO: pull these in using a config file or something
    [
        {
            env,
            serviceImage: "image",
            serviceName: "service1",
            stackName: "myTestStack",
            // ex: computeLimits: { cpu: "1 vcpu", mem: "4 GB" },
            // ex: replicas: 4,
            serviceType: ServiceType.EXTERNAL,
        },
        {
            env,
            serviceImage: "image2",
            serviceName: "service2",
            stackName: "myTestStack2",
            serviceType: ServiceType.INTERNAL,
        },
        {
            env,
            serviceImage: "image3",
            serviceName: "service3",
            stackName: "myTestStack3",
            serviceType: ServiceType.PRIVATE,
        },
    ].map(s => {
        const service = new HappyService(app, `happyapp-${s.serviceName}`, s);

        // TODO: figure out how the user configures their provider
        // I feel like we could pull this in from the config file
        new AwsProvider(service, "AWS", {
            region: "us-west-2",
        });
        new RandomProvider(service, "Random", {
        })
    })
})

// Add additional stack resources here:
// define custom resource 1 (maybe an lambda)
// define custom resource 2 (maybe a migration task)

app.synth();
