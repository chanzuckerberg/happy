# Image Pipeline   

This example shows how you can use happy to create an base container image that can be utilized by other Happy applications. Use in conjunction
with the nodejs_downstream example to test and end-to-end image pipeline example.

## Prerequistes

* Install the latest version of happy: `brew tap chanzuckerberg/tap && brew install happy`
* Make sure you have access to the czi-playground AWS environment

## Notes

* All stacks in this examples folder will be automatically cleaned up within 24 hours of creation; it is not intended for production usage
* All stacks are created in the czi-playground environment; all CZI employees should have access to this environmnet

## Process

This image_pipeline project is setup as a service_type ["IMAGE_TEMPLATE"](./.happy/terraform/envs/rdev/main.tf). This means it will not deploy an application, but only push container images to ECRS. This is ideal for using happy to develop a base container image that others can pull from. The flow should be:

* Develop your base image
* Push to your dev environment with `happy push <stackname> --env rdev`
* Test
* Iterate
* Promote image to the next environment with `happy push <stackname> --env prod`
* Test
* Tag your image with semantic version: `happy push <stackname> --env prod --tag v1.0.1`

Once your base image is ready to be consumed by a downstream happy project, that happy project can pull from the prod ECR of the base image using
the semantic versioning tags. The nodejs_downstream folder does this. Here is an example of pulling the base image in a Dockerfile:

~~~
FROM 626314663667.dkr.ecr.us-west-2.amazonaws.com/test/prod/base:v1.0.1
~~~

## Exercise

### Create a base image in rdev

~~~
cd nodejs_base
happy build
happy push test # test can be subed for another stack name
~~~

Take notes of the ECR that is created and pushed to, the happy output should look something like this:

~~~
The push refers to repository [626314663667.dkr.ecr.us-west-2.amazonaws.com/test/prod/base]
e001cc98f077: Layer already exists 
bc22ffdf8988: Layer already exists 
e09f336caa5d: Layer already exists 
fa916c369c49: Layer already exists 
ee420dfed78a: Layer already exists 
jheath-chanzuckerberg.com-2023-04-04T16-38-50: digest: sha256:f6b51e6d27430bf4126826cf475fcea594b6cfa3dca65178a818c1d9f2b7ccb3 size: 1369
~~~

### Create a human-readable tag for image in prod

~~~
happy push test --tag v1.0.1 --env prod
~~~

### Use the base image in another project

~~~
cd nodejs_downstream
# update the Dockerfile to point to the ECR image
echo "FROM 626314663667.dkr.ecr.us-west-2.amazonaws.com/test/prod/base:v1.0.1" >> ./frontend/Dockerfile
echo "FROM 626314663667.dkr.ecr.us-west-2.amazonaws.com/test/prod/base:v1.0.1" >> ./backend/Dockerfile

# build the backend
cd backend
happy create myapp
~~~