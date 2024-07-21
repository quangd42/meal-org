### Prepare
# Create user
POST {{host}}/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","email":"{{email}}","password":"{{password}}"}
HTTP 201
[Captures]
user_id: jsonpath "$['id']"

# Login
POST {{host}}/v1/auth/login
Content-Type: application/json; charset=utf-8
{"email":"{{email}}","password":"{{password}}"}
HTTP 200
[Captures]
token: jsonpath "$['token']"

### Tests
# List Ingredients - expect none
GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" == null

# Create Ingredient 1
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Beef"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Beef"
[Captures]
id1: jsonpath "$['id']"

# Create Ingredient 3
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Asparagus"
[Captures]
id3: jsonpath "$['id']"

# Create Duplicate of Ingredient 3
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 403
[Asserts]
jsonpath "$.error" exists

# Create Ingredient 4
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Vegetables"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Vegetables"
[Captures]
id4: jsonpath "$['id']"

# Update Ingredient 3
PUT {{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 200
[Asserts]
jsonpath "$.name" == "Asparagus"

# List Ingredients
GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 3

GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 3

GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 3


### Clean up

# Delete Ingredient 3
DELETE {{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 1
DELETE {{host}}/v1/ingredients/{{id1}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 4
DELETE {{host}}/v1/ingredients/{{id4}}
Authorization: Bearer {{token}}
HTTP 204

GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" == null

# ForgetMe (while logged in)
DELETE {{host}}/v1/users
Authorization: Bearer {{token}}
HTTP 204
