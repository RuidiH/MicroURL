# Lambda function
resource "aws_lambda_function" "this" {
    function_name = var.name
    filename = var.filename
    handler = var.handler
    runtime = var.runtime
    role = var.role_arn

    memory_size = var.memory_size
    timeout = var.timeout

    environment {
        variables = var.environment_variables
    }

    source_code_hash = filebase64sha256(var.filename)
}

# Allow invocation from API Gateway
resource "aws_lambda_permission" "apigw" {
    statement_id = "AllowAPIGatewayInvoke"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.this.function_name
    principal = "apigateway.amazonaws.com"
}