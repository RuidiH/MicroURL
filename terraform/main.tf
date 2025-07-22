module "url_table" {
  source = "./modules/dynamodb"
  name   = var.url_table_name
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
  name     = "create-micro-url"
  filename = "${path.module}/../src/create_url/bin/create.zip"
  handler  = "create_micro"
  role_arn = aws_iam_role.lambda_exec.arn

  environment_variables = {
    TABLE_NAME = var.url_table_name
  }
}

# module "lookup_micro" {
#   source = "./modules/lambda"
# }
