# syntax=docker/dockerfile:1

FROM alpine:3.19
ENTRYPOINT [ "/fabric" ]
CMD [ "--help" ]
COPY fabric /fabric