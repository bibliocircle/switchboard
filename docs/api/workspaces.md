## Listing Mock Services

**Request**

```txt
GET /configuration/mock-services
```

**Response**

```
200 OK
```

```json
{
  "data": [
    {
      "id": "e8c5772e-686f-41de-aed5-2dc0dc3f5f65",
      "type": "REST",
      "title": "Orders mock service",
      "createdAt": "2022-05-23T19:58:39.814Z",
      "updatedAt": "2022-05-23T19:58:39.814Z"
    },
    {
      "id": "aca7751a-b529-4aa3-9f4f-d462164b4ada",
      "type": "REST",
      "title": "Auth mock service",
      "createdAt": "2022-05-23T19:58:39.814Z",
      "updatedAt": "2022-05-23T19:58:39.814Z"
    }
  ]
}
```

## Retrieving Mock Service Configuration

**Request**

```txt
GET /configuration/mock-service/:mockServiceId
```

**Response**

```
200 OK
```

```json
{
  "data": {
    "id": "e8c5772e-686f-41de-aed5-2dc0dc3f5f65",
    "type": "REST",
    "title": "Orders mock service",
    "config": {
      "cors": {
        "allowedOrigins": ["test.example.com"],
        "allowedMethods": ["GET", "POST", "DELETE", "PUT"],
        "allowedHeaders": []
      },
      "defaultResponseHeaders": {
        "x-mocked-response": true
      },
      "upstreams": [
        {
          "name": "Orders API Dev",
          "url": "https://dev.example.com/orders"
        },
        {
          "name": "Orders API Staging",
          "url": "https://staging.example.com/orders"
        }
      ]
    },
    "createdAt": "2022-05-23T19:58:39.814Z",
    "updatedAt": "2022-05-23T19:58:39.814Z"
  }
}
```

## Creating a Mock Service

**Request**

```txt
POST /configuration/mock-service

{
  "type": "REST",
  "title": "Orders mock service",
  "config": {
    "cors": {
      "allowedOrigins": ["test.example.com"],
      "allowedMethods": ["GET", "POST", "DELETE", "PUT"],
      "allowedHeaders": []
    },
    "defaultResponseHeaders": {
      "x-mocked-response": true
    },
    "upstreams": [
      {
        "name": "Orders API Dev",
        "url": "https://dev.example.com/orders"
      },
      {
        "name": "Orders API Staging",
        "url": "https://staging.example.com/orders"
      }
    ]
  }
}
```

**Response**

```
201 OK
```

```json
{
  "id": "e8c5772e-686f-41de-aed5-2dc0dc3f5f65"
}
```

## Updating a Mock Service

**Request**

```txt
PUT /configuration/mock-service/:mockServiceId

{
  "data": {
    "type": "REST",
    "title": "Orders mock service",
    "config": {
      "cors": {
        "allowedOrigins": ["test.example.com"],
        "allowedMethods": ["GET", "POST", "DELETE", "PUT"],
        "allowedHeaders": []
      },
      "defaultResponseHeaders": {
        "x-mocked-response": true
      },
      "upstreams": [
        {
          "name": "Orders API Dev",
          "url": "https://dev.example.com/orders"
        },
        {
          "name": "Orders API Staging",
          "url": "https://staging.example.com/orders"
        }
      ]
    }
  }
}
```

**Response**

```
200 OK
```

```json
{
  "data": {
    "id": "e8c5772e-686f-41de-aed5-2dc0dc3f5f65",
    "type": "REST",
    "title": "Orders mock service",
    "config": {
      "cors": {
        "allowedOrigins": ["test.example.com"],
        "allowedMethods": ["GET", "POST", "DELETE", "PUT"],
        "allowedHeaders": []
      },
      "defaultResponseHeaders": {
        "x-mocked-response": true
      },
      "upstreams": [
        {
          "name": "Orders API Dev",
          "url": "https://dev.example.com/orders"
        },
        {
          "name": "Orders API Staging",
          "url": "https://staging.example.com/orders"
        }
      ]
    },
    "createdAt": "2022-05-23T19:58:39.814Z",
    "updatedAt": "2022-05-23T19:58:39.814Z"
  }
}
```

## Deleting a Mock Service

**Request**

```txt
DELETE /configuration/mock-services/:mockServiceId
```

**Response**

```txt
200 Ok
```
