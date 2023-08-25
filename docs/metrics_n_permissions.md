---
layout: default
title: Metrics and Permissions
nav_order: 5
---
 
# Metrics
 
SASTisfaction collects some data for analytics purposes:
 
* runtime errors of the Github App
* Github webhook information, including a link to the PR that triggered the code scan
* static analysis rule violations and their code locations
* code scan duration
* interactions with SASTisfaction comments
 
# Permissions
 
![permissions screenshot](https://user-images.githubusercontent.com/76011913/146958173-acc58ea5-4334-476b-a399-3b32a4883c86.png)
 
The Github app requests the following permissions:
 
* read/write access [to checks](https://docs.github.com/v3/apps/permissions/#permission-on-checks)
  * used to update the list of checks that are performed on a PR
* read/write access [to content](https://docs.github.com/v3/apps/permissions/#permission-on-contents)
  * used to perform a shallow clone of the PR to be analyzed locally by semgrep
* read/write access [to discussions](https://docs.github.com/v3/apps/permissions/)
  * used to collect interactions with comments made on PRs
  * SASTisfaction interact with these comments and also uses them to gauge issue validity
* read-only [to metadata](https://docs.github.com/v3/apps/permissions/#metadata-permissions)
  * a required permission of all GH apps
* read/write [to pull requests](https://docs.github.com/v3/apps/permissions/#permission-on-pull-requests)
  * used to know when a pull request is made on the repo
