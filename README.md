# URL Shortener

## Overview

This is a service that will shorten long URLs into easily shareable URLs. This service is built in Go, using MongoDB and Redis as data stores.

## Design

### Capacity Estimates

We estimate there will be about 2000 daily active users accessing this service. Out of this, there is assumed to be a 4:1 split between write and read operations -- i.e. we assume there are ~400 new URLs being generated daily and ~1600 reads on generated URLs.


## Data model

This is write-heavy system. The data itself is also not strongly relational, as we are not storing many relationships between entities. Therefore, we can go with a NoSQL data model, which also has the added benefit of easy horizontal scalability, should we need it. For the actual shortened identifier, we will use a 4 character alphanumeric combination -- this allows for 62<sup>4</sup> = 14776336 unique combinations, which should reduce the likelihood of collisions, while also generating a short, shareable URL.

## API

We will create a REST API that will allow for URL manipulations, displaying analytics, and general health checks. This is developed via the [Gin framework](https://gin-gonic.com/docs/). These API's are grouped by the particular domain:
* URL Manipulation (CRUD operations for URL generation / redirection):
```json

Create a new shortened URL
POST /api/v1/urls

Request:
{
    "redirectUrl": <URL to shorten>,
    "expiryDate": <optional expiry date for this short URL>
}

Response:
{
    "shortenedId": <4 character alphanumeric code>,
    "redirectUrl": ...,
    "expiryDate": ...,
    "createdAt": ...
}

Get URL
GET /api/v1/urls/<:id>

Response:
302 with Location header

Delete URL
DELETE /api/v1/urls/<:id>

Response:
204 with no content
```

* URL analytics:
```json

GET /api/v1/analytics/urls/<:id>

Response:
{
    "shortenedId": ...,
    "count": ...
}
```

* Health check
```json

GET /api/ping

Response
{
    "message": "pong"
}
```

## Running

This whole project has been containerized. There is a helper script called `run.sh`. One can run the application via
```sh
./run.sh -r
```

This will then launch the MongoDB and Redis containers (these containers will NOT be exposed outside of the internal Docker network), after which it will launch the API container on port 8000.

## Testing

The previously mentioned helper script can also invoke the test suite (this has also been containerized):
```sh
./run -t
```