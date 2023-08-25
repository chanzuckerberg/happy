---
parent: Rule Categories
title: Command Injection
layout: default
has_toc: true
---
 
## Command Injection
 
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
 
Attackers can execute arbitrary operating system commands on the server.
 
### Scenario
 
Jane is running the website `https://janesite.com`. Jane doesn't like to set up complicated debugging software. Instead, she creates an endpoint at `https://janesite.com/debug`. This endpoint copies any query parameters and pastes them into a shell and returns the result. Bill navigates to `https://janesite.com/debug?arg=cat+/home/jane/secret_journal.txt` to read Jane's secret journal.
 
### Reasoning
 
An important principle of software engineering is ["Don’t Repeat Yourself (DRY)"](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself). There is also a principle, popular in the Linux community, to have composable, small programs instead of large, encompassing ones. Command injection happens when these principals are taken to their extremes.
 
Most programming languages have APIs to execute operating system commands via a shell. These APIs are usually pretty simple to use, too: provide a string of commands to execute and receive their output. It is common to use these APIs to invoke software that already exists to accomplish a task or implement a feature and also adhere to the DRY principle.
 
However, the shell is very powerful and potentially dangerous. It can execute any program on the computer and can execute more than one program at a time. Most use cases only need to execute one shell command. Additionally, shells come with their own complex programming language with lots of special characters and shortcuts that allow programmers to do many things with limited keystrokes. Without very strict input validation, the allowed values that are supplied to the shell could allow users to execute other commands that were not intended. Even with very rigorous tests, which are not always a priority for developers with tough deadlines, edge cases slip through. Overall, there are lots of checks to securely allow users to execute shell commands and it is common all these checks aren't performed, leading to arbitrary code execution on the underlying operating system.
 
#### Example
 
For example, Bill might write the following server application to zip the files on his computer:
 
~~~golang
package main
 
import (
   "fmt"
   "net/http"
   "os/exec"
   "strings"
)
 
func zipHandler(rw http.ResponseWriter, r *http.Request) {
   filenames, ok := r.URL.Query()["filename"]
   filename := "/tmp/default.zip"
   if ok {
       filename = filenames[0]
   }
 
   // validate the file ends with .zip
   if !strings.HasSuffix(filename, ".zip") {
       return
   }
   myfavcommand := fmt.Sprintf("zip -r %s /home/bill/files >> /logs", filename)
   cmd := exec.Command("bash", "-c", myfavcommand)
   if err := cmd.Run(); err != nil {
       panic(err)
   }
}
 
func main() {
   http.HandleFunc("/compress", zipHandler)
   http.ListenAndServe(":8090", nil)
}
 
~~~
 
Bill is more familiar with using a shell program called `zip` instead of writing Go code to compress files. The above code allows him to provide a `filename` query parameter to change the output of the zip file and run his favorite `zip` program from the shell. Bill also wants to use the [Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)) special character `>>` to append any logs from the `zip` program to his log directory. Bill knows this can be dangerous so he tries to validate the filename before it is used in his shell command. However, an attacker could abuse this code to execute arbitrary commands by navigating to Bill’s website at `https://billsite.com/compress?filename=test.zip+/home/bill/files+&&+cat+/home/bill/secret_journal.txt+cat`. The command to be run now looks like this:
 
~~~bash
zip -r test.zip /home/bill/files && cat /home/bill/secret_journal.txt && cat /home/bill/files
~~~
 
### Solution
 
#### Avoid the Shell
 
Most programming languages will have the capability to build the same functionality as most shell programs. Making use of the programming language instead of an external program gives the developer better control about the inputs to that functionality whereas the shell programs are limited by the parameters exposed by the program.
 
#### Parse the User Input
 
If a shell script is necessary, find an API that will split the shell string based on a space character and only execute the first word. For example, the split words for `cat /etc/passwd` are `cat` and `/etc/passwd`. Programs that split shell commands into words will only execute the first token as a shell program and the rest of the tokens are only supplied as arguments to that program. All other shell symbols and characters are not executed. Luckily, a lot of shell APIs do this by default now.
 
### Fun Facts
 
* The PHP programming language is notorious for these types of issues because its shell API opens a complete Bash program and copies the string into it.
* The Ruby programming language has a built-in semantic for executing shell commands by putting a string inside enclosing backticks (\`cat /etc/passwd\`).
 
## References
 
### Read
 
These are good resources to read more advanced docs about this issue:
 
* [Advanced command injection definition](https://portswigger.net/web-security/os-command-injection)
* [OWASP command injection definition](https://owasp.org/www-community/attacks/Command_Injection)
* [Python Shell API](https://docs.python.org/3/library/subprocess.html#security-considerations)
* [Open3 - Safer Ruby Command Injection Library](https://docs.ruby-lang.org/en/2.0.0/Open3.html)
 
### Play
 
These are good interactive labs that allow you to play with real command injection:
 
* [Portswigger Vulnerable Labs](https://portswigger.net/web-security/all-labs)
* [Google Gruyere](https://google-gruyere.appspot.com/)