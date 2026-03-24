FROM minio/minio:latest

COPY minio/init.sh /init.sh

RUN chmod +x /init.sh

ENTRYPOINT ["/bin/sh", "/init.sh"]