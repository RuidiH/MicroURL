variable "name" {
  description = "Dynamodb table name"
  type        = string
  default     = "table-name"
}

variable "code_keyname" {
  description = "Key name of the table"
  type = string
  default = "HashCode"
}

variable "url_keyname" {
  description = "Key name of the user url"
  type = string
  default = "LongURL"
}

variable "gsi_keyname" {
  description = "Key name of the Global Secondary Index"  
  type = string
  default = "LongURLIndex"
}