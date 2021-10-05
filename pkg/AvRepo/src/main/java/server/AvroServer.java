/*
//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------
*/

package main.java.server;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.time.LocalDateTime;
import java.util.Date;
import java.util.concurrent.TimeUnit;

import com.google.api.client.util.DateTime;

import org.apache.commons.logging.Log;
import org.apache.hadoop.util.Time;
import org.apache.avro.Schema;
import org.apache.avro.Schema.Parser;
import org.apache.avro.file.CodecFactory;
import org.apache.avro.file.DataFileWriter;
import org.apache.avro.generic.GenericDatumWriter;
import org.apache.avro.io.DatumWriter;
import org.apache.avro.reflect.ReflectDatumWriter;
import org.apache.avro.specific.SpecificDatumWriter;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import main.java.avro.DeviceState;
import main.java.proto.AvroreposGrpc;
import main.java.proto.SaveDeviceStateRequest;
import main.java.proto.SaveDeviceStateResponse;

public class AvroServer {
    Log log = org.apache.commons.logging.LogFactory.getLog(AvroServer.class);
    private Server server;
    
    private void start() throws IOException {
        /* The port on which the server should run */
        int port = 50054; 
        server = ServerBuilder.forPort(port)
        .addService(new AvroreposImpl())
        .build()
        .start();
        log.info("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
          @Override
          public void run() {
            // Use stderr here since the logger may have been reset by its JVM shutdown hook.
            System.err.println("*** shutting down gRPC server since JVM is shutting down");
            try {
              AvroServer.this.stop();
            } catch (InterruptedException e) {
              e.printStackTrace(System.err);
            }
            System.err.println("*** server shut down");
          }
        });
      }

      private void stop() throws InterruptedException {
        if (server != null) {
          server.shutdown().awaitTermination(30, TimeUnit.SECONDS);
        }
      }
    
      /**
       * Await termination on the main thread since the grpc library uses daemon threads.
       */
      private void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
          server.awaitTermination();
        }
      }
    
      /**
       * Main launches the server from the command line.
       */
      public static void main(String[] args) throws IOException, InterruptedException {
        final AvroServer aserver = new AvroServer();
        aserver.start();
        aserver.blockUntilShutdown();
      }
    
      static class AvroreposImpl extends AvroreposGrpc.AvroreposImplBase{
        Log log = org.apache.commons.logging.LogFactory.getLog(AvroreposImpl.class);
        private DateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd");
        private Date date = new Date();
        Boolean status = false;
        @Override
        public void setDeviceState(SaveDeviceStateRequest request,
            StreamObserver<SaveDeviceStateResponse> responseObserver) {
              try{
                status = saveDatoToAvroFile(request);
              }
              catch(IOException e){
                System.out.println(e.getMessage());
              }
              SaveDeviceStateResponse response = SaveDeviceStateResponse.newBuilder().setStatus(status).build();
              responseObserver.onNext(response);
              responseObserver.onCompleted();
        }
        private Boolean saveDatoToAvroFile(SaveDeviceStateRequest req) throws IOException{
          // Serialize ds1, ds2 to disk
          /*Schema schema = null;
          try{
            schema = new Parser().parse(new File("pkg//AvRepo//src//main//java//schema//DeviceState.avsc"));
          }
          catch(FileNotFoundException e){
            log.error(e.getMessage());
            return false;
          }*/
          //Construct via builder
          DeviceState ds = DeviceState.newBuilder()
                          .setDeviceId(req.getDeviceId())
                          .setVersion(req.getVersion())
                          .setNamespaceId(req.getNamespaceId()) 
                          .setReportedState(req.getDs().getReportedState().toString())
                          .setDesiredState(req.getDs().getDesiredState().toString())
                          .build();
          
          File theDir = new File("//avrepo//data//"+req.getNamespaceId()+"//"+req.getDeviceId()+"//"+dateFormat.format(date)+"//"+LocalDateTime.now().getHour());
          if (!theDir.exists()){
            theDir.mkdirs();
          }
          File toWriteFile = new File(theDir+"//"+"device_state_"+req.getDeviceId()+".avro");    
          DatumWriter<DeviceState> datumWriter = new ReflectDatumWriter<DeviceState>(DeviceState.class);        
          DataFileWriter<DeviceState> dataFileWriter = new DataFileWriter<DeviceState>(datumWriter);
          try{
            //dataFileWriter.setMeta("version", req.getVersion());
            dataFileWriter.setMeta("creator", "Infinimesh-Avro");
            dataFileWriter.setCodec(CodecFactory.deflateCodec(5));
            if (toWriteFile.exists()){
              dataFileWriter.appendTo(toWriteFile);
            }else{
              dataFileWriter.create(ds.getSchema(), toWriteFile);
            }
            dataFileWriter.append(ds);
            dataFileWriter.close();
          }
          catch(IOException e){
            log.debug(e.getMessage());
            return false;
          }
          return true;
          /*
          // Deserialize users from disk
          DatumReader<DeviceState> datumReader = new GenericDatumReader<DeviceState>(schema);
          DeviceState ds3 = null;
          try(DataFileReader<DeviceState> dataFileReader = new DataFileReader<DeviceState>(toWriteFile, datumReader)){
            while (dataFileReader.hasNext()) {
              // Reuse user object by passing it to next(). This saves us from
              // allocating and garbage collecting many objects for files with
              // many items.
              ds3 = dataFileReader.next(ds3);
              System.out.println(ds3);
            }
          }
          return true;
          */
        }
      }
}
