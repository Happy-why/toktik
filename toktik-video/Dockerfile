FROM alpine
WORKDIR /app

RUN apk update && apk add ffmpeg

COPY ./bin/toktik-video .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-video

EXPOSE 8083 8883 9193
ENTRYPOINT ["./toktik-video"]
