FROM openjdk:8
ADD . /go/src/github.com/infinimesh/infinimesh/AvRepo
WORKDIR /go/src/github.com/infinimesh/infinimesh

COPY /go/src/github.com/infinimesh/infinimesh/AvRepo/target/AvRepo-0.0.1-SNAPSHOT.jar  /avrepo-server.jar

EXPOSE 8080
ENTRYPOINT ["java", "-jar", "/avrepo-server.jar"]
 