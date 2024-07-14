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
jsonpath "$.parent_id" not exists
[Captures]
id1: jsonpath "$['id']"

# Create Ingredient 2
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Ground Beef", "parent_id":"{{id1}}"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Ground Beef"
jsonpath "$.parent_id" == "{{id1}}"
[Captures]
id2: jsonpath "$['id']"

# Create Ingredient 3
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" not exists
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

# Create Ingredient - parent_id doesn't exists
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Pork", "parent_id": "c624bce3-2d1b-4ae8-87e2-af775be70077"}
HTTP 400
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
jsonpath "$.parent_id" not exists
[Captures]
id4: jsonpath "$['id']"

# Update Ingredient 3
PUT {{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus", "parent_id":"{{id4}}"}
HTTP 200
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" == {{id4}}

# Update Ingredient - parent_id doesn't exists
PUT {{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Pork", "parent_id": "c624bce3-2d1b-4ae8-87e2-af775be70077"}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# List Ingredients
GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 4

# Delete Ingredient 2
DELETE {{host}}/v1/ingredients/{{id2}}
Authorization: Bearer {{token}}
HTTP 204

GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 3

# Delete Ingredient 4: should fail because it's a parent
DELETE {{host}}/v1/ingredients/{{id4}}
Authorization: Bearer {{token}}
HTTP 403

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
