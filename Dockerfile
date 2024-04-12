# syntax=docker/dockerfile:1

FROM scratch
ENTRYPOINT [ "/fabric" ]
CMD [ "--help" ]
COPY fabric /fabric