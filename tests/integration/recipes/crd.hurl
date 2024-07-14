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

# Create Ingredient 1
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Beef"}
HTTP 201
[Captures]
ingre_id1: jsonpath "$['id']"
[Asserts]
jsonpath "$.name" == "Beef"
jsonpath "$.parent_id" not exists

# Create Ingredient 2
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Ground Beef", "parent_id":"{{ingre_id1}}"}
HTTP 201
[Captures]
ingre_id2: jsonpath "$['id']"
[Asserts]
jsonpath "$.name" == "Ground Beef"
jsonpath "$.parent_id" == "{{ingre_id1}}"

# Create Ingredient 3
POST {{host}}/v1/ingredients
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{"name":"Asparagus"}
HTTP 201
[Captures]
ingre_id3: jsonpath "$['id']"
[Asserts]
jsonpath "$.name" == "Asparagus"
jsonpath "$.parent_id" not exists

# List Cuisines
GET {{host}}/v1/cuisines
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" count > 0
[Captures]
cuisine_id1: jsonpath "$[0].id"
cuisine_id2: jsonpath "$[1].id"
cuisine_id3: jsonpath "$[2].id"

### Tests
# List Recipes - expect null result
GET {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
HTTP 200
[Asserts]
jsonpath "$" == null

# Create Recipe 1
# Its 3rd ingredient has the same id as 2nd
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "yield": "a whole pan of goodness",
  "cook_time_in_minutes": 20,
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    },
    {
      "id": "{{ingre_id3}}", 
      "amount": "1lb",
      "prep_note": "all this for garnish",
      "index": 3
    }
  ],
  "instructions": [
    {
      "step_no": 1,
      "instruction": "cook the meat till brown"
    },
    {
      "step_no": 2,
      "instruction": "put the broccoli in, stir then steam"
    }
  ]
}
HTTP 201
[Captures]
id1: jsonpath "$['id']"
ingre_in_re_id1: jsonpath "$.ingredients[1].id"
ingre_in_re_id2: jsonpath "$.ingredients[2].id"
[Asserts]
jsonpath "$.name" == "Beef and Broccoli Spaghetti"
jsonpath "$.external_url" == "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929"
jsonpath "$.servings" == 4
jsonpath "$.yield" == "a whole pan of goodness"
jsonpath "$.cook_time_in_minutes" == 20
jsonpath "$.notes" == "enjoy!"
jsonpath "$.ingredients" count == 3
jsonpath "$.ingredients[0].amount" == "1lb"
jsonpath "$.ingredients[0].prep_note" == "seasoned with salt and pepper"
jsonpath "$.ingredients[0].index" == 1
jsonpath "$.ingredients[0].parent_id" not exists
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 2
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "all this for garnish"
jsonpath "$.ingredients[2].index" == 3
variable "ingre_in_re_id1" == "{{ingre_in_re_id2}}"
jsonpath "$.instructions" count == 2
jsonpath "$.instructions[0].step_no" == 1
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 2
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id1}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"

# Create Recipe - ingredient_id does not exist
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "cook_time_in_minutes": 60,
  "yield": "a whole pan of goodness",
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "f689252d-9719-420f-8bc5-26e5c9bde366",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    }
  ],
  "instructions": [
    {
      "step_no": 1,
      "instruction": "cook the meat till brown"
    },
    {
      "step_no": 2,
      "instruction": "put the broccoli in, stir then steam"
    }
  ]
}
HTTP 403
[Asserts]
jsonpath "$.error" exists

# Create Recipe - duplicate step_no in instructions
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "cook_time_in_minutes": 60,
  "yield": "a whole pan of goodness",
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    }
  ],
  "instructions": [
    {
      "step_no": 2,
      "instruction": "just leave the meat raw, raw stuff is the best"
    },
    {
      "step_no": 2,
      "instruction": "put the broccoli in, stir then steam"
    },
    {
      "step_no": 4,
      "instruction": "put some garnish on"
    }
  ]
}
HTTP 403
[Asserts]
jsonpath "$.error" exists

# Create Recipe - missing servings and cook_time_in_minutes in host recipe
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "yield": "a whole pan of goodness",
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    }
  ],
  "instructions": [
    {
      "step_no": 1,
      "instruction": "cook the meat till brown"
    },
    {
      "step_no": 2,
      "instruction": "put the broccoli in, stir then steam"
    }
  ]
}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Recipe - no cuisines
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "cook_time_in_minutes": 60,
  "yield": "a whole pan of goodness",
  "notes": "enjoy!",
  "cuisines": [],
  "ingredients": [
    {
      "id": "f689252d-9719-420f-8bc5-26e5c9bde366",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper"
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size"
    }
  ],
  "instructions": [
    {
      "step_no": 1,
      "instruction": "cook the meat till brown"
    },
    {
      "step_no": 2,
      "instruction": "put the broccoli in, stir then steam"
    }
  ]
}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Recipe - no ingredients
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "yield": "a whole pan of goodness",
  "cook_time_in_minutes": 20,
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [],
  "instructions": [
    {
      "step_no": 2,
      "instruction": "just leave the meat raw, raw stuff is the best"
    },
    {
      "step_no": 3,
      "instruction": "put the broccoli in, stir then steam"
    },
    {
      "step_no": 4,
      "instruction": "put some garnish on"
    }
  ]
}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Recipe - no instructions
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "yield": "a whole pan of goodness",
  "cook_time_in_minutes": 20,
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    }
  ],
  "instructions": []
}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Get Recipe 1
# Its 3rd ingredient has the same id as 2nd
GET {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
HTTP 200
[Asserts]
jsonpath "$.name" == "Beef and Broccoli Spaghetti"
jsonpath "$.external_url" == "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929"
jsonpath "$.servings" == 4
jsonpath "$.yield" == "a whole pan of goodness"
jsonpath "$.cook_time_in_minutes" == 20
jsonpath "$.notes" == "enjoy!"
jsonpath "$.ingredients" count == 3
jsonpath "$.ingredients[0].amount" == "1lb"
jsonpath "$.ingredients[0].prep_note" == "seasoned with salt and pepper"
jsonpath "$.ingredients[0].index" == 1
jsonpath "$.ingredients[0].parent_id" not exists
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 2
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "all this for garnish"
jsonpath "$.ingredients[2].index" == 3
jsonpath "$.instructions" count == 2
jsonpath "$.instructions[0].step_no" == 1
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 2
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id1}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"

# Get Recipe not exists
GET {{host}}/v1/recipes/c624bce3-2d1b-4ae8-87e2-af775be70077
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
HTTP 404
[Asserts]
jsonpath "$.error" exists

# Get Recipe bad id
GET {{host}}/v1/recipes/some-bad-id
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Create Recipe 2
POST {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "General Tso's Chicken",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "yield": "a whole pan of goodness",
  "cook_time_in_minutes": 20,
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id2}}", "{{cuisine_id3}}"],
  "ingredients": [
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    }
  ],
  "instructions": [
    {
      "step_no": 1,
      "instruction": "cook the meat till brown"
    },
    {
      "step_no": 2,
      "instruction": "put the broccoli in, stir then steam"
    }
  ]
}
HTTP 201
[Asserts]
jsonpath "$.name" == "General Tso's Chicken"
jsonpath "$.external_url" == "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929"
jsonpath "$.servings" == 4
jsonpath "$.yield" == "a whole pan of goodness"
jsonpath "$.cook_time_in_minutes" == 20
jsonpath "$.notes" == "enjoy!"
jsonpath "$.ingredients" count == 2
jsonpath "$.ingredients[0].amount" == "1lb"
jsonpath "$.ingredients[0].prep_note" == "seasoned with salt and pepper"
jsonpath "$.ingredients[0].parent_id" not exists
jsonpath "$.ingredients[0].index" == 1
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 2
jsonpath "$.instructions" count == 2
jsonpath "$.instructions[0].step_no" == 1
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 2
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

# List Recipes
GET {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
HTTP 200
[Asserts]
jsonpath "$" count == 2
jsonpath "$[*].name" includes "General Tso's Chicken"

# Delete Recipe 1
DELETE {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
HTTP 204

# List Recipes
GET {{host}}/v1/recipes
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
HTTP 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[*].name" not includes "Beef and Broccoli Spaghetti"


### Clean up

# Delete Ingredient 3
DELETE {{host}}/v1/ingredients/{{ingre_id3}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 2
DELETE {{host}}/v1/ingredients/{{ingre_id2}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 1
DELETE {{host}}/v1/ingredients/{{ingre_id1}}
Authorization: Bearer {{token}}
HTTP 204

GET {{host}}/v1/ingredients
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$" == null

# Delete Recipe 1
DELETE {{host}}/v1/ingredients/{{id1}}
Authorization: Bearer {{token}}
HTTP 204

# ForgetMe (while logged in)
DELETE {{host}}/v1/users
Authorization: Bearer {{token}}
HTTP 204
