# 1. update this base image tag to update the base image
# 2. update this path based on the name of your stack (ie. /stackname/envname/servicename)
FROM ACCOUNTID.dkr.ecr.us-west-2.amazonaws.com/test1/rdev/base:v1.0.1
COPY ./ ./
RUN npm install
CMD [ "node", "app.js" ]