package main

import (
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
	if len(networks) == 0 {
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

	// delete network
	if networkID != "" {
		err := client.DeleteNetwork(networkID)
		die(err)
		fmt.Printf("deleted network with networkID: %s\n", networkID)
	}

}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
