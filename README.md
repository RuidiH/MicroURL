# MicroURL

## Functional Requirement

## Non-functional Requirement

## System Desgin

---

![Image](/sd_v1.0.png)

---

### API
POST /urls
{
   long_url : "https://example.com/very/long/path?with=params"
}

response
201 created
{ 
    "code": "abc123", 
    "short_url": "https://short.ly/abc123" 
}

GET /urls/{code}

response
302 found
Location: "https://example.com/very/long/path?with=params"


### NOTE

create lambda function: 
awslocal lambda create-function --function-name url-shortener --runtime provided.al2023 --zip-file fileb://function.zip --handler bootstrap --role arn:aws:iam::000000000000:role/lambda-role