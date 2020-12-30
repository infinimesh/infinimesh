FROM openjdk:8
ADD ./pkg/AvRepo/target /avrepo/ 

RUN ls
WORKDIR /avrepo/

RUN ls
COPY . avrepo-server.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "/avrepo-server.jar"]
 
