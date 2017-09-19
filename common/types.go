package common

import (
	openstack "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

type RosterType map[string]map[string]*openstack.Server
