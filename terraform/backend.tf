terraform {
    backend "s3" {
        bucket = "tf-state-bucket"
        key = "/tf-state-key"
        region = "us-west-2"
        # dynamodb_table = ""
    }
}