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
# List Cuisines - expect existing cuisines
GET {{host}}/v1/cuisines
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 23

# Create Cuisine empty - expect validation error
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":""}
HTTP 400
[Asserts]
jsonpath "$.error.name" exists

# Create Cuisine 1
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Beef"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Beef"
jsonpath "$.parent_id" not exists
[Captures]
id1: jsonpath "$['id']"

# Create Cuisine 2
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Ground Beef", "parent_id":"{{id1}}"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Ground Beef"
jsonpath "$.parent_id" == "{{id1}}"
[Captures]
id2: jsonpath "$['id']"

# Create Cuisine 3
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" not exists
[Captures]
id3: jsonpath "$['id']"

# Create Cuisine - parent_id doesn't exists
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Pork", "parent_id": "c624bce3-2d1b-4ae8-87e2-af775be70077"}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Cuisine - name empty
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":""}
HTTP 400
[Asserts]
jsonpath "$.error.name" exists

# Create Cuisine 4
POST {{host}}/v1/cuisines
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Vegetables"}
HTTP 201
[Asserts]
jsonpath "$.name" == "Vegetables"
jsonpath "$.parent_id" not exists
[Captures]
id4: jsonpath "$['id']"

# Update Cuisine 3
PUT {{host}}/v1/cuisines/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus", "parent_id":"{{id4}}"}
HTTP 200
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" == {{id4}}

# Update Cuisine - parent_id doesn't exists
PUT {{host}}/v1/cuisines/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Pork", "parent_id": "c624bce3-2d1b-4ae8-87e2-af775be70077"}
HTTP 400
[Asserts]
jsonpath "$.error.parent_id" exists

# Update Cuisine - name empty
PUT {{host}}/v1/cuisines/{{id3}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":""}
HTTP 400
[Asserts]
jsonpath "$.error.name" exists

# List Cuisines
GET {{host}}/v1/cuisines
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 27

# Delete Cuisine 2
DELETE {{host}}/v1/cuisines/{{id2}}
Authorization: Bearer {{token}}
HTTP 204

GET {{host}}/v1/cuisines
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 26

# Delete Cuisine 4: should fail because it's a parent
DELETE {{host}}/v1/cuisines/{{id4}}
Authorization: Bearer {{token}}
HTTP 403

GET {{host}}/v1/cuisines
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count == 26


### Clean up
# ForgetMe (while logged in)
DELETE {{host}}/v1/users
Authorization: Bearer {{token}}
HTTP 204
