resource "aws_dynamodb_table" "this" {
  name         = var.name
  billing_mode = "PAY_PER_REQUEST"
  hash_key = var.code_keyname
  # read_capacity = 10
  # write_capacity = 5

  attribute {
    name = var.code_keyname
    type = "S"
  }

  attribute {
    name = var.url_keyname
    type = "S"
  }

  global_secondary_index {
    name            = var.gsi_keyname
    hash_key        = var.url_keyname
    projection_type = "ALL"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = {
    Name = var.name
  }
}