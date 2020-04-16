FROM registry.access.redhat.com/ubi8/python-36

ENV LC_ALL="en_US.UTF-8" \
    LANG="en_US.UTF-8" \
    APP_ROOT="/opt/app-root" \
    USER_ID="cloudpi00"

# some basics
USER root

# Copy extra files to the image
COPY ./root/ /

WORKDIR ${APP_ROOT}

# add the code
COPY src/ src
COPY requirements.txt .

# upgrade pip
RUN pip install --upgrade pip

# install packages
RUN pip install -r requirements.txt

# Drop the root user and make the content of /opt/app-root owned by user 1001
RUN chown -R 1001:0 ${APP_ROOT} && chmod -R ug+rwx ${APP_ROOT}

USER 1001
CMD ["run"]