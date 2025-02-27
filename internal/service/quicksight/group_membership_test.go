// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quicksight_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/quicksight"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfquicksight "github.com/hashicorp/terraform-provider-aws/internal/service/quicksight"
)

func TestAccQuickSightGroupMembership_basic(t *testing.T) {
	ctx := acctest.Context(t)
	groupName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	memberName := "tfacctest" + sdkacctest.RandString(10)
	resourceName := "aws_quicksight_group_membership.default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, quicksight.EndpointsID),
		CheckDestroy:             testAccCheckGroupMembershipDestroy(ctx),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMembershipConfig_basic(groupName, memberName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMembershipExists(ctx, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccQuickSightGroupMembership_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	groupName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	memberName := "tfacctest" + sdkacctest.RandString(10)
	resourceName := "aws_quicksight_group_membership.default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, quicksight.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGroupMembershipDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMembershipConfig_basic(groupName, memberName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMembershipExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfquicksight.ResourceGroupMembership(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckGroupMembershipDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).QuickSightConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_quicksight_group_membership" {
				continue
			}
			awsAccountID, namespace, groupName, userName, err := tfquicksight.GroupMembershipParseID(rs.Primary.ID)
			if err != nil {
				return err
			}
			listInput := &quicksight.ListGroupMembershipsInput{
				AwsAccountId: aws.String(awsAccountID),
				Namespace:    aws.String(namespace),
				GroupName:    aws.String(groupName),
			}

			found, err := tfquicksight.FindGroupMembership(ctx, conn, listInput, userName)

			if tfawserr.ErrCodeEquals(err, quicksight.ErrCodeResourceNotFoundException) {
				continue
			}
			if err != nil {
				return err
			}
			if found {
				return fmt.Errorf("QuickSight Group (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckGroupMembershipExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		awsAccountID, namespace, groupName, userName, err := tfquicksight.GroupMembershipParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).QuickSightConn(ctx)

		listInput := &quicksight.ListGroupMembershipsInput{
			AwsAccountId: aws.String(awsAccountID),
			Namespace:    aws.String(namespace),
			GroupName:    aws.String(groupName),
		}

		found, err := tfquicksight.FindGroupMembership(ctx, conn, listInput, userName)
		if err != nil {
			return fmt.Errorf("Error listing QuickSight Group Memberships: %s", err)
		}

		if !found {
			return fmt.Errorf("QuickSight Group Membership (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccGroupMembershipConfig_basic(groupName string, memberName string) string {
	return acctest.ConfigCompose(
		testAccGroupConfig_basic(groupName),
		testAccUserConfig_basic(memberName),
		fmt.Sprintf(`
resource "aws_quicksight_group_membership" "default" {
  group_name  = aws_quicksight_group.default.group_name
  member_name = aws_quicksight_user.%s.user_name
}
`, memberName))
}
