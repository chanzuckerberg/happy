# Configuration

To work properly, Happy CLI needs access to several configuration values and secrets. These can typically be read from several sources.


## Precedence
Since we accept configuration from various sources, we will need a predefined composition order. From strongest biding to weakest:

1. CLI flag
2. Enfironment Variable
3. Files
4. Default values, where applicable

## Values
### Happy Project Root
A path to the root of your happy project. This is defined by the presence of a `.happy` directory. Typically as `<repo_root>/.happy`.

default: We will attempt to recurse up your current working directory until we find a `.happy` directory. If we find a match, we will use it as default. If we don't then you must specify it yourself. We recommend placing your `.happy` directory at the root of your project.

You can set it with:
flag: `--project-root`
env: `HAPPY_PROJECT_ROOT`

### Happy Config Path
A path to your happy project's configuration. Typically at `<happy_project_root>/.happy/config.json`.

default: If we were able to determine `happy_project_root` we default to `<happy_project_root>/.happy/config.json`.

You can set it with:
flag: `--config-path`
env: `HAPPY_CONFIG_PATH`

### Docker Compose Config Path
A path to your project's Docker Compose file. Typically sits next to your `.happy` directory as `<happy_project_root>/docker-compose.yml`.

default: If we were able to determine `happy_project_root` we default to `<happy_project_root>/docker-compose.yml`.

You can set it with:
flag: `--docker-compose-config-path`
env: `DOCKER_COMPOSE_CONFIG_PATH`

### Env
The happy env you want to interact with.

default: `renv`

You can set it with:
flag: `--env`
env: `HAPPY_ENV`
