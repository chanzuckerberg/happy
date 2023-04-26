# Deploy an Application Using Happy

**[ This is largely taken from old internal docs as a starting point. Needs to be freshened up. ]**

* Write a docker-compose file capable of building docker images that will represent the services in your application.

It’s also recommended to use docker-compose to launch a local version of your application for local development!

* Create new fogg components in your ${myteam}-infra repository that will manage the happy path dev/staging/prod environments. It’s ok if these workspaces are empty to begin with!

* Create a .happy/config.json file in your application repository.

* Add terraform code to your application repository as necessary to define your application stacks.

* Modify application code to look for environment variables that control access to stack resources
