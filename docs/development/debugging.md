---
parent: Development
layout: default
has_toc: true
---

# Debugging
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>


## Common Errors

* `executing from main: error executing root cmd: error running ghapp command: server killed or signal interupt: 1 error occurred: * interrupt`
  * The program was killed via a SIGTERM or SIGINT signal. This means someone pressed “CTR+C” to shut the program down (local case) or ECS shutdown the container using a SIGINT (production case).
  * This error usually doesn’t indicate a bug unless it was unexpected
* `context canceled`
  * The goroutines return this error when they are shutting down because their context was canceled. This is likely because the server failed or the program received a SIGTERM or SIGINT signal.
* `unable to parse webhook: unknown X-Github-Event in message: <some_event_type>`
  * We use the go-github module to parse the Github Events and this error shows up from time to time.
  * Here is the code [where we call it](https://github.com/chanzuckerberg/sastisfaction/blob/8c8b0b587d2f958594bd5c817876c91444d99956/pkg/ghapp/web.go#L36)
    * For some reason the abode code doesn’t seem to find some of the message types. I think its probably because the Github API is changing faster than the go-github module. Try to update to the latest version and see if the new version adds the missing messageType to the eventTypeMapping map.

## Debugging App

Debugging the Github App is much simpler than debugging a Github Action. One of the design goals of using the Github App vs the Github Action is that the local development environment should be as similar to the production environment. This makes it a lot easier to eliminate issues that are a result of a different environment locally vs. production. Here is a set of checklists to try and common error messages that might be useful when debugging SASTisfaction issues:

### Method

1. Identify the error
  * All errors are captured by Sentry.
  * All errors in Sentry come with a full stack trace of the error in question, with a small snippet of code. Use these stack traces to find the faulting function and line of code.
1. Identify the error’s source. SASTisfaction talks to many services so there are a few places to look when figuring out errors:
  * an error running the actual semgrep command
  * an error talking to Github
  * an error processing a message from Github
  * an error talking to Snowflake
  * an error caused by ECS shutting down the container
  * an error in the program logic
1. Reproduce locally
  * With the error and the source, try to reproduce the same error locally
    * Use VS code or similar IDE to run the SASTisfaction web server with debugging enabled (as in you should be able to set breakpoints) locally
    * Start a smee client and point it to your local development server
    * Create a test Github App and point it to your smee client
    * Send a webhook notification to your Github App through one of the following methods:
      * Through the Github UI, under the advanced section of your Github App
      * Installing the test Github App to a test repo and making a PR
      * Sending a crafted POST request to your local web server with something like postman
        * You might want to disable the crypto verification since that might be annoying to build yourself
    * Step through the event processing code until you reach the error location

## Debugging semgrep Rules

Debugging semgrep rules is not very difficult because of the tooling around semgrep. Here are a few tips for writing and debugging rules

### Method

The method I would suggest is a test-driven development style. Semgrep’s developer tooling makes it very easy to write tests, and run these tests quickly against any rules. This makes for a very tight development cycle.

1. Read the [semgrep developer documentation](https://semgrep.dev/docs/)
1. Navigate to the [semgrep playground](https://semgrep.dev/editor)
1. **Before writing any rule syntax**, write some example code that you would like to catch with your rule. 
  * Find cases online
  * Find cases in CZI codebases
  * Write some vulnerable code yourself!
1. Copy this code into the editor and annotate the lines you would like to be caught by your rule with a semgrep test annotation.
  * Since semgrep operates on a single file, you can add as many tests in a single file as you like
  * Make sure to include various types of tests and think of the different edge cases
1. Write your rule in the semgrep playground
1. Iterate on this rule until all tests pass
  * start with a simple rule and build on it
  * For more advanced rules, make sure to look at the AST that is generated to see how semgrep is parsing your code and how your rule checks against that AST
    * You may need to run semgrep locally to do this using the --dump-ast flag

Keep in mind too that semgrep is a very new product and growing every day. Keep an eye on their ongoing Github issues, engage with their engineers and watch for version changes to semgrep.

### Effectiveness

Once your rule is passing its tests, you can move it into the SASTisfaction repo. However, it is also important to add this rule to the repo’s [default policy file](https://github.com/chanzuckerberg/sastisfaction/blob/main/pkg/rulesets/policy.yaml). The policy file is what tells SASTisfaction what rules to run and what rules to simply collect data on. It is strongly advised to add new rules to the report-only-rules section of the policy file. This will tell SASTisfaction to run the rules, send metrics to our dashboards, but not to create Github comments or PRs. This allows us to see how noisy the rule is when run against our fleet of code.

After you have confirmed the rule produces high signal, remove the rule from report-only-rules section.