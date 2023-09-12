FROM alpine
WORKDIR /app
COPY ./bin/toktik-user .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-user

EXPOSE 8081 8881 9191
ENTRYPOINT ["./toktik-user"]
