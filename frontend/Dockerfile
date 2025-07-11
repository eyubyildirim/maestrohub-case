# Dockerfile for the Vue Frontend (Final Version)

# --- Stage 1: Build Stage ---
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# --- Stage 2: Serve Stage ---
FROM nginx:stable-alpine
# Remove the default Nginx configuration
RUN rm /etc/nginx/conf.d/default.conf
# Copy our custom Nginx configuration
COPY nginx/default.conf /etc/nginx/conf.d/
# Copy the built static files
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
