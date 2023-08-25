---
parent: Rule Categories
title: SSRF
layout: default
has_toc: true
---
 
## Server-Side Request Forgery (SSRF)
 
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
 
An attacker can make requests as if they were on your server. This might allow them to make requests to things that weren’t meant to be exposed to the Internet, such as a database, cloud configuration or internal application.
 
### Scenario
 
Jane is hosting a chat app at `https://janesite.com`. Whenever someone posts a link in her chat app, the app creates a preview box in the chat window of the linked website. Jane is also hosting her private, personal website at `https://127.0.0.1:8089` on the same network as `https://janesite.com`. `https://127.0.0.1:8089` is not exposed to the Internet. Jack sends the link `https://127.0.0.1:8089` in Jane's chat app. Jane's app downloads `https://127.0.0.1:8089` and displays a preview of Jane's private website to Jack's chat window.
 
### Reasoning
 
This often happens because when businesses host applications inside internal networks, it can be presumed that if you are on that internal network, you are automatically authenticated and authorized. For example, if a company has a VPN to interact with company applications and resources, it might be assumed by these company applications that if the user is on the VPN, they are automatically allowed to access the application. When SSRF is exploited on applications on a privileged network like this, they might be considered authenticated by other applications or resources on the same network. 
 
Attackers exploit this scenario to send requests from the authenticated server to internal resources that would normally not be exposed to the public Internet. This could include things like caches, databases, internal applications or other networks.
 
### Solution
 
#### Validate Links
 
Servers should validate any untrusted inputs from their applications before processing them. This includes user inputs used to build links before making a request to them. The specific validation will be dependent on your infrastructure. Generally, links should NOT resolve to any of the destination IP addresses:
 
* 127.0.0.1
* 10.0.0.0 — 10.255.255.255
* 172.16.0.0 — 172.31.255.255
* 192.168.0.0 — 192.168.255.255
* 169.254.169.254
 
#### Make the Request Through a Proxy
 
An alternative solution would be to deploy a web proxy outside of your applications network and route the request through this proxy. This works because any requests that are trying to access internal resources of your network will be blocked by your network's normal firewall rules.
 
#### Enable IMDSv2 for AWS
 
If your application is hosted in AWS, opt in to disable IMDSv1 and use [IMDSv2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html) instead. This prevents an SSRF issue from being able to exfiltrate the AWS credentials of the machine the website is hosted on. This feature will only limit the AWS impact of an SSRF attack and not solve the issue entirely. Services that are reachable from the vulnerable website could still be vulnerable.
 
### Fun Facts
 
* There was some debate about [adding SSRF to OWASP's top 10](https://owasp.org/Top10/A10_2021-Server-Side_Request_Forgery_%28SSRF%29/) for 2021 citing a low incidence rate.
 
## References
 
### Read
 
These are good resources to read more advanced docs about this issue:
 
* [Advanced SSRF Definition](https://portswigger.net/web-security/ssrf)
* [OWASP SSRF Definition](https://owasp.org/www-community/attacks/Server_Side_Request_Forgery)
 
### Play
 
These are good interactive labs that allow you to play with real SSRF:
 
* [Portswigger SSRF Lab](https://portswigger.net/web-security/ssrf/lab-basic-ssrf-against-localhost)
* [Portswigger Vulnerable Labs](https://portswigger.net/web-security/all-labs)
* [Google Gruyere](https://google-gruyere.appspot.com/)