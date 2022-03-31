FROM golang:1.18-bullseye
VOLUME /data
WORKDIR /data
COPY wait-for-it.sh /home/wait-for-it.sh
RUN chmod +x /home/wait-for-it.sh
CMD /home/wait-for-it.sh chat-mysql:3306 -- ./main serve