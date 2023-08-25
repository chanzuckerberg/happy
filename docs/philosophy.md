---
layout: default
title: Philosophy
nav_order: 3
---
 
# Philosophy
 
SASTisfaction is all about security education.
 
It switches the mentality from "find all the vulnerabilities" to "provide security education". By doing this, the complexity of the static analysis rules are relaxed. Instead of worrying about high signal bugs with complex [taint analysis](https://deepsource.io/glossary/taint-analysis/), it focuses on finding the places where it can easily insert a reminder about a general security principle. It makes the messaging clear and digestible. It communicates thoughtful, educational advice, as opposed to "this is a bug, block the build" messages. The author of the code has the option to either learn and revise or ignore the comment. Either way, the tool communicates the most important security concepts in a continuous, consistent manner. 
 
### Rules Engine with `semgrep`
 
Since the tool focuses on education, the rule engine needed
 
* a simple syntax to write rules so developers could contribute
* a simple way to execute rules
* lots of language support
* community adoption and contributions
 
Many of the commercial products for static analysis didn't tick all these boxes. Many are very heavy and focus a lot on taint analysis. Many require a complex environment to run them (for example, a virtual Windows environment). Many don't support Ruby.
 
Semgrep is a great open-source tool that smartly "greps" codebases for security signal. SASTisfaction uses that signal to educate on those points. In the long term, the rules engine could change or be supplemented with more comprehensive rules engines. SASTisfaction is written such that the rules engine is an implementation detail and we can swap it or add to its data as needed. Whether the rules engine changes or not, SASTisfaction to bring security education close to where developers work.