FROM node:11.6.0 as builder
RUN npm install -g @angular/cli && npm install @angular-devkit/build-angular

ADD ./frontend /frontend/
WORKDIR /frontend
RUN npm install && ng build --prod
RUN ls /frontend
RUN ls /frontend/dist

FROM nginx:1.15.8-alpine as production-stage
COPY --from=builder /frontend/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY frontend/entrypoint.sh /
EXPOSE 80
CMD ["/entrypoint.sh"]

