#!/bin/bash

kitex -module toktik-rpc   -service chat -I E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-rpc\\idl chat.proto

kitex -module toktik-rpc  -service interaction -I E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-rpc\\idl interaction.proto

kitex -module toktik-rpc   -service user -I E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-rpc\\idl user.proto

kitex -module toktik-rpc    -service video -I E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-rpc\\idl video.proto

kitex -module toktik-rpc    -service favor -I E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-rpc\\idl favor.proto

kitex -module toktik-rpc    -service comment -I E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-rpc\\idl comment.proto