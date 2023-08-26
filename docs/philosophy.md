---
layout: default
title: Philosophy
nav_order: 3
---
 
# Philosophy
 
Happy comes from the mindset of making the secure thing easy. 

At CZI, we have lots of little and medium sized applications. These apps are generally simple HTTP web applications.
Many times, developers don't have the capacity or time to wade through the 1000s of ways to deploy or stand up their application.
Happy is now the de-facto way of deploying applications in minutes at CZI. For folks with more advanced use cases, it is a great
starting point to iterate on, while being robust and flexible enough to add advanced infrastructure or features. For folks with simple 
us cases, it just works.


## Technologies 

It uses the following technologies:

* Containerized applications
* Docker compose
* EKS or ECS clusters
* Terraform
* Happy CLI tool
* Github Actions

## Features 

It gives your application out-of-the-box:

* Auditing and logging
* End-to-end encryption
* Incident response
* Load balancing and autoscaling
* Service discovery with other happy stacks
* Development environments
* Shell access to stacks
* Automatic deployment workflows
 
## Composition vs Prescription
 
Happy is more of a composition of tools rather than a rigid prescription used to deploy an application. We found
that each application is different in small ways and reflects the culture and style of each team. To try and accomidate
all the styles and teams would make the tool useless. And to force teams to fit a single style would stiffle new ideas.

We settled on a composition of tools that work together to provide an easy deployment path. If you are more familiar with
`kubectl`, there is nothing stopping you from exploring your cluster that way. If you like to look at your deployment 
through AWS console, go ahead. If you hate using docker-compose for local development, don't use it.

For those less familiar with infrastructure or find yourself saying "I just want to deploy my app and see it", then happy
provides a defined path and set of tools to do that. 