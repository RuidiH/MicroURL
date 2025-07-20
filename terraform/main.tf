module "url_table" {
    source = "./modules/dynamodb"
}

module "create_fn" {
    source  = "./modules/lambda"
}

module "lookup_fn" {
    source  = "./modules/lambda"
}