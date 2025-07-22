output "function_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.this.arn
}