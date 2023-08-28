---
parent: Stacks
layout: default
has_toc: true
---

# Service Discovery in EKS

# What is Service Discovery?

Modern applications often consist of many different services which need to be able to contact each other. When services are running on static VMs, and aren't changing much, this is a pretty simple problem to solve. You set up some configurations to point to the VM's IP and Port, and you're off! In a more complex VM-based setup, maybe you use a load balancer with a target group.

But how do you achieve the same thing in the Kubernetes world?

In an environment where processes all have their own IPs, and they're dynamically scaled from moment to moment, how do you reliably and deterministically contact these services? This class of problem is one of "Service Discovery". Fortunately, Kubernetes has some tools to help. They are often analogous to the world outside of Kubernetes, so we'll draw those parallels when it helps clarify things.

Note that Kubernetes setups are widely varied and can be extremely complex. Here we speak to common use cases and configurations.

# Core DNS

At the heart of Kubernetes' built-in service discovery is a system service called `Core DNS`. Core DNS is responsible for providing dynamic DNS services to the containers running on Kubernetes.

Core DNS is configured primarily through use of `Service` objects. These objects do more than configure DNS, but this is the primary integration you'll be using.

# Pods

You'll hear the word `pod` a lot in conversations about Kubernetes. What is a `pod`?

A pod is the basic unit of work that Kubernetes can schedule. Usually, the pod consists of one container which it manages. (There are more complicated setups which put multiple containers in a single pod, but let's keep it simple for now.)

A pod is scheduled on a specific worker node, and typically represents one instance of a given program running inside a container.

Kubernetes manages scaling by creating more replicas of a pod, so as you scale, you get more and more pods.

In Core DNS, each pod will have a unique hostname consisting of its name and some random characters. It will also have its own randomly assigned IP address. The assigned IP address and hostname will be different each time a pod is started.

But if pod hostnames and IPs are random, how do I put them in my configuration? Answer: *You don't!* If you're putting a pod hostname or IP in a configuration anywhere, you're not going to have a good time. For this, you want a `Service`.

# Services

So what is a `Service`?

In its usual form, a service is a proxy that accepts requests on behalf of a group of pod replicas, and forwards requests to those pods. It keeps track of all of the instances of your software services that are being run, and will distribute requests between them, more or less fairly.

In terms of Service Discovery, when you create a Service, you give it a `name`. This name is NOT random, and automatically becomes a DNS entry in CoreDNS. It is coupled with the service's `namespace` and a base domain - and is directly accessible to other pods in Kubernetes.

Normally, these DNS entries follow the form:

```
<servicename>.<namespace>.cluster.local
```

If you are in the same namespace as the service, you can just use the service name. If you are in another namespace, you can use `<servicename>.<namespace>` or `<servicename>.<namespace>.svc.cluster.local`.

In these names, the `servicename` and `namespace` portions should be adjusted appropriately for your software deployment. The `svc.cluster.local` portion is a base domain for cluster internal DNS, and is usually not changed. However, it is configurable by administrators, and may be changed in rare cases.

Services should be used for all cluster-internal communication. You should NOT go out to an external load balancer and back in. Not only does that cause additional latency, it also adds financial cost. You should also not attempt to directly address pods.

# Ingress

So we use Services and Core DNS for communication between processes inside the cluster, but what about traffic coming from outside the cluster?

This is where an `Ingress` comes in.

An `Ingress` is to a `Service`, what a `Service` is to a `Pod`, in the sense that it stands upstream and distributes requests to services. Just like a Service is responsible for knowing how to get traffic to one or more Pods, an Ingress is responsible for getting external requests routed to various Services inside the cluster.

```
Internet --> Ingresses --> Services --> Pod Replicas
```

There are many types of ingresses, but two of the more common ones are the `nginx` ingress and an `AWS Load Balancer Controller`-based ingress.

The `nginx` ingress is the default for most simple projects, and is extremely common for low to moderate traffic services. It is very configurable and easy to set up, and tends to be inexpensive.

The `AWS Load Balancer Controller`-based ingress is specific to AWS's EKS, and actually creates and configures Application Load Balancers and Network Load Balancers outside of the cluster to route traffic. They can use the AWS Certificate Manager, and can load balance between AWS Availability Zones. They are generally better integrated into AWS and are more robust, but they are also more expensive. In some use cases, they are also not as flexible as the `nginx` ingress.

A single Ingress can route based on virtual host names, path prefixes, and other predicates.

# ExternalName Service Types

Core DNS provides dynamic DNS for internal Kubernetes services, but what if we wanted to provide a DNS alias for an external service?

For example, let's say we have an RDS database whose hostname we would like to abstract in our configurations. We want to call it `database` everywhere. We can create an `ExternalName` Service, which is a Service object, but is NOT a proxy like we talked about above. This kind of Service simply creates the DNS entry in Core DNS, pointing to another hostname much like a DNS CNAME.

Once you've created the ExternalName service, you should be able to resolve it in pods inside the namespace.

ExternalName services are for resolving external DNS entries. They are not themselves visible externally.

# Service Meshes

A Service Mesh is a more advanced piece of software like [Istio](https://istio.io/) or [Linkerd](https://linkerd.io/). Service meshes tend to handle many concerns, like inter-service transport encryption, service discovery, monitoring and metrics. They may be deployed (sometimes automatically) as "sidecars" with your Kubernetes pods, and transparently provide services.

Services meshes are outside the scope of this document, other than to bring them up as a possible solution for advanced use cases and complex applications.

# Further Reading

* [Core DNS](https://coredns.io)
* [Core DNS Kubernetes Plugin](https://coredns.io/plugins/kubernetes/)

## Kubernetes

* [Kubernetes Networking](https://kubernetes.io/docs/concepts/services-networking/)
* [ExternalName Services](https://kubernetes.io/docs/concepts/services-networking/service/#externalname)

## Service Meshes

* [Linkerd](https://linkerd.io/)
* [Istio](https://istio.io/)
