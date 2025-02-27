FROM alpine:3.14

WORKDIR /usr/local/bin

RUN apk add gcompat
RUN mkdir -p /input
RUN mkdir -p /merged

VOLUME /input
VOLUME /merged

COPY build/archonpdf .

CMD ["/usr/local/bin/archonpdf", "-input=/input", "-merged=/merged"]
