FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/logs

COPY ./project-service.bin /app

# Run the server executable
CMD [ "/app/project-service.bin" ]