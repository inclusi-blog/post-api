FROM openjdk:12-alpine

RUN apk --no-cache add --update bash=4.4.19-r1

RUN apk upgrade openssl

RUN apk add git
RUN apk add jq

RUN apk add gettext
RUN apk add python3 py3-pip
RUN pip3 install --upgrade pip
RUN pip3 install awscli --upgrade

ENV FLYWAY_VERSION 9.22.0
ENV USERNAME gola
ENV HOME /home/$USERNAME
ENV WORKDIR /home/$USERNAME
ENV UID 553
ENV GID 994

RUN addgroup -g $GID $USERNAME && \
    adduser -u $UID -G $USERNAME -s /bin/sh -h $HOME -D $USERNAME && \
    chown -R $USERNAME:$USERNAME $HOME

RUN mkdir -p $HOME

# Change to the nearcast user
USER gola
WORKDIR $WORKDIR

RUN wget https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/${FLYWAY_VERSION}/flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && tar -xzf flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && mv flyway-${FLYWAY_VERSION}/* . \
  && rm flyway-commandline-${FLYWAY_VERSION}.tar.gz

ENV PATH="${WORKDIR}/flyway:${PATH}"

COPY database /home/gola/sql
COPY infrastructure/migrate.sh /home/gola


