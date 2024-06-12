FROM alpine:3.20.0
RUN mkdir /app

COPY ./project-service.bin /app

# Run the server executable
CMD [ "/app/project-service.bin" ]