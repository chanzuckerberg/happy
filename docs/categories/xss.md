---
parent: Rule Categories
title: XSS
layout: default
has_toc: true
---
 
## Cross-Site Scripting (XSS)
 
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
 
Attackers can write JavaScript that executes in users' browsers loading your application. This means when a user visits the site, the attacker can execute code that could control the user's session, read their data, perform phishing, and more all from within the trust of your application.
 
### Scenario
 
Jane is running the website `https://janesite.com`. Bill exploits an XSS issue in `https://janesite.com`, allowing him to insert a JavaScript payload into the page `https://janesite.com/about`. The payload reads any user's cookies and sends them to `https://billsite.com`. Whenever Jack visits Jane's about page at `https://janesite.com/about`, Bill's JavaScript payload executes in Jack's browser. 
 
### Reasoning
 
This issue happens because HTML and JavaScript are a part of the same data that is sent from the server to the browser. HTML is content that displays the viewable web page and JavaScript is the code that controls the page, but they are both delivered as part of the same HTML document. Since the browser cannot easily separate these two, valid HTML locations can accidently slip into valid JavaScript locations, making them executable as JavaScript code instead of displayed as HTML. In the case of XSS, attackers do this on purpose to move content that was supposed to be rendered as viewable HTML to be executed as JavaScript.
 
#### Example
 
For example, suppose a forum application allows you to type text in a text field and submit it. Then, the text shows as the latest post at the bottom of the forum thread. In this case, the forum post is rendered as viewable HTML. However, an XSS might happen if we can trick the browser into rendering our forum post as executable JavaScript instead of HTML.
 
Suppose we insert the following text as a forum post:
 
~~~
I think it's not great AT ALL
~~~
 
Let's look at the code for what the server might return:
 
~~~html
<DOCTYPE html>
<html>
   <body>
       <h1>Forum Posts</h1>
       <ul>
           <li>This forum is so cool!</li>
           <li>Ya totally</li>
           <li>I think it's not great AT ALL</li>
       </ul>
       <textarea></textarea>
   </body>
</html>
~~~
 
Our post is added to the bottom of the list as a child of the `<li>` element. Adding HTML to our forum post creates a different situation:
 
~~~
I think it's not great AT ALL<script>console.log('I am executing scripts now')</script>"
~~~
 
The server might respond with:
 
~~~html
<DOCTYPE html>
<html>
   <body>
       <h1>Forum Posts</h1>
       <ul>
           <li>This forum is so cool!</li>
           <li>Ya totally</li>
           <li>I think it's not great AT ALL<script>console.log('I am executing scripts now')</script></li>
       </ul>
       <textarea></textarea>
   </body>
</html>
~~~
 
Because the `<script>` tag tells the browser the child text is executable JavaScript, it will be executed, even though it was controlled by whoever wrote the forum post. Worse, this behavior will happen to any use that views these forum posts. A `<script>` tag is one way to do this, but there are many other methods of getting an XSS payload to execute.   
 
### Solution
 
#### User-controlled content is never rendered as executable JavaScript.
 
The de facto way of doing this is [HTML output encoding](https://portswigger.net/web-security/cross-site-scripting/preventing). In the example above, the HTML encoded string of
 
~~~
I think it's not great AT ALL<script>console.log('I am executing scripts now')</script>
~~~
 
becomes
 
~~~
I think its not great AT ALL&#x3c;script&#x3e;console.log('I am executing scripts now')&#x3c;/script&#x3e;
~~~
 
The browser will convert `&#x3e;` and `&#x3c;` to the `>` and `<` symbol character, respectively, and will keep all the children's text as renderable HTML instead of executable JavaScript.
 
Read more below about alternative solutions when output encoding is not viable for your application.
 
### Fun Facts
 
* JavaScript can be used to do anything the user would normally do such as click buttons, fill out forms, submit forms, and navigate pages.
* An XSS JavaScript payload will only execute while the victim's browser tab is open. The payload cannot do anything while the browser program is closed.
* The above example illustrates a "stored" XSS, but there are other variants of this issue such as "reflected" and "DOM".
* JavaScript is executed immediately for elements, even before they are written in the web page. For example, JavaScript that looks like `const a = new Image(); a.src="http://example.com"; a.onload="alert('loaded')"` will load that image and execute immediately rather than waiting for the image to be added to the web page.
 
## References
 
### Read
 
These are good resources to read more advanced docs about this issue:
 
* [XSS Cheatsheet](https://portswigger.net/web-security/cross-site-scripting/cheat-sheet)
* [Advanced XSS Definitions](https://portswigger.net/web-security/cross-site-scripting)
* [Stored XSS](https://portswigger.net/web-security/cross-site-scripting/stored)
* [Reflected XSS](https://portswigger.net/web-security/cross-site-scripting/reflected)
* [DOM XSS](https://portswigger.net/web-security/cross-site-scripting/dom-based)
 
### Play
 
These are good interactive labs that allow you to play with real XSS:
* [Portswigger Vulnerable Labs](https://portswigger.net/web-security/all-labs)
* [Google Gruyere](https://google-gruyere.appspot.com/)
* [XSS Firing Range](https://public-firing-range.appspot.com/)
* [XSS Game](https://xss-game.appspot.com/)