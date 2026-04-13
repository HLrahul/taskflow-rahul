# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy lockfile and package.json
COPY package.json bun.lock* ./
RUN bun install

COPY . .

# Pass VITE_API_URL as a build argument
ARG VITE_API_URL=http://localhost:4000
ENV VITE_API_URL=$VITE_API_URL

RUN bun run build

# Production stage
FROM nginx:stable-alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
