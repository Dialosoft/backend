# build the image
docker build -t database-mock -f ./Dockerfile .

# execute de container by the image
docker run -d --name database-postgres -p 5432:5432 database-mock
