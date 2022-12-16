output "ecs" {
  value = module.ecs-cluster.ecs
}

output "public-lbs" {
  value = { for service, alb in aws_lb.lb-public : service => {
    arn              = alb.arn
    dns_name         = alb.dns_name
    zone_id          = alb.zone_id
    route53_dns_name = aws_route53_record.services[service].name
  } }
}

output "security_groups" {
  value = [aws_security_group.happy_env_sg.id]
}

# Output the same interface we provide to the CLI.
output "integration_secret" {
  value = local.secret_string
}

output "github_ci_role_arns" {
  value = module.happy_github_ci_role
}
