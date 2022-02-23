/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/submariner-io/cloud-prepare/pkg/api"
)

var _ = Describe("AWS Peering", func() {
	Context("Accept Peering", testCreateVpcPeering)
})

func testCreateVpcPeering() {
	cloudA := newCloudTestDriver(infraID, region)
	_ = Describe("VPC Peering", func() {
		When("receiving a target Cloud", func() {
			It("is an unsupported Cloud", func() {
				invalidCloud := &fooCloud{}
				err := cloudA.cloud.CreateVpcPeering(invalidCloud, api.NewLoggingReporter())
				Expect(err).Should(MatchError("only AWS clients are supported"))
			})
		})
	})
}

type fooCloud struct{}

func (f *fooCloud) PrepareForSubmariner(input api.PrepareForSubmarinerInput, reporter api.Reporter) error {
	panic("not implemented")
}

func (f *fooCloud) CreateVpcPeering(target api.Cloud, reporter api.Reporter) error {
	panic("not implemented")
}

func (f *fooCloud) CleanupAfterSubmariner(reporter api.Reporter) error {
	panic("not implemented")
}

func getRouteTableFor(vpcID string) *ec2.DescribeRouteTablesOutput {
	rtID := vpcID + "-rt"

	return &ec2.DescribeRouteTablesOutput{
		RouteTables: []types.RouteTable{
			{
				VpcId:        &vpcID,
				RouteTableId: &rtID,
			},
		},
	}
}

type cloudTestDriver struct {
	fakeAWSClientBase
	cloud api.Cloud
}

func newCloudTestDriver(infraID, region string) *cloudTestDriver {
	t := &cloudTestDriver{}

	BeforeEach(func() {
		t.beforeEach()
		t.cloud = NewCloud(t.awsClient, infraID, region)
	})

	AfterEach(t.afterEach)

	return t
}

