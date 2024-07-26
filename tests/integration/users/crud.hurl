# Create user
POST {{host}}/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","email":"{{email}}","password":"{{password}}"}
HTTP 201
[Captures]
user_id: jsonpath "$['id']"

# Create user - duplicate email
POST {{host}}/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","email":"{{email}}","password":"{{password}}"}
HTTP 403

# Create user - wrong email format
POST {{host}}/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","email":"badEmail","password":"{{password}}"}
HTTP 400

# Login
POST {{host}}/v1/auth/login
Content-Type: application/json; charset=utf-8
{"email":"{{email}}","password":"{{password}}"}
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
  "email":"{{email}}",
  "password": "{{password2}}"
}
HTTP 200

# Update user - duplicate email: expect nothing to happen, email is not considered
PUT {{host}}/v1/users
Authorization: Bearer {{token}}
Content-Type: text/plain; charset=utf-8
[Options]
variable: password2=verySafePassword101
{
  "name": "julia2",
  "email":"{{email}}",
  "password": "{{password2}}"
}
HTTP 200

# Update user - wrong email format: expect nothing to happen, email is not considered
PUT {{host}}/v1/users
Authorization: Bearer {{token}}
Content-Type: text/plain; charset=utf-8
[Options]
variable: password2=verySafePassword101
{
  "name": "julia2",
  "email":"badEmail",
  "password": "{{password2}}"
}
HTTP 200

# Login post update
POST {{host}}/v1/auth/login
Content-Type: application/json; charset=utf-8
{"email":"{{email}}","password":"{{password2}}"}
HTTP 200
[Captures]
token2: jsonpath "$['token']"
refresh_token2: jsonpath "$['refresh_token']"

# ForgetMe (while logged in)
DELETE {{host}}/v1/users
Authorization: Bearer {{token2}}
HTTP 204

