# Configuration

To work properly, Happy CLI needs access to several configuration values and secrets. These can typically be read from several sources.


## Precedence
Since we accept configuration from various sources, we will need a predefined composition order. From strongest biding to weakest:

1. CLI flag
2. Enfironment Variable
3. Files
4. Default values, where applicable

## Values

### Happy Config Path
A path to your happy project's configuration. Typically at `<repo_root>/.happy/config.json`

You can set it with:
flag: `--config-path`
env: `HAPPY_CONFIG_PATH`

### Happy Project Root
A path to the root of your happy project.

You can set it with:
flag: `--project-root`
env: `HAPPY_PROJECT_ROOT`

### Docker Compose Config Path
A path to your project's Docker Compose file.

You can set it with:
flag: `--docker-compose-config-path`
env: `DOCKER_COMPOSE_CONFIG_PATH`

### Env
The happy env you want to interact with.

default: `renv`

You can set it with:
flag: `--env`
env: `HAPPY_ENV`
