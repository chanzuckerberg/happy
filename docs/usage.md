---
layout: default
title: Usage
nav_order: 4
---

# Usage

## Default Rules
 
SASTisfaction by default runs all the rules in the repository. The list of rules can be found in the [`semgrep` folder](https://github.com/chanzuckerberg/sastisfaction/tree/main/pkg/rulesets/semgrep). These rules include checks for things like secrets, command injection, cross-site scripting, secure default in TLS configurations, and server-side request forgery.
 
## Scan Types
 
SASTisfaction comes with three primary scan types:
 
* code scan
* secret scan
* report-only scan


### Code Scan
The code scan looks in business logic and configuration files for potentially bad programming practices or dangerous APIs that might introduce issues. During the code scan, SASTisfaction ignores files such as test folders, migration files, mocks and Github files.
 
### Secret Scan
The secret scan looks for secrets that have been committed to any part of the code base. No files are ignored during the secret scan.
 
### Report-Only Scan
The report-only scan is used for testing and development of new rules. SASTisfaction developers use this scan type to collect analytics on new rules to calculate their signal-to-noise ratio before releasing them as part of the default rules set.
 
## Making Rules
 
If you have ideas for things you want to be checked in your code base, antipatterns you'd like to teach your new hires, code smells that give you the goosebumps, let us know or write one yourself!
 
Currently, our rule engine is [`semgrep`](https://semgrep.dev/) and rule development [is very easy](https://semgrep.dev/editor). All rules are currently stored in [`/pkg/rulesets/semgrep/<category>/<language>`](https://github.com/chanzuckerberg/sastisfaction/tree/main/pkg/rulesets/semgrep). Add a YAML file with your semgrep rule in one of those folders that you think would match, test it in the semgrep editor with some code snippets and make a PR. We would LOVE your input! See more info in our [developer docs](/development).
 
## Ignoring Rules
 
SASTisfaction will look for a file called "sastisfaction-policy.yaml" in the root of your repository. Use this file to configure rules to ignore. This might look like the following:
 
~~~yaml
---
report-only:
- "detected-heroku-key"
~~~
 
where "detected-heroku-key" is the name of the rule you'd like to exclude. Also tell [@security-engineering](mailto:security-eng@chanzuckerberg.com)! We love talking about security and if a rule isn't working for you, we want to hear that. Maybe it can be improved for other folks too.
 
## About
 
### Non-blocking CI/CD
 
SASTisfaction does not block the PR from getting merged like a failed test. It tries very hard to provide as much security knowledge without getting in the developer’s way. SASTisfaction does this by leaving comments and educational material near the code that looks suspicious. 
 
SASTisfaction does listen to responses in the comments, so if you have suggestions or want to report a false-positive, leave us a note in the comment or reach us at [@security-engineering](mailto:security-eng@chanzuckerberg.com).
 
### Timing
 
As of this writing, SASTisfaction takes less than 1 minute. However, it will vary based on the size of your codebase since `semgrep` scans all files from the PR. It will likely also grow over time as more rules are added to the default set.
