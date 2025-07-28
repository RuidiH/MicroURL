output "zone_id" {
  description = "ID of the hosted zone"
  value       = aws_route53_zone.this.zone_id
}

output "record_fqdn" {
  description = "Fully qualified domain name of the record"
  value       = aws_route53_record.this.fqdn
}
