resource "random_password" "db_secret" {
  for_each = var.rds_dbs
  length   = 32
  special  = false
}

module "dbs" {
  for_each = var.rds_dbs
  source   = "github.com/chanzuckerberg/cztack//aws-aurora-postgres?ref=v0.49.0"

  project = var.tags.project
  env     = var.tags.env
  # some db names have underscores (this is fine)
  # but the service is used as the name of other things like the parameter groups
  # that do not allow underscores.
  service = replace("${var.tags.service}-${each.value["name"]}", "_", "-")
  owner   = var.tags.owner

  database_name              = each.value["name"]
  database_password          = random_password.db_secret[each.key].result
  database_username          = each.value["username"]
  database_subnet_group      = var.cloud-env.database_subnet_group
  engine_version             = coalesce(each.value["engine_version"], var.default_db_engine_version)
  ingress_security_groups    = [var.eks-cluster.worker_security_group]
  instance_class             = each.value.instance_class
  instance_count             = 1
  vpc_id                     = var.cloud-env.vpc_id
  ca_cert_identifier         = "rds-ca-2019"
  auto_minor_version_upgrade = false
  db_deletion_protection     = true
  rds_cluster_parameters     = each.value.rds_cluster_parameters
}
