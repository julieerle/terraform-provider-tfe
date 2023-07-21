// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTFEProjectPolicySet_basic(t *testing.T) {
	rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	tfeClient, err := getClientUsingEnv()
	if err != nil {
		t.Fatal(err)
	}

	org, orgCleanup := createOrganization(t, tfeClient, tfe.OrganizationCreateOptions{
		Name:  tfe.String(fmt.Sprintf("tst-terraform-%d", rInt)),
		Email: tfe.String(fmt.Sprintf("%s@hashicorp.com", randomString(t))),
	})
	t.Cleanup(orgCleanup)

	// Make a project
	prj := createProject(t, tfeClient, org.Name, tfe.ProjectCreateOptions{
		Name: randomString(t),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFEProjectPolicySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEProjectPolicySet_basic(org.Name, prj.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFEProjectPolicySetExists(
						"tfe_project_policy_set.test"),
				),
			},
			{
				ResourceName:      "tfe_project_policy_set.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/policy_set_test", org.Name, prj.ID),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTFEProjectPolicySet_incorrectImportSyntax(t *testing.T) {
	rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	tfeClient, err := getClientUsingEnv()
	if err != nil {
		t.Fatal(err)
	}

	org, orgCleanup := createOrganization(t, tfeClient, tfe.OrganizationCreateOptions{
		Name:  tfe.String(fmt.Sprintf("tst-terraform-%d", rInt)),
		Email: tfe.String(fmt.Sprintf("%s@hashicorp.com", randomString(t))),
	})
	t.Cleanup(orgCleanup)

	// Make a project
	prj := createProject(t, tfeClient, org.Name, tfe.ProjectCreateOptions{
		Name: randomString(t),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEProjectPolicySet_basic(org.Name, prj.ID),
			},
			{
				ResourceName:  "tfe_project_policy_set.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s/tst-terraform-%d", org.Name, rInt),
				ExpectError:   regexp.MustCompile(`Error: invalid project policy set input format`),
			},
		},
	})
}

func testAccCheckTFEProjectPolicySetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(ConfiguredClient)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		policySetID := rs.Primary.Attributes["policy_set_id"]
		if policySetID == "" {
			return fmt.Errorf("No policy set id set")
		}

		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("No project id set")
		}

		policySet, err := config.Client.PolicySets.ReadWithOptions(ctx, policySetID, &tfe.PolicySetReadOptions{
			Include: []tfe.PolicySetIncludeOpt{tfe.PolicySetProjects},
		})
		if err != nil {
			return fmt.Errorf("error reading polciy set %s: %w", policySetID, err)
		}
		for _, project := range policySet.Projects {
			if project.ID == projectID {
				return nil
			}
		}

		return fmt.Errorf("Project (%s) is not attached to policy set (%s).", projectID, policySetID)
	}
}

func testAccCheckTFEProjectPolicySetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(ConfiguredClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tfe_policy_set" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, err := config.Client.PolicySets.Read(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Policy Set %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

// testAccTFEProjectPolicySet_basic

func testAccTFEProjectPolicySet_base(orgName string) string {
	return fmt.Sprintf(`
resource "tfe_policy_set" "test" {
  name         = "policy_set_test"
  description  = "a test policy set"
  global       = false
  organization = "%s"
}
`, orgName)
}

func testAccTFEProjectPolicySet_basic(orgName string, prjID string) string {
	return testAccTFEProjectPolicySet_base(orgName) + fmt.Sprintf(`
resource "tfe_project_policy_set" "test" {
  policy_set_id = tfe_policy_set.test.id
  project_id      = "%s"
}
`, prjID)
}
