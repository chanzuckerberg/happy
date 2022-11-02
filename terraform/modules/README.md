# How to develop in this repo

## Tests

* Use main-test.ts as an example for how to test your code. There are 3 different tests you should run with `yarn test`:
  * unit tests
    * make sure that the resources in your constructs are made as you expect
  * snapshot tests
    * jest will take a snaphot of the synthesized output. this is to make sure why you are refactoring, the output is the exact same as before. If it needs to be different be sure to run `yarn test -u` to modify the snapshot. Make sure to commit all the snapshots. If your snapshot is different and you didn't expect the module to be any different you might have done something wrong
  * terraform tests
    * make sure `terraform validate` and `terraform plan` execute without errors
* Writing these tests first. It will greatly increase your confidence and speed while iterating on these modules
* To successfully execute a `terraform plan` you will need to use something like `AWS_PROFILE=czi-playground yarn test` so that the provider is configured with AWS credentials, otherwise you will get an error

## Structure

* Isolate your non-module components into classes that inherit from Construct. Think of Constructs that are reusable blocks that can be ingested by other CDK constructs, not just HCL. Ideally, we'll want to invoke these CDK constructs outside of HCL so isolate the dependencies and make your constructs more general purpose
* If you are making a module to be consumed by other terraform HCL files, make a module component that inherits from TerraformStack.
* Add all your TerraformStacks to the application that needs to be synthesized.
* The caller of the modules should configure the provider, don't add a provider for every construct. For example, our test code configures the AWS provider and the HCL that uses these modules configures the AWS provider
