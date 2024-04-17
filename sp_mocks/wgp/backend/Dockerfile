## backend
FROM node:lts-alpine
WORKDIR /usr/app
COPY package*.json ./
RUN npm cache clean --force && rm -rf node_modules && npm install && npm i layer8_middleware
COPY . .
# COPY ../../middleware/dist/middleware.wasm ./middleware.wasm
EXPOSE 8000
# RUN npm run build
CMD ["node", "./server.js"]