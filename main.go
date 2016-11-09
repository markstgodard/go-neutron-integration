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
	network_name := fmt.Sprintf("%s", space)
	networks, err = client.NetworksByName(network_name)
	die(err)

	fmt.Printf("Found %d networks\n", len(networks))

	var networkID string
	if len(networks) == 0 {
		// create network
		fmt.Printf("creating network: %s\n", network_name)
		net := neutron.Network{
			Name:         network_name,
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
	subnets, err := client.SubnetsByName(network_name)
	die(err)

	fmt.Printf("Found %d subnets\n", len(subnets))

	if len(subnets) == 0 {
		// create subnet
		fmt.Printf("creating subnet: %s\n", network_name)
		sn := neutron.Subnet{
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

}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}