FROM node:8.17.0-alpine

RUN mkdir /app
WORKDIR /app

COPY package.json .
RUN npm install

COPY . .

# Application port
EXPOSE 3000

# Remote debugging port
EXPOSE 9229

CMD ["npm", "start"]
