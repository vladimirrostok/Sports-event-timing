# Start from Node base image.
FROM node:14.9

# Add Maintainer info.
LABEL maintainer="Vladimir Andrianov"

# Set the current working directory inside the container.
WORKDIR /usr/src/app

# Install app dependencies.
COPY package*.json ./

RUN npm install

# Add project files.
COPY . ./

# Command to start the app.
CMD ["npm", "start"]