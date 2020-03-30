FROM	continuumio/miniconda3:4.6.14
RUN	conda install --yes \
	tensorflow \
	future \
	numpy==1.16.3 \
	pandas==0.24.2 \
	&& conda install -c conda-forge fbprophet \
	&& apt-get -y update && apt-get -y install vim bash curl \
	&& pip install prometheus_api_client
