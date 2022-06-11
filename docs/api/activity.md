## Retrieve Workspace Activity Log

**Request**

```txt
GET /workspace/:workspaceId/activity
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
      "request": {
        "hostname": "mock.example.com",
        "port": 443,
        "path": "/some/mock-endpoint",
        "method": "POST",
        "headers": {
          "Accept": "application/json",
          "Authorization": "Bearer d9571856-623c-4dc6-bf0e-afa627b998b1"
        },
        "bodyBase64Encoded": "eyAidGVzdCI6IDEyMyB9"
      },
      "response": {
        "statusCode": 200,
        "headers": {
          "Content-Type": "application/json",
          "Cache-Control": "nocache",
          "Server": "Switchboard/1.0.0"
        },
        "bodyBase64Encoded": "eyAidGVzdCI6IDEyMyB9",
        "responseTimeMs": 100
      },
      "metadata": {
        "scenarioId": "63dd809d-8430-445a-98ef-a5e67b85169b",
      },
      "timestamp": "2022-05-23T19:58:39.814Z"
    }
  ]
}
```

## Clear Workspace Activity

**Request**

```txt
DELETE /workspace/:workspaceId/activity
```

**Response**

```
200 OK
```
