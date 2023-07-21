
FROM golang:1.17
# Set destination for COPY
WORKDIR /ont-management-srv

# Copy the source code
COPY . .

# Build
RUN go build -o /ont-management-srv

#run
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh
CMD ["./docker-entrypoint.sh"]