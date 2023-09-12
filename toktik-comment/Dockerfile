FROM alpine
WORKDIR /app
COPY ./bin/toktik-comment .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-comment

EXPOSE 8086 8886 9196
ENTRYPOINT ["./toktik-comment"]
