# NFTVC

Нам пришлось вынуждено ускориться и понизить качество кода, но мы сделали

# Запуск полноценный
```git clone https://github.com/BrockenBots/NFTVC.git```
```docker-compose up --build -d```

# У нас есть:
- Фронтенд, где можно запросить сертификат, верифицировать, выпустить нфт
- Файловый сервер для хранения файлов и ссылок достижений (порт 2999)
- Сервис авторизации пользователя
- Сервис профилей для настройки своего / поиска чужого
- Сервис достижений (его не успели доделать, но там должна быть работа с контрактом по rest)
- Крипто сервис, где реализован основной функционал заявки, создания, изменения статуса, получения данных сертификата, добавление и удаление должности

Всё разворачивалось локально, некоторые моменты пришось захардкодить


## На фронтенде
две страницы:
/auth
/profile/1

## На бекенде
для сваггера добавить к адресу /swagger


# frontend
## Запуск проекта

Docker:
```docker-compose up --build -d```

Locally
```bash
bun install
bun dev | bun run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Архитектура

В проекте имеется next app router, поэтому ```src/app``` принадлежит для routing и layout

Но основной архитектурой является FSD - feature sliced design.

## Стек

- Nextjs
- NextUi
- Effector


#Backend
# NFTVC-auth
NFTVC-Auth is an authorization and token management service for the digital profile platform using NFT and Verifiable Credentials. The project is based on authorization via an Ethereum wallet, as well as session management and access tokens using JWT.
# Stack
Golang: The main implementation language
JWT: Generation and validation of access and update tokens
Ethereum: Verification of signatures using the public key of the Ethereum wallet
MongoDB: Data Storage for users and refresh tokens
Redis: A repository for nonces, checking their uniqueness, access token and blacklist tokens
Swagger: Automatic generation of API documentation
# API
1. Sign In With Wallet
POST /api/auth/sign-in

Authorization via wallet, returns nonce for subsequent signature.
Parameters:

WalletPub is the public address of the user's wallet
Answer:

nonce — Generated UUID for signature verification.
2. Verify Signature
POST /api/auth/verify-signature

Wallet signature verification, returns access and refresh tokens.
Parameters:

WalletPub is the public address of the user's wallet
Signature — A signature generated based on nonce
Answer:

access_token — JWT access token
refresh_token — JWT refresh token
3. Refresh Tokens
POST /api/auth/refresh-tokens

Updating tokens with a valid refresh token.
Parameters:

refresh_token — Refresh token
Answer:

access_token — Updated JWT access token
refresh_token — Updated JWT update token
4. Sign Out
POST /api/auth/sign-out

User logout, token revocation.
Parameters:

access_token (in the header) is the JWT token that needs to be revoked.
Answer:

An empty response body for a successful operation.
# Swagger
Access to swagger: http://localhost:<port>/swagger/
# Quickstart
1. git clone <repository_url>
2. Configure the config.yml
3. docker-compose up --build 
