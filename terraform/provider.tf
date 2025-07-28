provider "aws" {
  region     = var.aws_region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key

  dynamic "endpoints" {
    for_each = var.endpoint_url != "" ? [var.endpoint_url] : []
    content {
      iam        = endpoints.value
      lambda     = endpoints.value
      dynamodb   = endpoints.value
      apigateway = endpoints.value
      s3         = endpoints.value
      sts        = endpoints.value
      route53    = endpoints.value
    }
  }
}
