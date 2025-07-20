variable "name" {
    description = "Lamda resources base name"
    type = string
}

variable "filename" {
    description = "Zip path"
    type = string
}

variable "handler" {
    description = "Lambda hanlder"
    type = string
    default = null
}

variable "runtime" {
    description = "Lambda runtime"
    type = string
    default = "provided.al2023"
}

variable "role_arn" {
    description = "IAM role ARN that the lambda assumes"
    type = string
}

variable "environment_variables" {  
    description = "Environment variables Map"
    type = map(string)
    default = {}
}

variable "memory_size" {
    description = "Lambda memory in MB"
    type = number
    default = 128
}

variable "timeout" {
    description = "Lambda timeout in seconds"
    type = number
    default = 10
}