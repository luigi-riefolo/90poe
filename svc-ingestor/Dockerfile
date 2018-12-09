FROM scratch

ARG SERVICE
ARG VERSION
ARG PORT

LABEL service=$SERVICE \
    version=$VERSION

COPY $SERVICE /
COPY data/ /data/

EXPOSE $PORT

ENTRYPOINT [ "/ingestor" ]
