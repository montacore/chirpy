## JWT 
- A JWT is a JSON Web Token. It's a cryptographically signed JSON object that contains information about the user. You'll learn about how the cryptography of JWTs work in our `Learn Cryptography` course, for now, it's just important to know that once the token is created by the server, the data in the token can't be changed without the server knowing.

When your server issues a JWT to Bob, Bob can use that token to make requests as Bob to your API, Bob won't be able to change the token to make requests as Alice.

## Assignment
The first building blocks you'll write are the functions for creating and validating JWTs, which will be used in the next lesson to authenticate users.

    1. [ ] Add a`MakeJWT` funtion to your `auth` package.
    ```go 
        func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error)
    ```

    Create and return a JWT using this [JWT library](https://github.com/golang-jwt/jwt), which you can import into your code by running:
    `{bash} go get -u github.com/golang-jwt/jwt/v5`
    2. Create a new token.
        1. Use `{go}jwt.NewWithClaims`
        2. Use `{go}jwt.SigningMethodSH256`
        3. Use `{go} jwt.RegisteredClaims`
            - Set the `Issuer` to "chirpy-access"
            - Set `IssuedAt` to the current time in UTC
            - Set `ExpiresAt` to the current time plus the expiration time(`ExpiresIn`)
            - Set the `Subject` to a stringified version of the user's ID
        4. Use `{go} token.SignedString` to sign the token with the secret key from `tokenSecret`. Refer to the docs for an overview of the different signing methods and their respective key types.
    3. Add a `ValidateJWT` function to your `auth` package:

```go title: Signature
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error)
```
Use the `{go}jwt.ParseWithClaims` function to validate the signature of the JWT and extract the claims into a `{go}*jwt.Token` struct. The `{go}keyFunc` callback must return the same key type (`{go}[]byte`) used when the token was signed. An error will be returned if the token is invalid or has expired.

If all is well with the token, use the `{go} token.Claims` interface to get access to the users Id as a uuid.UUID.

add some unit tests to the auth package. Make sure you can create and validate JWTs and that expired tokens are rejected and JWTs signed wrong secret are rejected.
