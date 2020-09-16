FROM node:latest

ADD ./frontend /frontend/
WORKDIR /frontend
RUN yarn

RUN yarn run build
ENTRYPOINT [ "yarn", "start" ]