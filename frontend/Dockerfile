FROM node:16.0.0 AS builder

COPY . .

RUN npm install

RUN npm run build

# -------------------------------------------------

FROM nginx:1.19

RUN mkdir -p /etc/nginx/html/gitee/

COPY --from=builder dist /etc/nginx/html/gitee

COPY --from=builder conf.d /etc/nginx/conf.d

CMD nginx -g 'daemon off;'