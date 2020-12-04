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

package test.java;

import java.io.File;
import java.io.IOException;

import org.apache.avro.Schema;
import org.apache.avro.Schema.Parser;
import org.apache.avro.file.DataFileReader;
import org.apache.avro.file.DataFileWriter;
import org.apache.avro.generic.GenericData;
import org.apache.avro.generic.GenericDatumReader;
import org.apache.avro.generic.GenericDatumWriter;
import org.apache.avro.generic.GenericRecord;
import org.apache.avro.io.DatumReader;
import org.apache.avro.io.DatumWriter;

import com.google.protobuf.Timestamp;

public class GenericMain {
  public static void main(String[] args) throws IOException {
    Schema schema = new Parser().parse(new File("src/main/java/schema/DeviceState.avsc"));

    GenericRecord ds1 = new GenericData.Record(schema);
    ds1.put("DeviceId", "dde");
    ds1.put("Date", "dde");
    ds1.put("Timestamp", Timestamp.newBuilder().setSeconds((60 * 60 * 24) - 1).build().toString());
    ds1.put("NamespaceId", "Alyssa");
    // Leave favorite color null

    GenericRecord ds2 = new GenericData.Record(schema);
    ds2.put("NamespaceId", "Ben");
    ds2.put("Date", "dde");
    ds2.put("DeviceId", "dd1f");
    ds2.put("Timestamp", Timestamp.newBuilder().setSeconds((60 * 60 * 24) - 1).build().toString());

    // Serialize user1 and user2 to disk
    File file = new File("devicestate.avro");
    DatumWriter<GenericRecord> datumWriter = new GenericDatumWriter<GenericRecord>(schema);
    DataFileWriter<GenericRecord> dataFileWriter = new DataFileWriter<GenericRecord>(datumWriter);
    dataFileWriter.create(schema, file);
    dataFileWriter.append(ds1);
    dataFileWriter.append(ds2);
    dataFileWriter.close();

    // Deserialize users from disk
    DatumReader<GenericRecord> datumReader = new GenericDatumReader<GenericRecord>(schema);
    GenericRecord ds3 = null;
    try(DataFileReader<GenericRecord> dataFileReader = new DataFileReader<GenericRecord>(file, datumReader)){
      while (dataFileReader.hasNext()) {
        // Reuse user object by passing it to next(). This saves us from
        // allocating and garbage collecting many objects for files with
        // many items.
        ds3 = dataFileReader.next(ds3);
        System.out.println(ds3);
      }
    }


  }
}
