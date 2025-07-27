module "url_table" {
  source = "./modules/dynamodb"
  name   = var.url_table_name
  code_keyname = var.ddb_code_keyname
  url_keyname = var.ddb_url_keyname
  gsi_keyname = var.ddb_gsi_keyname
}

resource "aws_iam_role" "lambda_exec" {
  name               = "url-shortener-lambda-exec"
  assume_role_policy = <<EOF
    {
    "Version":"2012-10-17",
    "Statement":[
        {
        "Action":"sts:AssumeRole",
        "Principal":{"Service":"lambda.amazonaws.com"},
        "Effect":"Allow"
        }
    ]
    }
    EOF
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

module "create_micro" {
  source   = "./modules/lambda"
  name     = var.create_func_name
  filename = "${path.module}/../src/${var.create_func_folder}/bin/function.zip"
  handler  = "create"
  role_arn = aws_iam_role.lambda_exec.arn  

  environment_variables = {
    TABLE_NAME = var.url_table_name
    AWS_REGION = var.aws_region
    BASE_URL = var.base_url
    CODE_KEYNAME = var.ddb_code_keyname
    URL_KEYNAME = var.ddb_url_keyname
    GSI_KEYNAME = var.ddb_gsi_keyname
  }
}

module "redirect_micro" {
  source = "./modules/lambda"
  name = var.redirect_func_name
  filename = "${path.module}/../src/${var.redirect_func_folder}/bin/function.zip"
  handler  = "redirect"
  role_arn = aws_iam_role.lambda_exec.arn

  environment_variables = {
    TABLE_NAME = var.url_table_name
    AWS_REGION = var.aws_region
    BASE_URL = var.base_url
    CODE_KEYNAME = var.ddb_code_keyname
    URL_KEYNAME = var.ddb_url_keyname
    GSI_KEYNAME = var.ddb_gsi_keyname
  }
}

module "api_gateway" {
  source = "./modules/apigateway"
  api_name = var.api_gateway_name
  stage_name = var.stage_name
  create_lambda_arn = module.create_micro.function_arn
  redirect_lambda_arn = module.redirect_micro.function_arn
}
