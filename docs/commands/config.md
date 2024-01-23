---
parent: Commands
layout: default
has_toc: true
---

# App Configs

One of the features of Happy allows you to manage application configuration values (key-value pairs) for your application that are injected as environment variables into your application's containers. This allows you to configure your application without having to rebuild your container image.

See `happy config -h` for usage information.

## Scope & Hierarchy of Configs

App configs are scoped to an app and environment (and sometimes a stack, see below). This uses the value of `app` from your config.json file and the value passed to the `--env` flag (or the default environment if no `--env` flag is passed). This means that you can have different configs for different environments of the same app, and configs for different apps will not conflict with each other (even when the apps are deployed to the same "happy-env" (as defined by the `happy-env-eks` Terraform module)).

App configs are defined hierarchically where configs can be set at the environment level or at the stack level for your app. Every stack inherits the configs of its parent environment and can be overridden at the stack level. For example, if you have an environment called `rdev` and a stack called `foo`, the `foo` stack will inherit all the configs of the `rdev` environment. When a config with the same key is defined in both the environment and the stack, the stack's config value will take precedence.

## V1

The V1 version of app configs is deprecated and will be removed in a future release. Please use V2.

## V2

This is the recommended version of app configs. It is currently in beta and is not enabled by default.

Depending on whether you are starting a new app or have an existing app, the steps to enable V2 are different. See below.

### Enabling V2 for a new app
1. Make sure the version of the `happy-stack-eks` Terraform module that you are using is at least `v4.26.0`
2. Set `features.enable_happy_config_v2` to `true` in your `config.json` file:
```json
{
  ...
  "features": {
    ...
    "enable_happy_config_v2": true
  }
}
```

### Migrating an existing app to V2

Run `happy config migrate-all` to migrate all keys for the given environment and stack. This will need to be done for each environment and each stack. Your configs are now migrated from the DB used by V1 to Kubernetes secrets in the same cluster/namespace as your app / environment. (If you only want to migrate a subset of keys, use `happy config migrate <KEY>` to migrate a single key from V1 to V2)

Once you have migrated all your configs, you will need to upgrade your `happy-stack-eks` Terraform module to at least `v4.26.0`. This will ensure that your stack's containers get the new secrets injected into them.

#### Validating V2 Configs

Once configs are migrated and `happy-stack-eks` version upgraded, you can validate that your configs are being injected into your containers from V2 instead of V1 by setting `skip_config_injection = true` in your `happy-stack-eks` Terraform module invocation. This skips the V1 config injection, allowing you to validate that your containers are getting the correct configs from V2.

#### Finalizing the Migration

Once you have validated that your containers are getting the correct configs from V2, you can remove the `skip_config_injection = true` line from your `happy-stack-eks` Terraform module invocation.

Next, delete your V1 configs from the DB by running `happy config delete <KEY>` for each of them. Alternatively, reach out to @core-infra and we can delete them for you.

Finally, make V2 your default by setting `features.enable_happy_config_v2` to `true` in your `config.json` file.
```json
{
  ...
  "features": {
    ...
    "enable_happy_config_v2": true
  }
}
```

NOTE: Prior to making V2 your default config version (ie: prior to setting `features.enable_happy_config_v2` to `true`), you can manage V2 configs with all the usual `happy config` commands by passing the `--v2` flag. For example, `happy config get --v2 <KEY>` will get the V2 config for the given key.
