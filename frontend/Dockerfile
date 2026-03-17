FROM nginx:alpine

COPY frontend/web/index.html /usr/share/nginx/html/index.html
COPY frontend/web/src/ /usr/share/nginx/html/src/
COPY frontend/css/ /usr/share/nginx/html/css/

RUN rm /etc/nginx/conf.d/default.conf

COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
