output "table_arn" {
  description = "dynamodb table arn"
  value       = aws_dynamodb_table.this.arn
}