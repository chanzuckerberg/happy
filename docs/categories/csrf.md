---
parent: Rule Categories
title: CSRF
layout: default
has_toc: true
---
 
## Cross-site Request Forgery (CSRF)
 
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
 
Attackers trick a user to visit their site, where the attacker can make changes to other sites that user might be logged into.
 
### Scenario
 
Jane is running the website `https://janesite.com`. Jack logs into Jane's site with his email. Bill is trying to exploit a CSRF issue in `https://janesite.com`, allowing him to change another user's email. Bill hosts a malicious website at `https://billsite.com` and asks Jack to check it out. Whenever Jack visits Bill's web site with his web browser, Bill's site uses JavaScript to make a request to Jane's site from Jack's web browser. Bill changes Jack's email on `https://janesite.com`.
 
### Reasoning
 
This issue happens because of the history of the web. When the web was designed, it was primarily text pages with hyperlinks to other pages. These pages and links were connected together to create a "web". For sensitive resources, some pages needed to be authenticated. To make a smooth user experience, when a user clicked a link, the browser would automatically send the authentication cookies it had remembered for **reading** that page. This way, the user would not need to reauthenticate on every link click.
 
The issue started to arise when the web evolved to something other than **reading** text and links. As the web stated to become a platform for applications, web pages needed to **write** as well as **read** text. For a logged in user, a link might unsubscribe them from an email list, change their login email address, or even make a bank transfer. In these examples, navigating to that resource changes the state of the application. Another way of thinking of this is in the original web, everything was **read-only**. Today, we can perform **writes** as well as **reads** and CSRF targets these **writes**.
 
When a user visits a malicious web page attempting to abuse a CSRF issue, the attacker is using the fact that the browser will submit the user's remembered credentials with any request. The attacker's website sends a request from the user's browser to an authenticated resource that **writes** a resource. Since the credentials for the resource will be submitted with the request, the request goes through.
 
#### Example
 
For example, an attacker might make the following malicious website:
 
~~~html
<DOCTYPE html>
<html>
   <head>
       <script>
           fetch("https://yourbank.com/transfer", {
               "method": "POST",
               "body": JSON.stringify({
                   "to_acount": "attack_account_111",
                   "from_account": "victim_account_000",
                   "amount": "$1000000",
               }),
               "mode": "no-cors",
               "credentials": "include"
           })
       </script>
   </head>
   <body>
       <h1>Gotcha!</h1>
   </body>
</html>
~~~
 
This website makes a HTTP POST to `https://yourbank.com/transfer`. If Jack navigates to this site, that request will come from Jack's browser. If Jack is logged into `https://yourbank.com`, Jack's browser will send any authentication cookies to `https://yourbank.com`. `https://yourbank.com` will think Jack wanted to transfer money to the attacker's account.   
 
### Solution
 
#### Verify the Origin of All Stateful Requests
 
The classic way of doing this is using a [CSRF token](https://portswigger.net/web-security/csrf/tokens). This is a value that only your browser and the server know. The server verifies that its users all submit their token when making stateful requests. A malicious webpage won't know this value, so won't be able to submit it.
 
Another way of doing this is called [Double Submit](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#double-submit-cookie). This is where the CSRF token **is** the authentication cookie. The authentication cookie submitted both as a cookie value and as an argument to the stateful request. The server verifies that both are the same before allowing the request to go through. This is an efficient method because it doesn't require a separate secret value to be stored in a database. CSRF tokens can be verifed in-memory instead of with a database lookup.
 
#### Limit the Power of Cookies
 
Cookies now come with the [SameSite](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#samesite-cookie-attribute) attribute. This allows website designers to instruct the browser if other domains are allowed to send cookies for your application. Use the SameSite attribute on all authentication cookies with the value of “Strict” to prevent this cookie from being sent by other sites. The “Strict” setting will prevent a CSRF issue in all cases. A “Lax” option is also available where cookies are sent from other sites when a user is navigating, but not sent when submitting a form or using JavaScript. This might be an acceptable option, but still leaves applications potentially vulnerable to CSRF if state-changing requests are initiated through a link click.
 
### Fun Facts
 
* If CSRF is interesting to you, read [Weaving the Web](https://www.amazon.com/Weaving-Web-Original-Ultimate-Destiny/dp/006251587X) by Tim Burners Lee (inventor of the web).
* Most modern web frameworks come with CSRF protections out of the box, but need to be specifically enabled.
* The default behavior of cookies without SameSite for modern browsers, [like Chrome, is "Lax"](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite#cookies_without_samesite_default_to_samesitelax).
 
## References
 
### Read
 
These are good resources to read more advanced docs about this issue:
 
* [Advanced CSRF Definition and Mitigation](https://portswigger.net/web-security/csrf)
* [CSRF Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
* [OWASP CSRF](https://owasp.org/www-community/attacks/csrf)
* [Rails CSRF Defenses](https://guides.rubyonrails.org/security.html#cross-site-request-forgery-csrf)
* [Django CSRF Defenses](https://docs.djangoproject.com/en/4.0/ref/csrf/)
* [State Changing Requests](https://www.cloudflare.com/learning/security/threats/cross-site-request-forgery/)
 
### Play
 
These are good interactive labs that allow you to play with real CSRF:
* [Portswigger Vulnerable Labs](https://portswigger.net/web-security/all-labs)
* [Google Gruyere](https://google-gruyere.appspot.com/)