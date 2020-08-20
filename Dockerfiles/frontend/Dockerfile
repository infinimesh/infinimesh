FROM node:latest

ADD ./frontend /frontend/
WORKDIR /frontend
RUN npm install
RUN npm run build

EXPOSE 80

ENTRYPOINT [ "npm", "start" ]
