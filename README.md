# go-oidc-idp-example

OIDC IdP example with Golang

## Overview

This IdP supports only **Implicit flow**.

## How to use

### IdP

1. `go mod tidy`
2. `go run .`
3. Open `http://localhost:8080`
4. Register account
5. Login with the account

### Client

1. `pnpm i`
2. `pnpm next`
3. Open `http://localhost:3000`
4. Push **Sign In** button
5. Push **Sign in with Go OIDC IdP Example** button
6. Successfully signed in with IdP ID!

## License

Copyright 2023 Rokoucha

Released under the Apache License, Version 2.0, see LICENSE.
