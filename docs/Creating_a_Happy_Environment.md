# Creating a Happy Environment

## Set up a Fogg Project

`Fogg` is a tool used by Chan Zuckerberg Initiative to manage the base AWS accounts where your Happy environments will run. It helps you generate basic Terraform which you can add to in order to create the infrastructure for your application. 



```
mkdir <projectdir> && cd <projectdir>
git init
fogg init

<Answer questions>

```

During the init interview, you will be asked for the following information:

* project name: A name for your project (e.g. `hello-world`)
* aws region: AWS region (eg `us-east-1`)
* infra bucket name: Valid S3 bucket name to store Terraform states
* infra dynamo table: Valid DynamoDB table name (e.g. `infra-data`)
* auth profile: AWS profile name
* owner: An email address who will be listed as a primary contact for the project
* AWS Account ID: An AWS account ID which will be the primary account used for managing infrastructure

You will then have a `fogg.yml` file in your base project directory.

You can now "apply" this fogg configuration to generate Terraform.

```
fogg apply
```

You should now have a directory structure with default terraform and other scripts and files.

