variable "api_name" {
  description = "Friendly name for the API Gateway REST API"
  type        = string
}

variable "stage_name" {
  description = "Name of the stage to deploy (e.g. dev, prod)"
  type        = string
  default     = "dev"
}

variable "create_lambda_arn" {
  description = "ARN of the Lambda for POST /urls"
  type        = string
}

variable "redirect_lambda_arn" {
  description = "ARN of the Lambda for GET /urls/{code}"
  type        = string
}

variable "tags" {
  description = "Tags to apply to all API Gateway resources"
  type        = map(string)
  default     = {}
}
