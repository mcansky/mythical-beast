FROM iron/base
WORKDIR /app
COPY main /app/
EXPOSE $PORT
CMD ["/app/main"]
