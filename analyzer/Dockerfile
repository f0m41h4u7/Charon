FROM	continuumio/miniconda3:4.6.14

ENV	PYTHONDONTWRITEBYTECODE=true

RUN	conda install --yes \
	nomkl \
	dask==1.2.2 \
	numpy==1.16.3 \
	pandas==0.24.2 \
	tini==0.18.0 \
	&& conda install -c conda-forge fbprophet \
	&& conda clean -afy \
	&& apt-get -y update && apt-get -y install vim bash curl \
	&& pip install prometheus_api_client tsfresh\
	&& find /opt/conda/ -follow -type f -name '*.a' -delete \
	&& find /opt/conda/ -follow -type f -name '*.pyc' -delete \
	&& find /opt/conda/ -follow -type f -name '*.js.map' -delete \
	&& find /opt/conda/lib/python*/site-packages/bokeh/server/static -follow -type f -name '*.js' ! -name '*.min.js' -delete
