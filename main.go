package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/markstgodard/go-neutron/neutron"
)

func main() {
	var err error

	fmt.Printf("args: %v\n", os.Args)

	url := os.Args[1]
	token := os.Args[2]
	space := os.Args[3]

	log.Printf("URL: %s\n", url)
	log.Printf("token: %s\n", token)
	log.Printf("space: %s\n", space)

	client, err := neutron.NewClient(url, token)
	die(err)

	// get networks
	networks, err := client.Networks()
	die(err)

	fmt.Println("Networks")
	for _, n := range networks {
		fmt.Printf("\tname: %s\n", n.Name)
	}

	fmt.Println("Looking up network by name: ", space)
	networkName := fmt.Sprintf("%s", space)
	networks, err = client.NetworksByName(networkName)
	die(err)

	fmt.Printf("Found %d networks\n", len(networks))

	var networkID string
	switch len(networks) {
	case 0:
		// create network
		fmt.Printf("creating network: %s\n", networkName)
		net := neutron.Network{
			Name:         networkName,
			Description:  fmt.Sprintf("space: %s\n", space),
			AdminStateUp: true,
		}
		n, err := client.CreateNetwork(net)
		die(err)
		fmt.Printf("created network: %s\n", n.Name)
		networkID = n.ID
	case 1:
		networkID = networks[0].ID
	default:
		die(fmt.Errorf("could not find network with name: %s", networkName))
	}

	// find subnet
	fmt.Println("Looking up subnet by name :", space)
	subnets, err := client.SubnetsByName(networkName)
	die(err)

	fmt.Printf("Found %d subnets\n", len(subnets))

	if len(subnets) == 0 {
		// create subnet
		fmt.Printf("creating subnet: %s\n", networkName)
		sn := neutron.Subnet{
			Name:      networkName,
			NetworkID: networkID,
			IPVersion: 4,
			CIDR:      "10.0.3.0/24",
			AllocationPools: []neutron.AllocationPool{
				{
					Start: "10.0.3.20",
					End:   "10.0.3.150",
				},
			},
		}
		s, err := client.CreateSubnet(sn)
		die(err)
		fmt.Printf("created subnet: %s for networkID: %s\n", s.CIDR, networkID)
	}

	// create port
	port := neutron.Port{
		NetworkID:    networkID,
		Name:         "container-id-123",
		DeviceID:     "d6b4d3a5-c700-476f-b609-1493dd9dadc0",
		AdminStateUp: true,
	}

	p, err := client.CreatePort(port)
	if err != nil {
		log.Fatal(err)
	}

	pretty, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created port: %s\n", pretty)

	// delete port
	cleanup := false
	if cleanup {
		err := client.DeletePort(p.ID)
		die(err)
		fmt.Printf("deleted port with ID: %s\n", p.ID)

		// delete network
		if networkID != "" {
			err := client.DeleteNetwork(networkID)
			die(err)
			fmt.Printf("deleted network with networkID: %s\n", networkID)
		}
	}

}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
