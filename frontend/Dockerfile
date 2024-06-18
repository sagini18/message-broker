FROM node:18.12.0-alpine AS build

WORKDIR /app

COPY package.json /app/package.json

RUN npm install --silent

COPY . /app

RUN npm run build

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
