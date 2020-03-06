FROM	centos/python-36-centos7
USER	root
WORKDIR /home
COPY	requirements.txt /home
RUN	yum -y install pip3-python vim bash curl && pip install --upgrade pip && pip3 install -r requirements.txt
