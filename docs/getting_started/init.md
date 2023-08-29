---
parent: Getting Started
layout: default
has_toc: true
---

# Init

There are several [configuration options](../config/index.md) right now for creating happy. Hopefully, this number decreases over time, but
for now there are handful of things that need to be set in the right place for happy to know where and what to deploy from you application.

## Bootstrap a Greenfield Project

When starting a new greenfield project from an empty repository, use the following command to have happy initialize a project for you:

~~~
$ happy bootstrap --force
~~~

The `happy bootstrap` command will walk you through a survey of questions and code generate all the correct configuration files. This is the 
easiest way to get started. For this command to run properly, you will need to have AWS credentials loaded in your environment and you will
also need a [happy environment](./deploy_first_env.md) already set up.

## Copy an Example as a Template

All the folders in the [examples](https://github.com/chanzuckerberg/happy/tree/main/examples) directory are fully tested and fleshed out
happy projects. They can be used as templates for new projects. Copy the directoy that fits your needs into your application and start
tinkering!