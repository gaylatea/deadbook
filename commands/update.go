package commands

import (
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
	"os"
	"strings"
)

type UpdateCommand struct {
	Ui cli.Ui
}

func (uc *UpdateCommand) Help() string {
	return "Update /etc/hosts.d files with EC2 instance hostnames."
}

func (uc *UpdateCommand) Run(args []string) int {
	// Connect to AWS.
	credentials, err := aws.GetAuth("", "")
	if err != nil {
		uc.Ui.Error(err.Error())
		return 1
	}

	// Figure out which region this instance is in, if we have one.
	// Do some sanity-checking on the arguments we're given.
	var realRegion string

	metadataRegion, err := aws.GetMetaData("placement/availability-zone")
	if err != nil {
		uc.Ui.Error(err.Error())
		realRegion = ""
	} else {
		// Gotta take out the last character of the AZ for the proper
		// region name.
		realRegion = string(metadataRegion[:(len(metadataRegion) - 1)])
	}

	// We'll iterate through all the regions. You need to keep DNS names
	// unique across EC2 for this to work properly.
	for regionName, region := range aws.Regions {
		regionFileName := fmt.Sprintf("/etc/hosts.d/hosts-%s", regionName)

		o := fmt.Sprintf("Updating with information from EC2 region %s", regionName)
		uc.Ui.Info(o)

		region_connection := ec2.New(credentials, region)

		// Multiple hostnames can map to a single IP.
		regionHostEntries := make(map[string][]string)

		// Grab all instances that have a DNS key from the region.
		instanceFilter := ec2.NewFilter()
		instanceFilter.Add("tag-key", "dns")
		instances, _ := region_connection.Instances(nil, instanceFilter)

		for _, reservation := range instances.Reservations {
			for _, instance := range reservation.Instances {
				// Make the tags easier to work with.
				instanceTags := make(map[string]string)
				for _, tag := range instance.Tags {
					instanceTags[tag.Key] = tag.Value
				}

				ipAddress := ""
				if realRegion == regionName {
					ipAddress = instance.PrivateIpAddress
				} else {
					ipAddress = instance.PublicIpAddress
				}

				regionHostEntries[ipAddress] = []string{instanceTags["dns"], instance.InstanceId}
			}
		}

		regionHostOutput := ""
		for ip, hostnames := range regionHostEntries {
			regionHostOutput += fmt.Sprintf("%s %s\n", ip, strings.Join(hostnames, " "))
		}

		// Output the file to /etc/hosts.d/ so they can be concatenated.
		regionFile, err := os.Create(regionFileName)
		if err != nil {
			uc.Ui.Error(err.Error())
			return 1
		}

		regionFile.WriteString(regionHostOutput)
		regionFile.Close()
	}

	return 0
}

func (uc *UpdateCommand) Synopsis() string {
	return "Update /etc/hosts.d files with EC2 instance hostnames."
}
