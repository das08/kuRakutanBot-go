FROM node:lts

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

COPY . .

CMD ["npm", "run", "test"]
