---
layout: home
title: Design
nav_order: 1
---
 
Static analysis to drive better, continuous security education.
 
 
# Design
 
SASTisfaction is a static-analysis-static-testing (SAST) tool used to create and deliver security education for developers. SASTisfaction is a Github App that scans PRs with [semgrep](https://semgrep.dev/). These scans are used to write security educational material as a peer reviewer comment in the PR.
 
![screenshot](https://user-images.githubusercontent.com/76011913/145272418-74cc247e-ca48-4f66-a4c9-c3c8dfa1283b.png)
 
This allows for
 
* reaching developers quickly and from inside Github
* studying real-world examples
* monitoring code changes over the long term
 
SASTisfaction scans for a [few things right now](categories):
 
* potentially harmful APIs, libraries and programming patterns 
* poor configuration options of the Django and Rails web frameworks
* secrets checked into Github
