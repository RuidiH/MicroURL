data "aws_region" "current" {}

# 1. The REST API
resource "aws_api_gateway_rest_api" "this" {
  name        = var.api_name
  description = "URL shortener API"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
  tags = var.tags
}

# 2. /urls
resource "aws_api_gateway_resource" "urls" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = "urls"
}

# 3. /urls/{code}
resource "aws_api_gateway_resource" "code" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_resource.urls.id
  path_part   = "{code}"
}

# 4. POST /urls → create_lambda
resource "aws_api_gateway_method" "create" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.urls.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "create" {
  rest_api_id             = aws_api_gateway_rest_api.this.id
  resource_id             = aws_api_gateway_method.create.resource_id
  http_method             = aws_api_gateway_method.create.http_method
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.region}:lambda:path/2015-03-31/functions/${var.create_lambda_arn}/invocations"
}

resource "aws_lambda_permission" "allow_create" {
  statement_id  = "AllowInvokeCreate"
  action        = "lambda:InvokeFunction"
  function_name = var.create_lambda_arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.this.execution_arn}/*/POST/urls"
}

# 5. GET /urls/{code} → lookup_lambda
resource "aws_api_gateway_method" "redirect" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.code.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "redirect" {
  rest_api_id             = aws_api_gateway_rest_api.this.id
  resource_id             = aws_api_gateway_method.redirect.resource_id
  http_method             = aws_api_gateway_method.redirect.http_method
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.region}:lambda:path/2015-03-31/functions/${var.redirect_lambda_arn}/invocations"
}

resource "aws_lambda_permission" "allow_redirect" {
  statement_id  = "AllowInvokeLookup"
  action        = "lambda:InvokeFunction"
  function_name = var.redirect_lambda_arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.this.execution_arn}/*/GET/urls/*"
}

# Create a deployment
resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  # Forces a new deployment on any change to integration resources:
  depends_on = [
    aws_api_gateway_integration.create,
    aws_api_gateway_integration.redirect,
  ]
}

# Publish that deployment to a stage
resource "aws_api_gateway_stage" "this" {
  rest_api_id    = aws_api_gateway_rest_api.this.id
  deployment_id  = aws_api_gateway_deployment.this.id
  stage_name     = var.stage_name
}

resource "aws_api_gateway_domain_name" "custom" {
  domain_name = var.domain_name
  regional_certificate_arn = var.certificate_arn
  endpoint_configuration {
    types = ["REGIONAL"]
  }
  security_policy = "TLS_1_2"
}

resource "aws_api_gateway_base_path_mapping" "default" {
  api_id = aws_api_gateway_rest_api.this.id
  domain_name = aws_api_gateway_domain_name.custom.domain_name
  stage_name  = var.stage_name  
}
