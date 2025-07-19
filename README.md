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
