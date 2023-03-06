FROM centos:7
RUN mkdir /app
ADD server /app/server
ADD worker /app/worker
COPY frontend /app/frontend
COPY script /app/script
COPY sql /app/sql
COPY start_server.sh /app
COPY start_worker.sh /app
WORKDIR /app
ENTRYPOINT ["/bin/bash", "-c"]
