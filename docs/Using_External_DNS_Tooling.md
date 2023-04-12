# External DNS Tooling with EKS

Inside of our cluster, Core DNS, Services, and Ingresses [help us with service discovery](Service_Discovery_In_EKS.md), but how do users and services on the internet, or in other clusters find our EKS instance?

# Ingress Endpoints

Kubernetes Ingresses, regardless of type, must expose some sort of externally facing hostname. Load Balancer ingresses on EKS will expose the hostname for the Load Balancer which was created to service the Ingress. Nginx ingresses will expose a hostname which points to the cluster master nodes. But these names are typically unfriendly, and, in the case of AWS, will be reflective of the AWS infrastructure, not something recognizable as our services or application.

Kubernetes ingresses are also capable of "virtual hosting" -- transparently serving many DNS names at once -- so we can't simply give out the bare ingress endpoint (more on that later).

# Simple Hosting

In the simplest case, we can route all requests to our ingress endpoint as a single API, either sending all requests to a single service, or matching on `/url/path/` and sending requests to the appropriate service. In this case, the default Ingress Endpoing MAY be used, because the hostname is not involved in routing. However, it is not very pretty, and you'll want to have a nicer name outside of testing environments.

So how do I get a nicer name for my ingress endpoint?

The simple answer is to add a CNAME record for it in your DNS zone. You can use Route53 for this if you maintain your DNS zones on AWS, but this works regardless of DNS server type.

> TIP: A `CNAME` is basically a hostname in DNS which functions as an "alias" to another DNS name. The alias and the target do not need to be in the same zone. You don't even need to own the target to create a CNAME alias to it!

Once you've added a CNAME pointing to the ingress endpoint's hostname, you will be able to use your new name in URLs (possibly after a DNS caching delay).

In the case of simple hosting where we are only matching on URL path, or sending all traffic to a single service, this will work almost immediately, and simply serves as an easy-to-remember way to name your endpoint services.

# Virtual Hosting

Simple hosting is easy to understand and set up, but what if we have many logical APIs to host for our app, and we want them all on different DNS names? We don't really want to set up new ingress endpoints for every service. ALB-based ingresses in particular can get pretty pricey if allowed to multiply.

Fortunately, Kubernetes Ingresses support what's called "Domain Based Virtual Hosting". What this means is that the ingress router will look at the hostname portion of the requested URL, in addition to its path. That hostname will be used in routing the request, along with any other matching your ingress does.

The trick is to make sure that your ingress is matching on hostname, and then go back to your DNS software and create as many CNAMEs as you need, pointing to your common ingress endpoint.

When a client looks up your friendly names, the CNAMEs will all be resolved by the global DNS system to point to the same, single Kubernetes ingress hostname exposed by our cluster. All requests for all of those "apparently" different APIs will funnel into the same place.

The Ingress definition will first match the hostname requested, and then proceed to match URL paths as normal. This allows each "apparent" API server to have identical paths that route to different services in the backend.
