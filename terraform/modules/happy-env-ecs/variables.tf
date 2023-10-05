
variable "tags" {
  type        = map(string)
  description = "tags to associate with env resources"
}

variable "name" {
  type = string
}

variable "base_zone" {
  description = "base route53 zone"
  type        = string
}

variable "cloud-env" {
  type = object({
    public_subnets        = list(string)
    private_subnets       = list(string)
    database_subnets      = list(string)
    database_subnet_group = string
    vpc_id                = string
    vpc_cidr_block        = string
  })
}

variable "public_lb_services" {
  description = "Create a public-facing ALB for these services"
  type        = set(string)
  default     = []
}

variable "private_lb_services" {
  description = "Create a private load balancers for a set of services sitting behind Okta/OAuth proxy"
  type        = set(string)
  default     = []
}

variable "ecr_repos" {
  description = "Map of ECR repositories to create. These should map exactly to the service names of your docker-compose"
  type = map(object({
    name           = string,
    read_arns      = optional(list(string), []),
    write_arns     = optional(list(string), []),
    tag_mutability = optional(bool, true),
    scan_on_push   = optional(bool, false),
  }))
  default = {}
}

variable "rds_dbs" {
  description = "set of DB's to create"
  type        = map(object({ name = string, username = string, instance_class = string }))
  default     = {}
}

variable "s3_buckets" {
  description = "S3 buckets to create"
  type        = map(object({ name = string }))
  default     = {}
}

variable "swipe_module_envs" {
  description = "set of swipe-batch envs to create"
  type = map(object({
    name                   = string,
    ami_id                 = string,
    mock                   = bool,
    job_policy_arns        = list(string),
    spot_min_vcpus         = number,
    spot_max_vcpus         = number,
    on_demand_min_vcpus    = number,
    on_demand_max_vcpus    = number,
    instance_types         = list(string),
    extra_env_vars         = map(string),
    sqs_queues             = map(map(string)),
    workspace_s3_prefix    = string,
    wdl_workflow_s3_prefix = string,
  }))
  default = {}
}

variable "swipe_envs" {
  description = "set of swipe-batch envs to create"
  type        = map(object({ version = string, name = string, job_policy_arns = list(string), min_vcpus = number, max_vcpus = number, spot_desired_vcpus = number, ec2_desired_vcpus = number, instance_type = list(string), ami_id = string }))
  default     = {}
}

variable "batch_envs" {
  description = "set of batch envs to create"
  type = map(object({
    version         = string,
    name            = string,
    job_policy_arns = list(string),
    min_vcpus       = number,
    max_vcpus       = number,
    desired_vcpus   = number,
    instance_type   = list(string),
    init_script     = optional(string),
  volume_size = number }))
  default = {}
}

variable "services" {
  description = "set of services to prebuild LB's for"
  type        = map(object({ idle_timeout = number }))
  default     = {}
}

variable "min_servers" {
  type        = number
  description = "Minimum number of instances for the cluster"
  default     = 2
}

variable "max_servers" {
  type        = number
  description = "Maximum number of instances for the cluster. Must be at least var.min_servers + 1."
  default     = 5
}

variable "instance_type" {
  type    = string
  default = "m5.large"
}

variable "datadog_api_key" {
  type        = string
  default     = ""
  description = "A datadog api key to enable the datadog agent on the instance"
}

// TODO: remove this variable. It is not really used and confusing as a required input parameter
variable "ssh_key_name" {
  type        = string
  default     = "happy_key"
  description = "Deprecated"
}

variable "roll_interval_hours" {
  type        = number
  default     = 8
  description = "how often to roll hosts"
}

variable "app_ports" {
  type        = set(number)
  default     = [80, 8080, 8000, 5000, 9000]
  description = "What ports do tasks need to be able to reach each other on?"
}

variable "additional_secrets" {
  type        = any
  default     = {}
  description = "Any extra secret key/value pairs to make available to services"
}

variable "oauth_dns_prefix" {
  type        = string
  default     = ""
  description = "DNS prefix for oauth-proxied stacks. Leave this empty if we don't need a prefix!"
}

variable "oauth_bypass_paths" {
  type        = set(string)
  default     = []
  description = "Bypass these paths in the oauth proxy"
}

variable "ssh_users" {
  description = "Okta groups that should have SSH access to ECS nodes"
  type        = list(object({ username : string, sudo_enabled : bool }))
  default     = []
}

variable "extra_security_groups" {
  description = "Security groups that need access to RDS DB's"
  type        = list(string)
  default     = []
}

variable "db_engine_version" {
  description = "The Aurora Postgres engine version"
  type        = string
  default     = "14.3"
}

variable "extra_proxy_args" {
  description = "Add to the proxy's default arguments."
  type        = set(string)
  default     = []
}

variable "github_actions_roles" {
  description = "Roles to be used by Github Actions to perform Happy CI."
  type = set(object({
    name = string
    arn  = string
  }))
  default = []
}

variable "hapi_base_url" {
  type        = string
  description = "The base URL for HAPI"
  default     = "https://hapi.hapi.prod.si.czi.technology"
}

variable "regional_wafv2_arn" {
  type        = string
  description = "A WAF to protect the happy env if needed"
  default     = null
}

variable "oauth2_proxy_registry_image" {
  type    = string
  default = "626314663667.dkr.ecr.us-west-2.amazonaws.com/oauth2_proxy"
}

variable "oauth2_proxy_image_version" {
  type    = string
  default = "sha-0c26cccf"
}

variable "lifecycle_policy" {
  description = "ECR ifecycle policy, as JSON string. If provided, max_images is ignored."
  type        = string
  default     = ""
}
