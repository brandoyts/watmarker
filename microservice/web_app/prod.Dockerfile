# build stage
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci --omit=dev

COPY . .

RUN npm run build

# production stage
FROM nginx:1.25-alpine

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx/conf.d/default.conf.template /etc/nginx/conf.d/default.conf.template

# add startup script
COPY entrypoint.sh /docker-entrypoint.d/99-replace-env.sh
RUN chmod +x /docker-entrypoint.d/99-replace-env.sh

RUN chmod -R 755 /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]