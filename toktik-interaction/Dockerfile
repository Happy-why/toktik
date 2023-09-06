FROM alpine
WORKDIR /app
COPY ./bin/toktik-interaction .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-interaction

EXPOSE 8082 8882
ENTRYPOINT ["./toktik-interaction"]
