FROM alpine
WORKDIR /app
COPY ./bin/toktik-api .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-api

EXPOSE 8080
ENTRYPOINT ["./toktik-api"]
