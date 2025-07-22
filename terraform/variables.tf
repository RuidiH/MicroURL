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
  default = "url_table"
}