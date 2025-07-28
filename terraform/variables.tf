variable "aws_access_key" {
  type      = string
  sensitive = true
}

variable "aws_secret_key" {
  type      = string
  sensitive = true
}

variable "aws_region" {
  type = string
}

variable "endpoint_url" {
  type    = string
  default = ""
}

variable "url_table_name" {
  type    = string
  default = "ddb_url_table"
}

// lambda function variables

variable "create_func_folder" {
  type = string
  default = "create"
}
variable "redirect_func_folder" {
  type = string
  default = "redirect"
}

variable "create_func_name" {
  type = string
  default = "create-func"
}
variable "redirect_func_name" {
  type = string
  default = "redirect-func"
}

// dynamodb variables

variable "ddb_code_keyname" {
  type = string
  default = "HashCode"
}

variable "ddb_url_keyname" {
  type = string
  default = "LongURL"
}

variable "ddb_gsi_keyname" {
  type = string
  default = "LongURLIndex"
}

// api gateway variables
variable "api_gateway_name" {
  type = string
  default = "api_gateway"
}

variable "stage_name" {
  type = string
  default = "local"
}

# Route 53 

variable "zone_name" {
  type = string
  default = "micro.url"
}

variable "record_name" {
  type = string
  default = "micro.url"
}

variable "certificate_arn" {
  type = string
}