FROM node:23.10.0-alpine3.21

WORKDIR /opt/sorrel/

RUN pwd

ADD src/ src/
ADD static/ static/
ADD backend/ backend/
COPY package.json svelte.config.js tsconfig.json vite.config.ts yarn.lock ./

RUN yarn install

RUN yarn run build

WORKDIR /opt/sorrel/backend/

RUN apk add --no-cache go

RUN go build

CMD ["./backend"]
