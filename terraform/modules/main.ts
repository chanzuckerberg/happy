import { App } from "cdktf";
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";
import { RandomProvider } from "@cdktf/provider-random/lib/provider"
import { HappyService, HappyServiceProps } from "./service/main";

const env = "rdev"
// TODO: pull these in using a config file or something
const services: HappyServiceProps[] = [
    {
        env,
        serviceImage: "image",
        serviceName: "service1",
        stackName: "myTestStack",
        // ex: computeLimits: { cpu: "1 vcpu", mem: "4 GB" },
        // ex: replicas: 4,
        serviceType: "EXTERNAL",
    },
    {
        env,
        serviceImage: "image2",
        serviceName: "service2",
        stackName: "myTestStack2",
        serviceType: "INTERNAL",
    },
    {
        env,
        serviceImage: "image3",
        serviceName: "service3",
        stackName: "myTestStack3",
        serviceType: "PRIVATE",
    },
]

// the below code represents what used to be everything in the "stack"
// module. For each stack, this will be used to make the stack.
// A happy service is the module that builds the networking, compute,
// and container definitions.
const app = new App();
services.map(s => {
    const stack = new HappyService(app, `happyapp-${s.serviceName}`, s);
    new AwsProvider(stack, "AWS", {
        region: "us-west-2",
    });
    new RandomProvider(stack, "Random", {
    })
})
// Add additional stack resources here:
// define custom resource 1
// define custom resource 2

app.synth();
