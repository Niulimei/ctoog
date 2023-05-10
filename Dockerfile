FROM centos-7-with-git:0.0.1
RUN mkdir -p /app
RUN mkdir -p /app/log
COPY worker /app/worker
COPY server /app/server
COPY frontend /app/frontend
COPY script /app/script
COPY sql /app/sql
COPY start_server.sh /app/start_server.sh
COPY start_worker.sh /app/start_worker.sh
WORKDIR /app
