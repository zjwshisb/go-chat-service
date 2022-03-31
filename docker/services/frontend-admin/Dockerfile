FROM node:16-buster
VOLUME /data
WORKDIR /data
RUN git clone https://github.com/zjwshisb/service-frontend.git app \
    && cd app \
    && yarn config set registry https://registry.npm.taobao.org \
    && yarn install
WORKDIR /data/app
CMD npm start
