// Package network contains operations for building social network graphs from
// object data. This package is still in progress and is not implemented in the
// current version. Generating the social network graph from the data points
// derived from this package will be done by Python's NetworkX package.
package network

import (
	"fmt"

	"github.com/elections/source/donations"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

// PersistNodes saves a list of Nodes to the "network_db.db" database
func PersistNodes(nodes []*Node) error {
	err := createBucket("nodes")
	if err != nil {
		fmt.Println("PersistNodes failed: ", err)
		return fmt.Errorf("PersistNodes failed: %v", err)
	}

	// for each obj
	for _, obj := range nodes {
		err := PutNode(obj)
		if err != nil {
			fmt.Println("PersistNodes failed: ", err)
			return fmt.Errorf("PersistNodes failed: %v", err)
		}
	}
	return nil
}

// PutNode saves a Node object to the "network_db.db" database as a protobuf object
func PutNode(node *Node) error {
	// convert obj to protobuf
	data, err := convNodeToProto(*node) // FIX THIS
	if err != nil {
		fmt.Println("PutNode failed: ", err)
		return fmt.Errorf("PutNodes failed: %v", err)
	}
	// open/create bucket in db/network_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open("db/network_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: PutNode failed: 'network_db.db' failed to open")
		return fmt.Errorf("PutNode failed: 'network_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if err := b.Put([]byte(node.ID), data); err != nil { // serialize k,v
			fmt.Printf("PutNode failed: network_db.db': failed to store donor: %s\n", node.ID)
			return fmt.Errorf("PutNode failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: PutNode failed: 'network_db.db': 'nodes' bucket failed to open")
		return fmt.Errorf("PutNode failed: 'network_db.db': 'nodes' bucket failed to open: %v", err)
	}

	return nil
}

func GetNode(id string) (*Node, error) {
	db, err := bolt.Open("db/network_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetCandidate failed: 'network_db.db' failed to open")
		return nil, fmt.Errorf("GetCandidate failed: 'network_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("nodes")).Get([]byte(id))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetCandidate failed: 'network_db.db': 'nodes' bucket failed to open")
		return nil, fmt.Errorf("GetCandidate failed: 'network_db.db': 'nodes' bucket failed to open: %v", err)
	}

	node, err := convProtoToNode(data)
	if err != nil {
		fmt.Println("GetCandidate failed: decodeCand failed: ", err)
		return nil, fmt.Errorf("GetCandidate failed: decodeCand failed: %v", err)
	}

	return &node, nil
}

func createBucket(name string) error {
	db, err := bolt.Open("db/network_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Printf("FATAL: createBucket failed: 'network_db.db' failed to open\n")
		return fmt.Errorf("createBucket failed: 'network_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			fmt.Printf("FATAL: createBucket failed: 'network_db.db': '%s' bucket failed to open\n", name)
			return fmt.Errorf("'main': FATAL: 'network_db.db': '%s' bucket failed to open: %v", name, err)
		}
		return nil
	}); err != nil {
		fmt.Printf("FATAL: createBucket failed: 'network_db.db': '%s' bucket failed to open\n", name)
		return fmt.Errorf("createBucket failed: 'network_db.db': '%s' bucket failed to open: %v", name, err)
	}

	return nil
}

// convNodeToProto encodes LogData structs as protocol buffers
func convNodeToProto(node Node) ([]byte, error) { // move conversions to protobuf package?
	entry := &NodeProto{
		ID:      node.ID,
		Name:    node.Name,
		Type:    node.Type,
		Edges:   node.Edges,
		Weights: node.Weights,
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("convNodeToProto failed: ", err)
		return nil, fmt.Errorf("convNodeToProto failed: %v", err)
	}
	return data, nil
}

func convProtoToNode(data []byte) (Node, error) {
	node := &NodeProto{}
	err := proto.Unmarshal(data, node)
	if err != nil {
		fmt.Println("convProtoToNode failed: ", err)
		return donations.Candidate{}, fmt.Errorf("convProtoToNode failed: %v", err)
	}

	entry := Node{
		ID:            node.GetID(),
		Name:          node.GetName(),
		Type:          node.GetType(),
		WeightedEdges: decodeInnerMap(node.GetWeightedEdges()),
	}

	return entry, nil
}

func decodeInnerMap(m map[string]*NodeProto_InnerMap) map[string]map[string]float32 {
	we := make(map[string]map[string]float32)
	for k, v := range m {
		we[k] = make(map[string]float32)
		we[k] = v.GetWeights()
	}
	return we
}
