
# Create user
POST http://localhost:7331/v1/users
Content-Type: application/json; charset=utf-8
{"name":"julia","email":"test@m.c","password":"testpassword"}

# Login
POST http://localhost:7331/v1/auth/login
Content-Type: application/json; charset=utf-8
{"email":"test@m.c","password":"testpassword"}
HTTP 200


