ROM alpine:3.7
EXPOSE 80 8080 443
ENV LISTEN_ADDR 0.0.0.0:8080
RUN apk add --update bash
RUN apk --no-cache add ca-certificates
RUN mkdir -p /news
ADD miniflux /usr/local/bin/miniflux
ADD 3-codes-news.json /news/3-codes-news.json
ADD country-json-as-per-Canada.json /news/country-json-as-per-Canada.json
USER nobody
CMD ["/usr/local/bin/miniflux"]