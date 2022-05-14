# -------------------------
FROM node:16-alpine AS deps
ARG LEAPX_GITLAB_TOKEN

WORKDIR /app
ENV NPM_CONFIG_LOGLEVEL=error
RUN apk add --no-cache libc6-compat
COPY package*.json ./
RUN npm ci

# -------------------------
FROM node:16-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# -------------------------
FROM nginx:latest
ENV NODE_ENV=production
COPY --from=builder /app/dist /usr/share/nginx/html
COPY entry.sh /docker-entrypoint.d/99-entry.sh
