# Use Ubuntu as the base image
FROM debian:latest

# Install system dependencies
RUN apt-get update && apt-get install -y build-essential openjdk-17-jdk python3 python3-pip g++ gcc golang 
RUN apt-get install -y  nodejs npm
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*
# RUN useradd -u 1000 -m myuser
# USER myuser
# RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh
# RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
# ENV NVM_DIR /usr/local/nvm
# ENV NODE_VERSION 20
# # RUN source ~/.bashrc
# RUN nvm install 20
# # Set up a working directory
WORKDIR /app

# Install Python packages (if needed)
# RUN pip3 install --no-cache-dir numpy pandas                     # Example Python package

# Install Node.js packages (if needed)
RUN npm install -g yarn
RUN apt-get update && apt-get install -y locales && \
    echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && \
    locale-gen en_US.UTF-8 && \
    update-locale LANG=en_US.UTF-8

# Set environment variables for the container
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

ENV GOPATH /go
ENV DONT_POLLUTE_OUTPUT_WITH_UTF8=true
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN echo 'export LANG=en_US.UTF-8' >> /etc/bash.bashrc && \
    echo 'export DONT_POLLUTE_OUTPUT_WITH_UTF8=true' >> /etc/bash.bashrc && \
    echo 'locale-gen en_US.UTF-8' >> /etc/bash.bashrc && \
    echo 'update-locale LANG=en_US.UTF-8' >> /etc/bash.bashrc
# Explicitly set the locale environment variables
ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8
ENV TERM=xterm-256color
# Copy any necessary scripts or files
COPY . .


# Default command (can be overridden)
CMD ["bash"]