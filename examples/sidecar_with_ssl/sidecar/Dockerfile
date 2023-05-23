FROM node:18-bullseye-slim
COPY ./ ./
RUN npm install
RUN chown -R node:node /var
USER node
CMD [ "node", "app.js" ]