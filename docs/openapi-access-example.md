# OpenAPI访问示例

- 入口地址：`<云效主域名>/oapi/v1`
- 请求示例

```bash
curl -X 'GET' \
  'https://<your-host>/oapi/v1/platform/users?status=enabled' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'x-yunxiao-token: <your-token>'
```

- 返回结构

```json
[
  {
    "email": "string",
    "id": "string",
    "joinDate": "string",
    "name": "string",
    "nickName": "string",
    "staffId": "string",
    "sysDeptIds": [
      "string"
    ],
    "username": "string"
  }
]
```
