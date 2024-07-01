### Prepare
# Create user
POST http://{{host}}/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","username":"{{username}}","password":"{{password}}"}
HTTP 201
[Captures]
user_id: jsonpath "$['id']"

# Login
POST http://{{host}}/v1/auth/login
Content-Type: application/json; charset=utf-8
{"username":"{{username}}","password":"{{password}}"}
HTTP 200
[Captures]
token: jsonpath "$['token']"

### Tests
# Create Ingredient 1
POST http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Beef"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Beef"
jsonpath "$.parent_id" == null
[Captures]
id1: jsonpath "$['id']"

# Create Ingredient 2
POST http://{{host}}/v1/ingredients
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
POST http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" == null
[Captures]
id3: jsonpath "$['id']"

# Create Duplicate of Ingredient 3
POST http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Ingredient - parent_id doesn't exists
POST http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Pork", "parent_id": "c624bce3-2d1b-4ae8-87e2-af775be70077"}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Ingredient 4
POST http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Vegetables"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Vegetables"
jsonpath "$.parent_id" == null
[Captures]
id4: jsonpath "$['id']"

# Update Ingredient 3
PUT http://{{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus", "parent_id":"{{id4}}"}
HTTP 200
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" == {{id4}}

# Update Ingredient - parent_id doesn't exists
PUT http://{{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Pork", "parent_id": "c624bce3-2d1b-4ae8-87e2-af775be70077"}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# List Ingredients
GET http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 4

# Delete Ingredient 2
DELETE http://{{host}}/v1/ingredients/{{id2}}
Authorization: Bearer {{token}}
HTTP 204

GET http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 3

# Delete Ingredient 4: should fail because it's a parent
# DELETE http://{{host}}/v1/ingredients/{{id4}}
# Authorization: Bearer {{token}}
# HTTP 500

GET http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 3


### Clean up

# Delete Ingredient 3
DELETE http://{{host}}/v1/ingredients/{{id3}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 1
DELETE http://{{host}}/v1/ingredients/{{id1}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 4
DELETE http://{{host}}/v1/ingredients/{{id4}}
Authorization: Bearer {{token}}
HTTP 204

GET http://{{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" == null

# ForgetMe (while logged in)
DELETE http://{{host}}/v1/users
Authorization: Bearer {{token}}
HTTP 204
