FROM nginx:1.16.1

RUN mkdir -p /etc/nginx/html/gitee/

COPY dist /etc/nginx/html/gitee

COPY conf.d /etc/nginx/conf.d

CMD nginx -g 'daemon off;'