output "happy_route53_zone_name" {
  value = local.hosted_zone_name
}

output "happy_route53_zone_zone_id" {
  value = aws_route53_zone.happy_route53_zone.zone_id
}
