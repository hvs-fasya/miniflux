FROM alpine:3.7
EXPOSE 80 8080 8081 443
ENV LISTEN_ADDR 0.0.0.0:8081
#ENV LISTEN_ADDR :https
RUN apk add --update bash
RUN apk --no-cache add ca-certificates
RUN apk add --update libcap
RUN mkdir -p /news
ADD miniflux /usr/local/bin/miniflux
ADD 3-codes-news.json /news/3-codes-news.json
ADD country-json-as-per-Canada.json /news/country-json-as-per-Canada.json
RUN setcap 'cap_net_bind_service=+ep' /usr/local/bin/miniflux
USER nobody
CMD ["/usr/local/bin/miniflux"]
