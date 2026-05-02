FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:22-alpine
WORKDIR /app
COPY --from=builder /app ./
RUN if [ -f dist/server/index.js ] && [ ! -f dist/server/server.js ]; then cp dist/server/index.js dist/server/server.js; fi

EXPOSE 80
ENV NODE_ENV=production
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "80"]
