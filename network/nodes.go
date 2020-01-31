package network

import (
	"fmt"
	"github.com/elections/donations"
	"github.com/elections/persist"
)

// Node represents a node on a social network graph
type Node struct {
	ID      string
	Name    string
	Type    string
	WeightedEdges map[string]map[string]float32 // type: edge ID: direct %
}

// NetworkGraph is a social network graph comprised of Node objects
type NetworkGraph struct {
	ID string
	Nodes map[string]*Node
}

func CreateGraph(source interface{}, graphID string) (*NetworkGraph, error) {
		// initialize graph
		graph := &NetworkGraph{
			ID: graphID,
			Nodes: make(map[string]*Node)
		}
		graph, err = populateGraph(source, graph)
		if err != nil {
			fmt.Println("CreateGraph failed: ", err)
			return nil, fmt.Errorf("CreateGraph failed: %v", err)
		}
}


func populateGraph(sourceNodes []interface{}, graph *NetworkGraph) (*NetworkGraph, error) {
	// create root nodes & add to graph
	for _, source := range sourceNodes {
		switch r := source.(type) {
		case *donations.Individual:
			root, err := CreateDonorNode(source)
			if err != nil {
				fmt.Println("CreateGraph failed: ", err)
				return nil, fmt.Errorf("CreateGraph failed: %v", err)
			}
		case *donations.Committee:
			root, err := CreateCmteNode(source)
			if err != nil {
				fmt.Println("CreateGraph failed: ", err)
				return nil, fmt.Errorf("CreateGraph failed: %v", err)
			}
		// case Candidate
		// case DisbRecipient
		}
		graph.Nodes[root.ID] = root

		// find adjacent nodes
		adj, err := findAdjacent(source)
		if err != nil {
			fmt.Println("CreateGraph failed: ", err)
			return nil, fmt.Errorf("CreateGraph failed: %v", err)
		}

		// add adjacent nodes to graph or update if already exist in graph
		for id, node := range adj {
			if graph.Nodes[id] != nil {

			}
		} 

	}	

	
	
	// repeat with each node as new root (BFS)
} 

func findAdjacent(root *Node, adj map[string]*Node) (map[string]*Node, error) {
	// find all nodes adjacent to root node
	for nodeType, weights := range root.WeightedEdges {
		// create nodes for adjacent nodes
		if nodeType == "indvOUT" || "cmteOUT" {  // indv -> cmte  || cmte -> cmte
			for id, amt := range weights {
				// skip if node already exists
				if adj[id] != nil {
					continue
				}

				cmte, err := persist.GetCommittee(id)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				node, err := CreateCmteNode(cmte)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				adj[node.ID] = node
				continue
			}
		}

		if nodeType == "disbOUT" {  // cmte -> disb. rec.
			for id, amt := range weights {
				// skip if node already exists
				if adj[id] != nil {
					continue
				}

				cmte, err := persist.GetDisbRecipient(id)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				node, err := Create  // create DisbRec Node
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				adj[node.ID] = node
				continue
			}
		}

		if nodeType == "indvIN" {  // cmte <- indv
			for id, amt := range weights {
				// skip if node already exists
				if adj[id] != nil {
					continue
				}
				indv, err := persist.GetIndvDonor(id)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				node, err := CreateDonorNode(indv)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				adj[node.ID] = node
				continue
			}
		}

		if nodeType == "cmteIN" {  // cmte <- cmte
			for id, amt := range weights {
				// skip if node already exists
				if adj[id] != nil {
					continue
				}
				cmte, err := persist.GetCommittee(id)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				node, err := CreateCmteNode(cmte)
				if err != nil {
					fmt.Println("CreateGraph failed: ", err)
					return nil, fmt.Errorf("CreateGraph failed: %v", err)
				}
				adj[node.ID] = node
				continue
			}
		}

	// find adjacent nodes recursively 
	for id, node := adj {
		adj, err = findAdjacent(node, adj)
		if err != nil {
			fmt.Println("CreateGraph failed: ", err)
			return nil, fmt.Errorf("CreateGraph failed: %v", err)
		}
	}

		return adj, nil
	}
}

/* Directed Graph - OUT -> (Donor -> Cmte -> Cand/DisbRec) */
func CreateDonorNode(donor *donations.Individual) (*Node, error) {
	weights := make(map[string]map[string]float32)

	// find edges and weights
	weights["indvOUT"] = make(map[string]float32)
	for id, amt := range donor.RecipientsAmt {
		edges = append(edges, id)
		cmte, err := persist.GetCommittee(id)
		if err != nil {
			fmt.Println("CreateDonorNode failed: ", err)
			return nil, fmt.Errorf("CreateDonorNode failed: %v", err)
		}
		weight := amt / cmte.TotalReceived
		weights["indvOUT"][id] = weight
	}

	// initialize node object
	node := &Node{
		ID:      donor.ID,
		Name:    donor.Name,
		Type:    "indv",
		WeightedEdges: weights,
	}

	return node, nil
}

// CreateCmteNode creates a node with edges in both directions
func CreateCmteNode(cmte *donations.Committee) (*Node, error) 
	weights := make(map[string]map[string]float32)

	// find edges and weights of Affiliate Committees (tx out ->)
	weights["cmteOUT"] =  make(map[string]float32)
	for id, amt := range cmte.AffiliatesAmt {
		aff, err := persist.GetCommittee(id)
		if err != nil {
			fmt.Println("CreateDonorNode failed: ", err)
			return nil, fmt.Errorf("CreateDonorNode failed: %v", err)
		}
		weight := amt / aff.TotalReceived
		weights["cmteOUT"][id] = weight
	}

	// find edges and weights of Disbursement Recipients (tx out ->)
	weights["disbOUT"] =  make(map[string]float32)
	for id, amt := range cmte.TopDisbRecipientsAmt {
		rec, err := persist.GetDisbRecipient(id)
		if err != nil {
			fmt.Println("CreateCmteNode failed: ", err)
			return nil, fmt.Errorf("CreateCmteNode failed: %v", err)
		}
		weight := amt / rec.TotalReceived
		weights["disbOUT"][id] = weight
	}

	// find edges and weights of Indvidual Donors (tx in <-)
	weights["indvIN"] =  make(map[string]float32)
	for id, amt := range cmte.TopIndvDonorsAmt {
		if err != nil {
			fmt.Println("CreateCmteNode failed: ", err)
			return nil, fmt.Errorf("CreateCmteNode failed: %v", err)
		}
		weight := amt / cmte.TotalReceived
		weights["indvIN"][id] = weight
	}

	// find edges and weights of Committee Donors (tx in <-)
	weights["cmteIN"] =  make(map[string]float32)
	for id, amt := range cmte.TopCmteDonorsAmt {
		if err != nil {
			fmt.Println("CreateCmteNode failed: ", err)
			return nil, fmt.Errorf("CreateCmteNode failed: %v", err)
		}
		weight := amt / cmte.TotalReceived
		weights["cmteIN"][id] = weight
	}

	// initialize node object
	cmteType := "pac"
	if cmte.Candidate != nil {  // determine committee type
		cmteType = "pcc"
	}
	node := &Node{
		ID:      cmte.ID,
		Name:    cmte.Name,
		Type:    cmteType,
		WeightedEdges: weights,
	}

	return node, nil
}

func CreateCandNode(cand *donations.Candidate) (*Node, error) {
	weights := make(map[string]float32)

	// find edges and weights of Individual donors (tx in <-)
	for id, amt := range cand.TopIndvDonorsAmt {
		edges = append(edges, id)
		cmte, err := persist.GetCommittee(id)
		if err != nil {
			fmt.Println("CreateDonorNode failed: ", err)
			return nil, fmt.Errorf("CreateDonorNode failed: %v", err)
		}
		weight := amt / cmte.TotalReceived
		weights[id] = weight
	}

	// initialize node object
	node := &Node{
		ID:      cand.ID,
		Name:    cand.Name,
		Type:    "cand",
		WeightedEdges: weights,
	}

	return node, nil
}
