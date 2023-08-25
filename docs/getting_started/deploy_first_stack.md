---
parent: Happy Stacks
layout: default
has_toc: true
---

# Deploy Your First Stack
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Summary

Generally, a happy environment would need to be set up before you can start deploying stacks. But that is boring. It is much more fun to 
to start using happy and make your first application. Luckily, we have a sandbox environment to do just that. This sandbox environment is
strictly for CZI use only and all stacks in this environment are automatically cleaned up at the end of the week. Here's how to get started:

* Clone https://github.com/chanzuckerberg/happy/tree/main
* Navigate to ./happy/tree/main/examples/typical_app
* Execute `happy create <stackname>` where `<stackname>` is the name of the stack
* You're done!

This stack will deploy two services. To see the endpoints to access them, use `happy list` and it should look like the following:

~~~
$ happy list
[INFO]: Listing stacks from the happy api
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│ NAME     │ OWNER    │ APP     │ REPO                                    │ BRANCH                     │ HASH                           │ STATUS  │ URLS                                                          │ LASTUPDATED │
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
│ typical1 │ alokshin │ typical │ https://github.com/chanzuckerberg/happy │ alokshin/global-name-check │ dirty git tree (PLEASE COMMIT  │ applied │ https://typical1.happy-playground-rdev.rdev.si.czi.technology │ 176h34m39s  │
│          │          │         │                                         │                            │ YOUR CHANGES)                  │         │                                                               │             │
│──────────│──────────│─────────│─────────────────────────────────────────│────────────────────────────│────────────────────────────────│─────────│───────────────────────────────────────────────────────────────│─────────────│
~~~

Change the Dockerfile, application code, configuration however you like and redeploy using `happy update <stackname>`