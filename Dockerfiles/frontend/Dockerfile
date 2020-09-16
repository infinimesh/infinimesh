FROM node:latest

ADD ./frontend /frontend/
WORKDIR /frontend
RUN yarn

RUN npm run build
ENTRYPOINT [ "npm", "start" ]