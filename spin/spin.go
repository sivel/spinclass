package spin

import (
	"fmt"
	"sync"
	"time"

	"github.com/rackspace/gophercloud"
	openstack "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/sivel/spinclass/common"
)

type Class struct {
	computeClient *gophercloud.ServiceClient
	Roster        common.RosterType
	Config        common.Config
}

func (c *Class) Create() {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: c.Config.OpenStack.IdentityEndpoint,
		Username:         c.Config.OpenStack.Username,
		APIKey:           c.Config.OpenStack.APIKey,
		Password:         c.Config.OpenStack.Password,
	}

	provider, _ := rackspace.AuthenticatedClient(opts)

	serviceClient, _ := rackspace.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: c.Config.OpenStack.Region,
	})

	c.computeClient = serviceClient
}

func (c *Class) create(co chan<- string, wg *sync.WaitGroup, prefix string, counter int) {
	defer wg.Done()
	server, _ := servers.Create(c.computeClient, servers.CreateOpts{
		Name:      fmt.Sprintf("%s-%04d", prefix, counter),
		ImageRef:  c.Config.OpenStack.ImageRef,
		FlavorRef: c.Config.OpenStack.FlavorRef,
	}).Extract()
	server.AdminPass = ""
	c.Roster[prefix][server.ID] = server
	co <- server.ID
}

func (c *Class) pedal(serverID string, prefix string, dismiss bool) {
	var server *openstack.Server
	for {
		time.Sleep(5 * time.Second)
		server, _ = servers.Get(c.computeClient, serverID).Extract()
		if server != nil {
			c.Roster[prefix][serverID] = server
		}

		if dismiss == false && (server.Status == "ACTIVE" || server.Status == "ERROR") {
			break
		} else if dismiss == true && (server == nil || server.Status == "" || server.Status == "DELETED" || server.Status == "ERROR") {
			c.Roster[prefix][serverID].Status = "DELETED"
			break
		}
	}
	return
}

func (c *Class) Dismiss(prefix string) {
	for serverID, _ := range c.Roster[prefix] {
		go func(serverID string) {
			servers.Delete(c.computeClient, serverID)
			c.pedal(serverID, prefix, true)
		}(serverID)
	}
}

func (c *Class) New(count int, prefix string) []string {
	var servers []string

	c.Roster[prefix] = make(map[string]*openstack.Server)

	co := make(chan string)
	wg := new(sync.WaitGroup)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go c.create(co, wg, prefix, i)
	}

	for i := 0; i < count; i++ {
		serverID := <-co
		go c.pedal(serverID, prefix, false)
		servers = append(servers, serverID)
	}

	wg.Wait()
	close(co)
	return servers
}
