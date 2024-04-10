# An Offline Job App

A typical application deployed using happy path, that is just a long running process (e.g. a queue consumer) that doesn't expose a port.

## Prerequistes

* Install the latest version of happy: `brew tap chanzuckerberg/tap && brew install happy`
* Make sure you have access to the czi-playground AWS environment

## Notes

* All stacks in this examples folder will be automatically cleaned up within 24 hours of creation; it is not intended for production usage
* All stacks are created in the czi-playground environment; all CZI employees should have access to this environmnet
