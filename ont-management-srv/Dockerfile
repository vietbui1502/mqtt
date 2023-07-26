
FROM golang:1.17 AS builder
# Set destination for COPY
WORKDIR /build_dir
# Copy the source code
COPY . .
# Build
RUN go build

FROM golang:1.17
WORKDIR /ont-management-srv/
# Copy service
COPY --from=builder /build_dir/ont-management-srv ./ont-management-srv
# Copy data
COPY ./data/domain_blacklist ./data/domain_blacklist
# Run services
RUN chmod +x ./ont-management-srv
CMD ["./ont-management-srv"]