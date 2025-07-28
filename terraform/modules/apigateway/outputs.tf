output "rest_api_id" {
  description = "ID of the REST API"
  value       = aws_api_gateway_rest_api.this.id
}

output "invoke_url" {
  description = "Full public URL including stage"
  value       = "https://${aws_api_gateway_rest_api.this.id}.execute-api.${data.aws_region.current.region}.amazonaws.com/${var.stage_name}"
}

output "regional_domain_name" {
  description = "API Gateway's regional endpoint" 
  value = aws_api_gateway_domain_name.custom.regional_domain_name
}

output "regional_zone_id" {
  description = "API Gateway's zone ID" 
  value = aws_api_gateway_domain_name.custom.regional_zone_id
}
