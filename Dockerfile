FROM centos:7

RUN echo "now building..."
RUN yum -y install httpd
RUN sed -i '/#ServerName/a ServerName www.example.com:80' /etc/httpd/conf/httpd.conf

ADD ./index.html /var/www/html/

EXPOSE 80
CMD ["/usr/sbin/httpd", "-D", "FOREGROUND"]

# FROM gcr.io/google.com/cloudsdktool/cloud-sdk:slim

# ARG VERSION=1.0.7

# RUN apt install unzip -y
# RUN curl https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_linux_amd64.zip -o terraform.zip
# RUN unzip terraform.zip && rm terraform.zip
# RUN mv terraform /usr/bin