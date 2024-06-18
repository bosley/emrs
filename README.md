# emrs
Environmental Monitoring and Response System

### Develop

```
docker build --tag emrs docker/Dockerfile .

docker run --publish 8080:8080 emrs
```

### Release

```
docker build -t emrs:multistage -f docker/Dockerfile.rel .

docker run --publish 8080:8080 emrs:multistage
```
