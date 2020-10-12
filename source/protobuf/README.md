Package protobuf contains the protobuf .proto & .pb.go files used to store and retreive objects from disk
while performing admin operations such as building datasets from raw data, searching/viewing datasets on local storage,
and uploading datasets from local storage to DynamoDB.
Data is serialized and encoded as protobuf to be stored in BoltDB NoSQL key/value store (go type []byte).
These protobuf messages are not used to send data between services and web clients. Objects are retreived
from DynamoDB tables and are encoded using the protobuf messages defined in package svc when sent over network.