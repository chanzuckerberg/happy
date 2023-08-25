---
parent: Development
title: Testing
layout: default
---

# Testing

SASTisfaction has local unit and integration tests. It is preferred to run tests by using our Github CI. Test automatically run on each PR and each CI run has its environment set up properly to execute the tests. Create a branch, write some code, create a PR, mark PR as draft, and watch the tests execute. If tests fail, fix the code, commit to your branch, and push to your branch. The tests will run on every push. Once all tests are passing in our CI, mark the PR ready for review and assign a reviewer.

## Local Testing

Tests can be run locally, however, keep in mind that some of the environment variables will need to be set up to perform integration testing (for instance, we have integration tests with an actual Snowflake database to ensure our queries are crafted properly). To do so, make use of chamber to execute these tests. Here is an example of how to run a test locally with chamber:

~~~bash
AWS_PROFILE=czi-sec chamber exec sec-czi-sec-sastisfaction -- go test -timeout 30s -run ^TestValidQuery$ github.com/chanzuckerberg/sastisfaction/pkg/ghapp
~~~

Or simply run all the tests:

~~~bash
AWS_PROFILE=czi-sec make test
~~~

If you’re developing new functionality, please add to these local tests! 

## Running Locally

If you want to run the server locally, you can do so like: 

~~~bash
AWS_PROFILE=czi-sec make run
~~~

This will start a web server on [localhost:3000](http://localhost:3000) and use the environment stored in the production SSM parameters. If you are developing against a test Github App, you’ll need to set up your own environment variables [as shown below](#forwarding-requests-locally). 

To receive Github webhooks, [you will need a proxy](#forwarding-requests-locally) that can forward them from a public endpoint.

### Forwarding Requests Locally

1. Create a Github App in the [Settings -> Developer Settings -> New Github App section](https://github.com/settings/apps)
1. Install smee locally:
  ~~~bash
  npm install -g smee-client
  ~~~
1. [Create a Smee channel](https://smee.io/new) and paste the channel URL in the following:
  ~~~bash
  smee -u [CHANNEL_URL] --path /event_handler --port 3000 
  ~~~
1. Set up the smee channel to be the "Webhook URL" of the Github App
1. Request the following permissions (this might change up as we play with this app):
  * read/write -- contents
  * read/write -- discussions
  * read-only -- metadata
  * read/write -- pull requests
1. Generate an RSA private key at the bottom of the app settings page
1. Create a webhook secret (any text will do, doesn't need to be super secret) and put it in webhook secret section
1. Grab the app ID from the top of the Github App settings page
1. Install the app on one of your private repo (if you don't have one, just create one and install your test app there)
1. Create a .env in pkg/ghapp/ file with the following values:
  * `GITHUB_PRIVATEKEY="<PRIVATE_KEY_ONE_LINE>"` (needs to be on 1 line)
  * `GITHUB_WEBHOOKSECRET="<WEBHOOK_SECRET>"`
  * `GITHUB_APPIDENTIFIER=<APP_ID>`
1. Run the `pkg/ghapp/app.go` file: `cd pkg/ghapp; go run app.go ghapp`

Now, whenever you trigger a check_suite event on your Github repo that has your development app set up, you should get data forwarding through to your app. This allows you to iterate by using the Github UI to modify files and commit changes, causing check webhooks to be sent to your development server.