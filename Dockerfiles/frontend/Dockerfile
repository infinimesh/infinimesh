FROM node:9.11.1-alpine as build-stage
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY ./frontend .
RUN npm run build

# production stage
FROM nginx:1.15.8-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
