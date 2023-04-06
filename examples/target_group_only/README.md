# Target Group Only

This is an advanced example, but one that is helpful for hooking up a new happy service with
existing infrastructure. Suppose you have the scenario where you application is already be hosted by an
ALB controlled by some infrastructure elsewhere. The usual happy behavior creates an ALB for your stack. 
However, since you already have an existing ALB, you'd prefer to hook up happy to it inside.

## Prerequistes

* Install the latest version of happy: `brew tap chanzuckerberg/tap && brew install happy`
* Make sure you have access to the czi-playground AWS environment

## Notes

* This is an advanced method for bridging two pieces of infrastructure not all managed through happy
* For brand new happy projects, its not advised to use this example. Utilize the INTERNAL, PRIVATE, or EXTERNAL service types

## Process

Normally, happy will automatically deploy your services using an ALB. However, if you don't need an ALB 
or want to hook up your services to an existing ALB, the
normal behavior is problematic. In this example, we are going to solve this problem by telling happy
we don't need an ALB created. Do this by setting the [service_type to "TARGET_GROUP_ONLY"](./.happy/terraform/envs/rdev/main.tf). Additionally, we pass in a new set of options:

~~~json
{
    path              = "/mypath"        // path to attach the target group of the ALB below
    health_check_path = "/mypath/health" // the healthcheck route should be below the path specified
    port              = 3000             // port of the service (see app.js)

    // the external ALB and listener to attach to
    alb = {
        name          = aws_lb.this.name
        listener_port = aws_lb_listener.this.port
    }
}
~~~

The required information to perform this operation is the name of the existing ALB and the listener port number
you want to attach you happy service. Additionally, provide a path that happy service will be reachable on.