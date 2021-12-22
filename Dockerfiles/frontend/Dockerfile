FROM node:14-slim

ADD ./frontend /frontend/
WORKDIR /frontend
RUN yarn

RUN yarn run build
ENTRYPOINT [ "yarn", "start" ]