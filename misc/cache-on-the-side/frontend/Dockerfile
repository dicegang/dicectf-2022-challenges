FROM node:17-buster-slim AS build
WORKDIR /app

COPY package.json .
RUN npm install
EXPOSE 80
COPY . .
USER node
CMD ["node", "index.js"]
