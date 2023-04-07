# a base image you want everyone to use
FROM node:16.11.1-alpine3.14
RUN apk update && apk add --no-cache bash
COPY ./ ./
RUN npm install
CMD [ "node", "app.js" ]