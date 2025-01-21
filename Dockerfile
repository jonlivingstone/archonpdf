FROM alpine:3.14

WORKDIR /usr/local/bin

RUN apk add gcompat
RUN mkdir -p /input
RUN mkdir -p /merged

VOLUME /input
VOLUME /merged

COPY mergepdf .

CMD ["/usr/local/bin/mergepdf", "-input=/input", "-merged=/merged"]
