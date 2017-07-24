// Copyright Aaron Zinman 2017, 2018
// Copyright Duck Research LLC 2017, 2018
// All rights reserved.
//
// This file is part of Magichaus.
//
// Magichaus is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Magichaus is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Magichaus.  If not, see <http://www.gnu.org/licenses/>.

package dns

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// AWSZoneID for magicha.us
const AWSZoneID = "XXXXXXXXXX"

func publicIP() (string, error) {
	cmd := exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(out.String())
	log.Println("Current IP is:", ip)
	return ip, nil
}

func ipIsCurrent(recordName string, ip string) (bool, error) {
	if strings.HasSuffix(recordName, ".") {
		recordName = strings.TrimSuffix(recordName, ".")
	}
	addrs, err := net.LookupHost(recordName)
	if err != nil {
		return false, err
	}
	for _, addr := range addrs {
		if addr == ip {
			return true, nil
		}
	}
	return false, nil
}

// AWSUpdateToPublicIP sets a route 53 dns entry to the current public IP
func AWSUpdateToPublicIP(zoneID string, recordName string, ttlSeconds int64) error {
	ip, err := publicIP()
	if err != nil {
		return err
	}
	isCurrent, err := ipIsCurrent(recordName, ip)
	if err != nil {
		return err
	}
	if isCurrent {
		log.Println("Record", recordName, "IP is still current at; not updating route53 entry", ip)
		return nil
	}
	c := credentials.NewSharedCredentials("", "DNSUpdateUser")
	if c == nil {
		return errors.New("can't load DNSUpdateUser from AWS credentials")
	}
	svc := route53.New(session.New(aws.NewConfig().WithCredentials(c)))
	params := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				&route53.Change{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: aws.String("A"),
						TTL:  aws.Int64(ttlSeconds),
						ResourceRecords: []*route53.ResourceRecord{
							&route53.ResourceRecord{Value: aws.String(ip)},
						},
					},
				},
			},
		},
	}
	log.Println("Calling to AWS:", params)
	resp, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		return err
	}
	fmt.Println("Changed DNS entry:", resp)
	return nil
}
