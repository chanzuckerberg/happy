# Creating a Happy Environment


## Set Up Your Terraform State Provider

Fogg's default configuration uses the S3 Terraform backend. This requires both an S3 bucket (for state data) and a DynamoDB table for state locking.

Documentation for the bucket and Dynamo setup can be found [on the Hashicorp web site](https://developer.hashicorp.com/terraform/language/settings/backends/s3). 

> Make sure that your DynamoDB table has a Partition Key named `LockID` with a type of `String`.

Once you have the bucket, table, and permissions set up according to the above documentation, you are ready to create your fogg project.


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
* infra bucket name: This is the S3 bucket you created above
* infra dynamo table: This is the table you created above
* auth profile: AWS profile name
* owner: An email address who will be listed as a primary contact for the project
* AWS Account ID: An AWS account ID which will be the primary account used for managing infrastructure

You will then have a `fogg.yml` file in your base project directory.

You can now "apply" this fogg configuration to generate Terraform.

```
fogg apply
```

You should now have a directory structure with default terraform and other scripts and files.

## Initialize Terraform

```
cd terraform/global
terraform init
```

This will initialize the S3 backend data and make it the default backend for all of your commands going forward.

You can test it by running:

```
$ terraform plan

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```


##  