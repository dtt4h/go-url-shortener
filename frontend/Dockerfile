FROM nginx:alpine

COPY css/ /usr/share/nginx/html/css/
COPY web/ /usr/share/nginx/html/web/

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
