output "kms_key_id" {
  value = aws_kms_key.service_user.key_id
}

output "oidc_config" {
  value = module.service_user.oidc_config
}
