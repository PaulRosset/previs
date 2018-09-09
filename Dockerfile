FROM node:latest
LABEL maintainer="jacqueslorentzdev@gmail.com"
COPY * /home/node/
WORKDIR /home/node/
RUN ["yarn", "--production=true"]
ENTRYPOINT ["node", "index.js"]