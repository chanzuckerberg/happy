---
parent: Rule Categories
title: Missing HSTS Header
layout: default
has_toc: true
---
 
## Missing HTTP Strict-Transport-Security (HSTS) Header
 
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
 
Attackers can trick users to navigate to an unencrypted version of a website, allowing them to view and modify all the traffic.
 
### Scenario
 
Jane is running the website `https://janesite.com`. Jane keeps her unencrypted HTTP version `http://janesite.com` available because she didn’t think it mattered, but redirects all requests to her new, encrypted site `https://janesite.com`. Jack logs in to an open Wi-Fi network near his house to view the old, unencrypted site at `http://janesite.com`. Bill owns the open Wi-Fi network Jack connected to and can see the unencrypted traffic from Jack. Bill sees the redirected response from `http://janesite.com` to `https://janesite.com`, but changes the response to be a copy of the response from `https://janesite.com`. Bill does this for all network requests Jack makes to `http://janesite.com`. Bill can now see all the activity of Jack.
 
### Reasoning
 
The reason this can happen is because the history of the web started with HTTP, not HTTPS. As the networks the web was built on got more complex with sensitive applications such as online banking, payment transactions, and PII, the web needed to be secured with strong encryption (i.e. HTTPS). However, it could not be enabled by default for all browsers without breaking many websites. So HTTPS is now an opt-in feature and browsers are allowed to navigate to both HTTP and HTTPS sites as a normal part of their functionality.
 
Attackers can abuse this feature of web browsers to force users to navigate to insecure versions of websites. On a network without encryption, an attacker can trivially read and write the traffic of users who navigate to HTTP sites. Even when websites try to redirect users' browsers to secure versions of their site, the attacker can read and write the user's traffic, allowing them to block the redirect and keep the user navigating over HTTP.
 
HSTS helps deal with this issue. HSTS tells the browser to never go back to that website over HTTP. Even if an attacker attempts to phish a user with an HTTP link, the browser will change the link to be HTTPS. Since the browser will no longer send traffic over HTTP, the attacker can no longer read and write their traffic, even on an unencrypted network.
 
### Solution
 
#### Enable HSTS
 
Add the following HTTP response header to all servers that are served over HTTPS:
 
~~~http
Strict-Transport-Security: max-age=31536000; includeSubDomains
~~~
 
The `max-age` value is a number representing seconds. The above time is one year (60 * 60 * 24 * 365 = 3153600). Make sure the `max-age` value is sufficiently long. Most security guides will recommend one year.
 
#### HSTS Preload List
 
Popular browsers such Chrome and Firefox preload lists of domains to be included in HSTS for the browser so that users never have to navigate to HTTP. [Add your domain to these preload lists.](https://hstspreload.org/)
 
### Fun Facts
 
* HSTS is not that old (2016). It was designed after a [proof-of-concept](https://www.youtube.com/watch?v=MFol6IMbZ7Y) by [Moxie0](https://twitter.com/moxie) for [SSLStrip](https://github.com/moxie0/sslstrip).
* HSTS is simple to enable if all websites are using HTTPS. It is surprisingly difficult to enable for organizations that host websites, such as blogs or marketing pages, over HTTP since these sites will break if HSTS is enabled.
 
## References
 
### Read
 
These are good resources to read more advanced docs about this issue:
 
* [Advanced HSTS Definition](https://portswigger.net/kb/issues/01000300_strict-transport-security-not-enforced)
* [Google Web Security](https://developers.google.com/web/fundamentals/security/?hl=en)
* [Mozilla Web Security](https://infosec.mozilla.org/guidelines/web_security)
* [HSTS Preload List](https://hstspreload.org/)
* [Rails Secure TLS](https://api.rubyonrails.org/classes/ActionDispatch/SSL.html)
* [Django HSTS Configuration](https://docs.djangoproject.com/en/4.0/ref/settings/#std:setting-SECURE_HSTS_INCLUDE_SUBDOMAINS)
 
### Play
 
These are good interactive labs that allow you to play with real HSTS:

* [Portswigger Vulnerable Labs](https://portswigger.net/web-security/all-labs)
* [Google Gruyere](https://google-gruyere.appspot.com/)
* [SSLStrip](https://github.com/moxie0/sslstrip)