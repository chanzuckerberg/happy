# Creating a Happy Environment

# Fogg

`Fogg` is a tool used by Chan Zuckerberg Initiative to manage the base AWS accounts where your Happy environments will run. It helps you generate basic Terraform which you can modify as needed to create the infrastructure for your application. 


## Set Up Your Terraform State Provider

Fogg's default configuration uses the S3 Terraform backend. This requires both an S3 bucket (for state data) and a DynamoDB table for state locking.

Documentation for the bucket and Dynamo setup can be found [on the Hashicorp web site](https://developer.hashicorp.com/terraform/language/settings/backends/s3). 

> Make sure that your DynamoDB table has a Partition Key named `LockID` with a type of `String`.

Once you have the bucket, table, and permissions set up according to the above documentation, you are ready to create your fogg project.


## Set up a Fogg Project

We will now initialize our Fogg project.

```
mkdir <projectdir> && cd <projectdir>
git init
fogg init

<Answer questions>

```

During the `fogg init` interview, you will be asked for the following information:

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

## Create a Happy Environment in Fogg

Now that we have a base `global` environment and a working Terraform, we can start setting up our Happy Stack environments.

These environments are also managed via Fogg, so we'll still be working with the `fogg.yml` file.

In `fogg.yml`, we will add our `env` top-level section containing our first environment.

```
envs:
  dev:
```

Running `fogg apply` will produce a `terraform/envs/dev` folder containing a `Makefile` and a `README`. 

You can add as many environments as you'd like under this `envs` key. Just remember to `fogg apply` each time.

But that's not very useful. Let's add the Happy built-in `route53` component. 

```
envs:
  dev:
    components:
      route53:
        depends_on:
          accounts: []
          components: []
        providers:
          aws:
```

(Again, running `fogg apply`.)

You should now see a `route53` folder inside of your `dev` folder, and it should contain some bits of Terraform code.

