# A Cerbos Sidecar-Enabled App

A custom version of a sidecar application with the cerbos deployed as a sidecar. See https://docs.cerbos.dev/cerbos/latest/deployment/k8s-sidecar.html for more details.

## Prerequistes

* Install the latest version of happy: `brew tap chanzuckerberg/tap && brew install happy`
* Make sure you have access to the czi-playground AWS environment

## Notes

* All stacks in this examples folder will be automatically cleaned up within 24 hours of creation; it is not intended for production usage
* All stacks are created in the czi-playground environment; all CZI employees should have access to this environmnet
* Sidecar container is a secondary container, deployed alongside the primary application container, typically an agent, an initializer, a bootstrapper, a database schema migrator, etc. Sidecar container resides within the same container network as the primary application container (make sure ports don't collide), and either application can call the other one. Some sidecars can intercept traffic heading to the application (for the purposes of TLS termination, of traffic offloading). In this specific setup an an example provided sidecar responds on port 80; main application, when invoked on `/sidecar` endpoint, proxies the call to the sidecar and responds to the client.