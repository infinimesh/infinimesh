/*
//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
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

import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;
import com.google.protobuf.ByteString;

import io.grpc.Channel;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.StatusRuntimeException;
import main.java.proto.AvroreposGrpc;
import main.java.proto.DeviceState;
import main.java.proto.SaveDeviceStateRequest;
import main.java.proto.SaveDeviceStateResponse;

/**
 * A simple AvroClient that requests to save the device state {@link AvroServer}.
 */
public class AvroClient {
  private static final Logger logger = Logger.getLogger(AvroClient.class.getName());

  private final AvroreposGrpc.AvroreposBlockingStub blockingStub;

  /** Construct client for accessing AvroServer using the existing channel. */
  public AvroClient(Channel channel) {
    // 'channel' here is a Channel, not a ManagedChannel, so it is not this code's responsibility to
    // shut it down.

    // Passing Channels to code makes code easier to test and makes it easier to reuse Channels.
    blockingStub = AvroreposGrpc.newBlockingStub(channel);
  }

  /** Save new device state in avro file at server. */
  public void newDeviceState() {
    logger.info("Will try to save the new Device State ");
    SaveDeviceStateRequest request = SaveDeviceStateRequest.newBuilder()
                    .setDeviceId("0x1321")
                    .setNamespaceId("0x3")
    								.setDs(DeviceState.newBuilder()
    										.setDesiredState(ByteString.EMPTY)
    										.setReportedState(ByteString.copyFromUtf8("I am first one"))
    										.build()
    										)
    								.build();
    SaveDeviceStateResponse response;
    try {
      response = blockingStub.setDeviceState(request);
    } catch (StatusRuntimeException e) {
      logger.log(Level.WARNING, "RPC failed: {0}", e.getStatus());
      return;
    }
    logger.info("Device State Updated: " + response.getStatus());
  }

  /**
   * Greet server. If provided, the first element of {@code args} is the name to use in the
   * greeting. The second argument is the target server.
   */
  public static void main(String[] args) throws Exception {
   
    String target = "localhost:50054";
    // Allow passing in the user and target strings as command line arguments
    if (args.length > 1) {
      target = args[1];
    }

    // Create a communication channel to the server, known as a Channel. Channels are thread-safe
    // and reusable. It is common to create channels at the beginning of your application and reuse
    // them until the application shuts down.
    ManagedChannel channel = ManagedChannelBuilder.forTarget(target)
        // Channels are secure by default (via SSL/TLS). For the example we disable TLS to avoid
        // needing certificates.
        .usePlaintext()
        .build();
    try {
      AvroClient client = new AvroClient(channel);
      client.newDeviceState();
    } finally {
      // ManagedChannels use resources like threads and TCP connections. To prevent leaking these
      // resources the channel should be shut down when it will no longer be used. If it may be used
      // again leave it running.
      channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
    }
  }
}
