FROM alpine:latest
RUN mkdir /app

COPY ./project-service.bin /app

# Run the server executable
CMD [ "/app/project-service.bin" ]