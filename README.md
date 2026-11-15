# Image-Randomizer

![License: MIT](https://img.shields.io/github/license/rl404/image-randomizer)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/rl404/image-randomizer)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/rl404/image-randomizer)](https://hub.docker.com/r/rl404/image-randomizer)
[![CI-BE](https://github.com/rl404/image-randomizer/actions/workflows/ci-be.yml/badge.svg?event=push)](https://github.com/rl404/image-randomizer/actions/workflows/ci-be.yml)
[![CI-FE](https://github.com/rl404/image-randomizer/actions/workflows/ci-fe.yml/badge.svg?event=push)](https://github.com/rl404/image-randomizer/actions/workflows/ci-fe.yml)

Just a simple website to generate random image from your chosen image list.

Can be used for randoming website background image by calling your url in CSS. Or randoming your forum's signature.

Inspired by [sig.grumpybumpers](http://sig.grumpybumpers.com/).

## Development

### Requirement

- [Go](https://golang.org/)
- [NodeJS](https://nodejs.org)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://docker.com) + [Docker compose](https://docs.docker.com/compose/) (optional)

### Step

1. Git clone this repo.

   ```sh
   git clone https://github.com/rl404/image-randomizer
   ```

2. Rename `backend/.env.sample` to `.env` and modify according to your configuration.
3. Rename `frontend/.env.sample` to `.env` and modify according to your configuration.
4. Start the backend service. (Assumed the database is ready to use).

   ```sh
   cd backend

   # migrate tables
   make migrate

   # start api server
   make api
   ```

5. Open new console/terminal and start the frontend service.

   ```sh
   cd frontend
   npm run dev
   ```

6. [http://localhost:31001](http://localhost:31001) and [http://localhost:31002](http://localhost:31002) are ready to use\*.

\*_Port depends on `port` in their `.env` file._

### With Docker + Docker compose

1. Do step 1-3 above.
2. Build and start docker containers.

   ```sh
   docker-compose up
   ```

## License

MIT License

Copyright (c) 2022 Axel
