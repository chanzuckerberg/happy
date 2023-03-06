output "name" {
    value = module.ecr.repository_name
}

output "url" {
  value = module.ecr.repository_url
}

output "arn" {
  value = module.ecr.repository_arn
}
