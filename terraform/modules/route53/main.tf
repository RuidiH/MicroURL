resource "aws_route53_zone" "this" {
  name = var.zone_name
  tags = var.tags
}

resource "aws_route53_record" "this" {
  zone_id = aws_route53_zone.this.zone_id
  name    = var.record_name
  type    = var.record_type

  alias {
    name  = var.regional_domain_name
    zone_id = var.regional_zone_id
    evaluate_target_health = false
  }
}
