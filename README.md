# MicroURL (<i>Work in Progress</i>)
MicroURL takes a long URL and returns a micro URL for you

## Installation

// TODO add installation instructions 

## Design

### Functional Requirement
1. Users can post a long url and get a short url.
2. Visiting the short url should redirect user to the long url.
3. The short url is unique for every long url provided.
4. (Optional) Users can set a expiration date for the link.
5. (Optional) Users can use custom alias instead if it does not already exist.

### Non-functional Requirement
1. User should receive the short url in 10 seconds.
2. The short url should be avaialble in 1 minute.
3. The service should support 30,000 DAU activities.


### System Desgin Diagram

---

![Image](/design.png)

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