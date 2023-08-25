---
parent: Development
title: Directory Structure
layout: default
---

# Directory Structure

Code for SASTisfaction is located at  [pkg/rulesets/](https://github.com/chanzuckerberg/sastisfaction/tree/main/pkg/rulesets) contains logic for stitching together all of the semgrep rules defined in that folder into a single semgrep config file, with read-only and excluded rules defined in the policy.yaml file located in SASTisfactions top level directory

* [pkg/semgrep/](https://github.com/chanzuckerberg/sastisfaction/tree/main/pkg/semgrep) 
  * contains logic for running semgrep, and for translating semgrep output into a format digestible by other core logic. 
* [pkg/ghapp/](https://github.com/chanzuckerberg/sastisfaction/tree/main/pkg/ghapp) 
  * contains logic for running the Github App web service and processing the webhook events.
* [pkg/cmd/](https://github.com/chanzuckerberg/sastisfaction/tree/main/pkg/cmd) 
  * contains logic for running the sastisfaction CLI.