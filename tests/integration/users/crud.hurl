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
refresh_token: jsonpath "$['refresh_token']"

# Update user
PUT {{host}}/v1/users
Authorization: Bearer {{token}}
Content-Type: text/plain; charset=utf-8
[Options]
variable: password2=verySafePassword101
{
  "name": "julia2",
  "username":"{{username}}",
  "password": "{{password2}}"
}
HTTP 200
[Asserts]
jsonpath "$.name" == "julia2"

# Login post update
POST {{host}}/v1/auth/login
Content-Type: application/json; charset=utf-8
{"username":"{{username}}","password":"{{password2}}"}
HTTP 200
[Captures]
token2: jsonpath "$['token']"
refresh_token2: jsonpath "$['refresh_token']"

# ForgetMe (while logged in)
DELETE {{host}}/v1/users
Authorization: Bearer {{token2}}
HTTP 204

