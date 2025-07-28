variable "zone_name" {
  description = "The domain name for the hosted zone (e.g. short.local)"
  type        = string
}

variable "record_name" {
  description = "The FQDN of the record (e.g. short.local or api.short.local)"
  type        = string
}

variable "record_type" {
  description = "DNS record type"
  type        = string
  default     = "A"
}

variable "tags" {
  description = "Tags to apply"
  type        = map(string)
  default     = {}
}

variable "regional_domain_name" {
 description = "API Gateway's regional domain name" 
 type = string
}

variable "regional_zone_id" {
  description = "API Gateway's regional zone id"
  type = string
}