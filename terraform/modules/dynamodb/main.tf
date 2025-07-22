resource "aws_dynamodb_table" "this" {
  name         = var.name
  billing_mode = "PAY_PER_REQUEST"
  # read_capacity = 10
  # write_capacity = 5

  hash_key = "HashCode"
  attribute {
    name = "HashCode"
    type = "S"
  }

  attribute {
    name = "LongURL"
    type = "S"
  }

  global_secondary_index {
    name            = "LongURLIndex"
    hash_key        = "LongURL"
    projection_type = "ALL"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = {
    Name = "URLTable"
  }
}