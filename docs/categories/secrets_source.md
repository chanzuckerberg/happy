---
parent: Rule Categories
title: Secrets in Source Code
layout: default
has_toc: true
---
 
## Secrets in Source Code
 
{: .no_toc }
 
<details open markdown="block">
 <summary>
   Table of contents
 </summary>
 {: .text-delta }
1. TOC
{:toc}
</details>
 
### TL;DR
 
Lots of developers, contractors, and potentially other people can see production secrets in the application's source code. While source code might be needed by many different people, production secrets should only be known by a select few people because they could allow access to sensitive production data, such as user sessions, databases, or other sensitive systems.
 
### Scenario
 
Jane is writing the code for `https://janesite.com`. She tests her access to her AWS database by copy/pasting her access and secret key into the source code and running the site locally. After everything is working, she commits her code, including the credentials used to test her access. Jack is hired on as a contractor and receives access to the source code. Jack views the credentials and uses them to log into Jane's AWS account.
 
### Reasoning
 
Often, the reason this happens is because storing secrets properly is not always easy. It is much easier to put a secret value in a string or configuration file and commit it to the source code. This might even be accidental in that the developer was only testing something locally and forgot to remove the value before committing the code. In this case, the secret value is in the repository history and any user with access to the repository can see it.
 
Secrets don't belong in source code because not everyone who has access to the source code needs access to secrets. This is a great example of [principle of least privilege](https://en.wikipedia.org/wiki/Principle_of_least_privilege). The source code of repositories have broad ranges of access from developers, to contractors, to even the open Internet. Secrets should only be accessed by a select few system administrators. 
 
### Solution
 
#### Limit Access to Secrets
 
Secrets should only be required by a select group of administrators. Secrets belong in a secure location that
 
* can enforce the principle of least privilege
* can be audited for their access
* can easily rotate the secrets
* offer encryption at rest.
 
#### Rotate Secrets Committed to Code
 
If code is found to have a secret in it, the secret should be removed from the repository, then rotated. Most version control systems will maintain a copy of every commit made to the repository. Even after the secret has been removed, its value is stored in an older commit copy. Anyone with access to the repository has access to this history. By rotating the secret, it is impossible for anyone to use the old secret value.
 
### Fun Facts
 
* If your project is open sourced, many bots are constantly scanning your repository for secrets, [including Github](https://github.blog/2021-06-08-securing-open-source-supply-chain-scanning-package-registry-credentials/).
* Many tools exist to search for secrets in Git histories, such as [trufflehog](https://github.com/trufflesecurity/truffleHog).
 
## References
### Read
 
These are good resources to read more advanced docs about this issue:
 
* [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
* [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
* [Hashicorp Vault](https://www.vaultproject.io/)
* [1Password](https://1password.com/)