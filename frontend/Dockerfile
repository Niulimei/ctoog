FROM node:14.16.0 AS builder

COPY . .

RUN npm install

RUN npm run build

# -------------------------------------------------

FROM nginx:1.16.1

RUN mkdir -p /etc/nginx/html/gitee/

COPY --from=builder dist /etc/nginx/html/gitee

COPY --from=builder conf.d /etc/nginx/conf.d

CMD nginx -g 'daemon off;'