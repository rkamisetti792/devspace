FROM node:8.17.0

RUN mkdir /app
WORKDIR /app

COPY package.json .
RUN npm install

COPY . .

CMD ["npm", "start"]
