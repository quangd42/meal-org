### Prepare
# Create user
POST {{host}}/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","username":"{{username}}","password":"{{password}}"}
HTTP 201
[Captures]
user_id: jsonpath "$['id']"

# Login
POST {{host}}/v1/auth/login
Content-Type: application/json; charset=utf-8
{"username":"{{username}}","password":"{{password}}"}
HTTP 200
[Captures]
token: jsonpath "$['token']"

### Tests
# List Cuisines
GET {{host}}/v1/cuisines
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count > 0

### Clean up
# ForgetMe (while logged in)
DELETE {{host}}/v1/users
Authorization: Bearer {{token}}
HTTP 204

