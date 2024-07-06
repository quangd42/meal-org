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
# Create Recipe 1
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
      "id": "{{ingre_id1}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "leave as is for decoration",
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
jsonpath "$.ingredients[0].parent_id" not exists
jsonpath "$.ingredients[0].index" == 1
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 2
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "leave as is for decoration"
jsonpath "$.ingredients[2].index" == 3
variable "ingre_in_re_id1" ==  "{{ingre_in_re_id2}}"
jsonpath "$.instructions" count == 2
jsonpath "$.instructions[0].step_no" == 1
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 2
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id1}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"

# Update Recipe 1
PUT {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti",
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
      "index": 2
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 3
    },
    {
      "id": "{{ingre_id3}}",
      "amount": "1lb",
      "prep_note": "leave as is for decoration",
      "index": 4
    }
  ],
  "instructions": [
    {
      "step_no": 2,
      "instruction": "cook the meat till brown"
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
jsonpath "$.ingredients[0].index" == 2
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 3
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "leave as is for decoration"
jsonpath "$.ingredients[2].index" == 4
jsonpath "$.instructions" count == 3
jsonpath "$.instructions[0].step_no" == 2
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 3
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.instructions[2].step_no" == 4
jsonpath "$.instructions[2].instruction" == "put some garnish on"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

# Update Recipe - ingredient_id does not exist
PUT {{host}}/v1/recipes/{{id1}}
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
      "id": "{{ingre_id2}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "leave as is for decoration",
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
HTTP 403
[Asserts]
jsonpath "$.error" exists

# Get Recipe 1 - post Update with ingredient_id that does not exist
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
jsonpath "$.ingredients" count == 2
jsonpath "$.ingredients[0].amount" == "1lb"
jsonpath "$.ingredients[0].prep_note" ==  "seasoned with salt and pepper"
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" ==  "cut to bite size"
jsonpath "$.instructions" count == 3
jsonpath "$.instructions[0].step_no" == 2
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 3
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.instructions[2].step_no" == 4
jsonpath "$.instructions[2].instruction" == "put some garnish on"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

# Update Recipe - duplicate step_no in instructions
PUT {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti back to 1",
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "yield": "a whole pan of goodness",
  "cook_time_in_minutes": 20,
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "{{ingre_id1}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "leave as is for decoration",
      "index": 3
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

# Get Recipe 1 - post Update with bad instructions - expect no change
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
jsonpath "$.ingredients[0].index" == 2
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 3
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "leave as is for decoration"
jsonpath "$.ingredients[2].index" == 4
jsonpath "$.instructions" count == 3
jsonpath "$.instructions[0].step_no" == 2
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 3
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.instructions[2].step_no" == 4
jsonpath "$.instructions[2].instruction" == "put some garnish on"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

# Update Recipe - missing params name in host recipe
PUT {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "external_url": "cooked.wiki/https://www.thekitchn.com/beef-and-broccoli-noodles-recipe-23656929",
  "servings": 4,
  "yield": "a whole pan of goodness",
  "cook_time_in_minutes": 20,
  "notes": "enjoy!",
  "cuisines": ["{{cuisine_id1}}", "{{cuisine_id2}}"],
  "ingredients": [
    {
      "id": "{{ingre_id1}}",
      "amount": "1lb",
      "prep_note": "seasoned with salt and pepper",
      "index": 1
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "2lb",
      "prep_note": "cut to bite size",
      "index": 2
    },
    {
      "id": "{{ingre_id2}}",
      "amount": "1lb",
      "prep_note": "leave as is for decoration",
      "index": 3
    }
  ],
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

# Get Recipe 1 - post Update with missing name - expect no change
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
jsonpath "$.ingredients[0].index" == 2
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 3
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "leave as is for decoration"
jsonpath "$.ingredients[2].index" == 4
jsonpath "$.instructions" count == 3
jsonpath "$.instructions[0].step_no" == 2
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 3
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.instructions[2].step_no" == 4
jsonpath "$.instructions[2].instruction" == "put some garnish on"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

# Update Recipe - no ingredients
PUT {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti back to 1",
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

# Get Recipe 1 - post Update with no ingredients - expect no change
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
jsonpath "$.ingredients[0].index" == 2
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 3
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "leave as is for decoration"
jsonpath "$.ingredients[2].index" == 4
jsonpath "$.instructions" count == 3
jsonpath "$.instructions[0].step_no" == 2
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 3
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.instructions[2].step_no" == 4
jsonpath "$.instructions[2].instruction" == "put some garnish on"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

# Update Recipe - no instructions
PUT {{host}}/v1/recipes/{{id1}}
Authorization: Bearer {{token}}
Content-Type: application/json; charset=utf-8
{
  "name": "Beef and Broccoli Spaghetti back to 1",
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
      "prep_note": "leave as is for decoration",
      "index": 3
    }
  ],
  "instructions": []
}
HTTP 400
[Asserts]
jsonpath "$.error" exists

# Get Recipe 1 - post Update with no instructions - expect no change
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
jsonpath "$.ingredients[0].index" == 2
jsonpath "$.ingredients[1].amount" == "2lb"
jsonpath "$.ingredients[1].prep_note" == "cut to bite size"
jsonpath "$.ingredients[1].index" == 3
jsonpath "$.ingredients[2].amount" == "1lb"
jsonpath "$.ingredients[2].prep_note" == "leave as is for decoration"
jsonpath "$.ingredients[2].index" == 4
jsonpath "$.instructions" count == 3
jsonpath "$.instructions[0].step_no" == 2
jsonpath "$.instructions[0].instruction" == "cook the meat till brown"
jsonpath "$.instructions[1].step_no" == 3
jsonpath "$.instructions[1].instruction" == "put the broccoli in, stir then steam"
jsonpath "$.instructions[2].step_no" == 4
jsonpath "$.instructions[2].instruction" == "put some garnish on"
jsonpath "$.cuisines" count == 2
jsonpath "$.cuisines[*].id" includes "{{cuisine_id2}}"
jsonpath "$.cuisines[*].id" includes "{{cuisine_id3}}"

### Clean up

# Delete Ingredient 3
DELETE {{host}}/v1/ingredients/{{ingre_id3}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 1
DELETE {{host}}/v1/ingredients/{{ingre_id1}}
Authorization: Bearer {{token}}
HTTP 204

# Delete Ingredient 2
DELETE {{host}}/v1/ingredients/{{ingre_id2}}
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
