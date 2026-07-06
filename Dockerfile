FROM node:24-slim AS builder

WORKDIR /app

ENV HUSKY=0
ENV DATABASE_URL=/tmp/build.db

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build && npm prune --omit=dev

FROM node:24-slim

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=8080
ENV DATABASE_URL=/app/data/wishlist.db

COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["node", "build"]
