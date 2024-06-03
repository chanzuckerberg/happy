resource "random_password" "db-secret" {
  for_each = var.rds_dbs
  length   = 32
  special  = false
}

module "db" {
  for_each = var.rds_dbs
  source   = "github.com/chanzuckerberg/cztack//aws-aurora-postgres?ref=v0.49.0"

  project = local.project
  env     = local.env
  service = local.service
  owner   = local.owner

  database_name         = each.value["name"]
  database_password     = random_password.db-secret[each.key].result
  database_username     = each.value["username"]
  database_subnet_group = var.cloud-env.database_subnet_group
  engine_version        = var.db_engine_version
  ingress_security_groups = concat([
    module.ecs-cluster.security_group_id,
    aws_security_group.happy_env_sg.id,
    ],
    var.extra_security_groups,
    [for name, batch in module.batch : batch.batch.security_group],
    [for name, batch in module.batch-swipe : batch.batch.security_group],
    [for name, swipe in module.swipe : swipe.compute_environment_security_group_id],
  )
  instance_class             = each.value.instance_class
  instance_count             = 1
  vpc_id                     = var.cloud-env.vpc_id
  ca_cert_identifier         = "rds-ca-rsa2048-g1"
  auto_minor_version_upgrade = false
  db_deletion_protection     = true
}
