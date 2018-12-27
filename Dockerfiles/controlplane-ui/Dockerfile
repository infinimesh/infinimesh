FROM node:11.5.0-alpine as builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY ./frontend .
RUN npm run build

# production stage
FROM nginx:1.15.8-alpine as production-stage
COPY --from=builder /app/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY frontend/entrypoint.sh /
EXPOSE 80
CMD ["/entrypoint.sh"]
