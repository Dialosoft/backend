
# README 'temporally' DialoSoft Backend WIP

This is a more detailed explanation of the steps to run the Dialosoft Backend application.

## Running the Application

### How to Run a Spring Boot Application with Gradle

To run a Spring Boot application using Gradle, follow these steps:

1. Open a terminal or command prompt.

2. Navigate to the root directory of your Spring Boot project.

3. (Optional) Compile the project:
   It's generally not necessary to compile separately as the run command will do this, but if you want to ensure everything compiles correctly first:

```
gradle compile
```

Or using the wrapper:

```
./gradlew compile
```

4. Run the application using one of the following commands:

   a) Using the Gradle wrapper (recommended):
    - On Unix-based systems (Linux, macOS):
      ```
      ./gradlew bootRun
      ```
    - On Windows:
      ```
      gradlew.bat bootRun
      ```

   b) If Gradle is installed globally on your system:
   ```
    gradle bootRun
   ```

5. Gradle will compile your project (if not done already), download necessary dependencies, and run the Spring Boot application.

6. Wait for the application to start. You'll see logs in the console indicating that the application is running.

### Additional Notes

- Ensure your `build.gradle` file includes the Spring Boot plugin:
  ```groovy
  plugins {
      id 'org.springframework.boot' version '2.5.5'
      id 'io.spring.dependency-management' version '1.0.11.RELEASE'
      id 'java'
  }
  ```

- To pass arguments to your application:
  ```
  ./gradlew bootRun --args='--server.port=8081'
  ```
  (Use `gradlew.bat` on Windows)

- To stop the application, press Ctrl+C in the terminal.

- To run the application in debug mode:
  ```
  ./gradlew bootRun --debug-jvm
  ```
  (Use `gradlew.bat` on Windows)

Remember to execute these commands from the root directory of
your Spring Boot project, where the `build.gradle` file is located.

**Note:** The Gradle wrapper (`gradlew` or `gradlew.bat`) is preferred
as it ensures everyone uses the same Gradle version, regardless of what's installed on their system.


### Application Configuration Files

In a Spring Boot microservices architecture, you typically have multiple configuration files. For the moment, we'll focus on 3 important ones:

- `application.yml`: The main configuration file with default settings.
- `application-local.yml`: Configuration for local development environment.
- `application-docker.yml`: Configuration for Docker environment. This one is only used in the container to adjust other configurations in a docker network environment.

These files are usually located in the `src/main/resources` directory of each microservice.

And for the moment these files are only in auth-service and gateway-service.

### Specifying a Profile

To run the application with a specific profile (e.g., local or docker), you can use the `spring.profiles.active` parameter. Here's how to do it with Gradle:

- For local profile:
    ```
  ./gradlew bootRun --args='--spring.profiles.active=local'
  ```

- For Docker profile:
    ```
  ./gradlew bootRun --args='--spring.profiles.active=docker'
    ```


On Windows, replace `./gradlew` with `gradlew.bat`.

This will ensure that Spring Boot loads the corresponding `application-{profile}.yml` file in addition to the main `application.yml`.


### Build and Start Containers

Go to gateway-service directory:

```bash
cd src/gateway-service
```

And do the previous steps to run the spring boot gateway application.
This service already has a docker-compose dependency to manage the
initialization of the docker-compose file.

Or the alternative is run the following command to build images
and start all containers defined in `docker-compose.yml`:

```bash
cd src/gateway-service
docker-compose up --build
```

Add the `-d` flag to run containers in detached mode:

```bash
docker-compose up --build -d
```

**Check the Status of Containers**

Use the following command to check the status of running containers:

```bash
docker-compose ps
```

**Access the Application**

**temporally**: Once the containers are running, you can access the application through
the exposed ports specified in `docker-compose.yml`.

For example, gateway is exposed on port 8080, you can access it at [http://localhost:8080](http://localhost:8080).

**Stop and Remove Containers**

To stop and remove all containers, networks, and volumes defined in `docker-compose.yml`, run:

```bash
docker-compose down
```

To also remove built images, add the `--rmi all` option:

```bash
docker-compose down --rmi all
```

### Additional Configuration

#### auth-service
It could require additional configuration,
such as environment variables, define them in a `.env` file located
in the root directory of **auth-service**.

An example `.env` file can be found in `example.env`.

### Add a New Route for a New Microservice and Configure Roles

To add a new route on the Gateway for a new microservice and configure the necessary roles, follow these steps:

#### Definir la Ruta en application.yml

Open the `src/gateway-service/application.yml` file and add a new entry in the `spring.cloud.gateway.routes` section. For example, for a microservice called `order-microservice`:
```yaml
spring:
  cloud:
    gateway:
      routes:
        - id: order_service_routes
          uri: lb://order-microservice
          predicates:
            - Path=/dialosoft-api/order/**
          filters:
            - StripPrefix=1
            - name: JwtAuthenticationFilter
            - name: RoleBasedAuthorizationFilter
              args:
                requiredRoles: ROLE_USER, ROLE_ADMIN
```

In this example:

- **id**: Unique identifier for the route.
- **uri**: The URI of the microservice, in this case, order-microservice.
- **predicates**: The condition that the request must meet to be directed to this microservice (in this case, any path that starts with /dialosoft-api/order/).
- **filters**: Filters applied to the request before it is redirected.
    - **StripPrefix=1**: Remove the '/dialosoft-api' prefix before sending the request to the microservice.
    - **name: JwtAuthenticationFilter**: Apply the JWT authentication filter.
    - **name: RoleBasedAuthorizationFilter**: Apply the role-based authorization filter.
    - **args**: Arguments for the authorization filter.
        - **requiredRoles**: Roles required to access the route. In this case, both ROLE_USER and ROLE_ADMIN have access.

### Explication of Open Routes:

If you open the file `src/gateway-service/application.yml`, The `openApi Endpoints` property is defined under `app.security.jwt` to specify routes that do not require authentication:

```yaml
app:
  security:
    jwt:
      secret-key: dialosoft-key-123
  # Open non-secured endpoints for the gateway
  openApiEndpoints:
    - /actuator
    - /dialosoft-api/auth/login
    - /dialosoft-api/auth/register
    - /dialosoft-api/auth/refresh-token
```

### Explication of Swagger Routes from Gateway Service:

If you open the file `src/gateway-service/application.yml` and goes to swagger section props beginning with `swagger-ui`,
you will see the following configuration:

```yaml
# SWAGGER
springdoc:
  enable-native-support: true
  api-docs:
    enabled: true
  swagger-ui:
    enabled: true
    path: /swagger-ui.html
    config-url: /v3/api-docs/swagger-config
    # Add microservices api-docs urls to expose them in swagger-ui
    urls:
      - url: /v3/api-docs
        name: API Gateway Service Swagger
        primaryName: API Gateway Service
      - url: /auth-service/v3/api-docs
        name: Auth Service Swagger
        primaryName: Auth Service Swagger
        belongServiceName: auth-microservice
```
In `urls` property, you can indicate the path of the api-docs/swagger-ui and the name 
of the microservice that will be exposed in the gateway swagger-ui. At the same time,
you can indicate the registered service name of the microservice through `belongServiceName` that will be used to
generate the swagger-ui.
So to add a new route in the gateway swagger-ui, you need to add a new entry in the `urls` property.

For example, if you want to add a new route for a supposed `order-microservice` in the gateway swagger-ui,
you need to add the following entry in the `urls` property:

```yaml
      - url: /order-microservice/v3/api-docs
        name: Order Microservice Swagger
        primaryName: Order Microservice
        belongServiceName: order-microservice
```

At the same time, you need obviously to configure openApi and swagger-ui in the microservice.
For instance in a supposed `order-microservice`. If this one is written in Java, you need to add the following configuration:

```yaml
# SWAGGER
springdoc:
  api-docs:
    enabled: true
    path: /auth-service/v3/api-docs
  swagger-ui:
    enabled: true
    path: /auth-service/swagger-ui.html
```

This configuration will let us indicate the path for OpenApi and Swagger-ui, 
which will be used in the gateway swagger-ui. But this configuration will
depend on how is done in each other languages.

**It's important that microservice must have been registered in the discovery service first.**