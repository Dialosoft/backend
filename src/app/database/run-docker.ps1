# Build the PostgreSQL image
docker build -t database-mock .

# Execute the PostgreSQL container by the image
docker run -d --name database-postgres -p 5432:5432 database-mock

# Build the Redis image using the custom Dockerfile
docker build -f redis.Dockerfile -t redis-mock .

# Execute the Redis container by the image
docker run -d --name database-redis -p 6379:6379 redis-mock
