FROM golang

#RUN apt-get update && apt-get install -y ca-certificates git-core ssh

#ADD keys/id_rsa_without_pass /root/.ssh/id_rsa
#ADD ~/.ssh/id_rsa_without_pass /root/.ssh/id_rsa
#RUN chmod 700 /root/.ssh/id_rsa
#RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
#RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
#RUN echo "====================================FROM DOCKERFILE\n"
RUN mkdir -p /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
#RUN  ls && make prog 
CMD ["./main"]
#CMD ["ls"]ls