FROM centos:7
RUN yum update -y && yum install -y subversion git curl
RUN git config --global user.name osc-admin
RUN git config --global user.email osc-admin@oschina.cn
