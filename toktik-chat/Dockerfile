FROM alpine
WORKDIR /app

COPY ./bin/toktik-chat .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-chat

EXPOSE 8084 8884
ENTRYPOINT ["./toktik-chat"]
