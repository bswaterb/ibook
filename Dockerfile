FROM ubuntu:20.04
LABEL authors="bswaterb"

COPY backend/build-bin/gint /app/gint
WORKDIR /app

ENTRYPOINT ["./gint"]