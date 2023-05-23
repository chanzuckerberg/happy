# An API exposed through an SSL-enabled sidecar

This app is composed of a single pod with a main container, and a sidecar that acts as an SSL proxy.  This example demonstrates a few things:

1. A service that is exposing the sidecar port (8443) as opposed to the main container port (3000).
2. Volume mounting a secret, and using this in the side car (SSL certs and key).
3. Volume mounting a configmap (the stacklist) and using this in the main container.

If you hit the following endpoint: https://sidecar-with-ssl-api.happy-playground-rdev.rdev.si.czi.technology/stacklist it will return
a list of stacks.

## Prerequistes

* Install the latest version of happy: `brew tap chanzuckerberg/tap && brew install happy`
* Make sure you have access to the czi-playground AWS environment
* This examples assumes cert-manager is provisioned.

## Notes

* All stacks in this examples folder will be automatically cleaned up within 24 hours of creation; it is not intended for production usage
* All stacks are created in the czi-playground environment; all CZI employees should have access to this environmnet
* Sidecar container is a secondary container, deployed alongside the primary application container, typically an agent, an initializer, a bootstrapper, a database schema migrator, etc. Sidecar container resides within the same container network as the primary application container (make sure ports don't collide), and either application can call the other one. Some sidecars can intercept traffic heading to the application (for the purposes of TLS termination, of traffic offloading).