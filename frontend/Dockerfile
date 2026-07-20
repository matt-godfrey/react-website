# ==================================================
# Stage 1: Install deps and build React app
# ==================================================
FROM node:24 AS build

WORKDIR /app

COPY package*.json ./

# Install project dependencies using npm ci (ensures a clean, reproducible install)
RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .

RUN npm run build

# ==================================================
# Stage 2: Serve generated static files
# ==================================================

FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
