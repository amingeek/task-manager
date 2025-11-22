FROM node:20-alpine as build

WORKDIR /app

COPY package*.json ./

RUN npm install --legacy-peer-deps --force

COPY . .

ENV VITE_API_URL=http://localhost:8080/api

RUN npm run build

# Install serve to run the app
RUN npm install -g serve

EXPOSE 3000

CMD ["serve", "-s", "dist", "-l", "3000"]