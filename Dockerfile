FROM golang:1.24.1-alpine3.21 AS gobuilder
RUN apk add --no-cache gcc musl-dev

COPY backend/go.mod backend/go.sum ./

RUN go mod download

WORKDIR /opt/sorrel

COPY backend/ .

ENV CGO_ENABLED=1
RUN go build

FROM node:23.10.0-alpine3.21

WORKDIR /opt/sorrel/

ADD src/ src/
ADD static/ static/
COPY package.json svelte.config.js tsconfig.json vite.config.ts yarn.lock ./

RUN yarn install
RUN yarn run build

COPY --from=gobuilder /opt/sorrel/backend .

CMD ["./backend"]
