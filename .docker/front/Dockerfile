FROM node:20-alpine AS builder
ENV NODE_ENV production

# Add a work directory
WORKDIR /app

# Cache and Install dependencies
COPY ./vue .

RUN npm i
RUN npm run build

# Bundle static assets with nginx
FROM nginx:1.25-alpine as production
ENV NODE_ENV production

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Add your nginx.conf
COPY ./.docker/front/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

WORKDIR /usr/share/nginx/html