FROM python:3.10.4-alpine3.15

RUN pip install mkdocs
RUN pip install mkdocs-material
RUN pip install mkdocs-render-swagger-plugin
RUN pip install mkdocs-build-plantuml-plugin

WORKDIR /app
COPY docs/ /app/docs
COPY mkdocs.yml /app/

EXPOSE 8000

CMD [ "mkdocs", "serve", "--dev-addr", "0.0.0.0:8000" ]