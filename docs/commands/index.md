---
layout: default
title: Commands
has_children: true
nav_order: 6
---

# Commands

Happy has a number of commands that you can use to manage your application. You can see a list of all the commands by running `happy -h`. You can also get help for a specific command by running `happy <COMMAND> -h`.

For each command, the `.happy/config.json` file is used to determine which app to run the command against. If you do not specify an environment with the `--env` flag, the default environment is used. You can set the default environment by setting the `default_env` field in your `.happy/config.json` file.


