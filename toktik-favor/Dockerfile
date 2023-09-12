FROM alpine
WORKDIR /app
COPY ./bin/toktik-favor .
COPY ./config .
RUN  mkdir config && mv config.yaml config/config.yaml && mv bootstrap.yaml config/bootstrap.yaml

RUN chmod +x ./toktik-favor

EXPOSE 8085 8885 9195
ENTRYPOINT ["./toktik-favor"]
