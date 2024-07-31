FROM node as base
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

FROM nginx
COPY --from=base /app/dist/ /usr/share/nginx/html/
EXPOSE 80


