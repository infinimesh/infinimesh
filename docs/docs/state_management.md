# Device State Management
Infinimesh offers powerful state management based on the JSON data format. Each device owns a `state document`, where it can store its state and transfer it to the platform. Some examples for valid states:

```
true
```

```
17
```

```
"POWER_ON"
```

```
["VALUE_1", "VALUE_2"]
```

```
{
  "color" : "RED",
  "on" : true
}
```

As you can see, a state can be any JSON document. This extremely flexible approach allows you to define your data format as needed. This state document is a primary means of data exchange between the device and business backend applications.

Whenever a state document is transferred to the platform, it is merged with the previous version of the document. This allows devices to send `delta` messages: they do not have transfer the whole state every time, but can just send an update when a state change occurred.

## Send states from a device
The device must publish the JSON 
## Send states to a device
