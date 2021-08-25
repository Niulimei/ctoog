FROM centos:7.6.1810
RUN echo "Asia/shanghai" > /etc/timezone;
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ADD ./main /home/root/
ADD ./frontend /home/root/frontend
ADD ./start.sh /home/root/
RUN chmod +x /home/root/start.sh

WORKDIR /home/root
ENTRYPOINT ["/bin/bash", "-c", "/home/root/start.sh"]
