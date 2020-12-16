FROM openjdk:8
ADD ./pkg /avrepo/ 
WORKDIR /avrepo

COPY target/lib /app/lib
COPY target/AvRepo-0.0.1-SNAPSHOT.jar /app/avrepo-server.jar

EXPOSE 8080
ENTRYPOINT ["java", "-jar", "/app/avrepo-server.jar"]
 