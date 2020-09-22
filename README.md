# Image-Randomizer

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
    ```
    git clone https://github.com/rl404/image-randomizer
    ```
2. Rename `backend/sample.env` to `.env` and modify according to your configuration.
3. Modify `frontend/.env` according to backend host and port.
4. Start the backend service. (Assumed the database is ready to use).
    ```
    cd backend
    go run main.go
    ```
5. Open new console/terminal and start the frontend service.
    ```
    cd frontend
    npm start
    ```
6. [http://localhost:31001](http://localhost:31001) and [http://localhost:31002](http://localhost:31002) are ready to use*.

**Port depends on `port` in their `.env` file.*

### With Docker + Docker compose

1. Do step 1-3 above.
2. Build and start docker containers.
    ```
    docker-compose up
    ````

## License

MIT License

Copyright (c) 2020 Axel