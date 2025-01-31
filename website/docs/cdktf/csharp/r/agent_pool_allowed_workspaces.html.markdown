---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_agent_pool_allowed_workspaces"
description: |-
  Manages allowed workspaces on agent pools
---


<!-- Please do not edit this file, it is generated. -->
# tfe_agent_pool_allowed_workspaces

Adds and removes allowed workspaces on an agent pool

~> **NOTE:** This resource requires using the provider with Terraform Cloud and a Terraform Cloud
for Business account.
[Learn more about Terraform Cloud pricing here](https://www.hashicorp.com/products/terraform/pricing).

## Example Usage

Basic usage:

```csharp
using Constructs;
using HashiCorp.Cdktf;
/*Provider bindings are generated by running cdktf get.
See https://cdk.tf/provider-generation for more details.*/
using Gen.Providers.Tfe;
class MyConvertedCode : TerraformStack
{
    public MyConvertedCode(Construct scope, string name) : base(scope, name)
    {
        var tfeOrganizationTestOrganization = new Organization.Organization(this, "test-organization", new OrganizationConfig {
            Email = "admin@company.com",
            Name = "my-org-name"
        });
        var tfeWorkspaceTestWorkspace = new Workspace.Workspace(this, "test-workspace", new WorkspaceConfig {
            Name = "my-workspace-name",
            Organization = Token.AsString(tfeOrganizationTestOrganization.Name)
        });
        var tfeAgentPoolTestAgentPool = new AgentPool.AgentPool(this, "test-agent-pool", new AgentPoolConfig {
            Name = "my-agent-pool-name",
            Organization = Token.AsString(tfeOrganizationTestOrganization.Name),
            OrganizationScoped = false
        });
        new AgentPoolAllowedWorkspaces.AgentPoolAllowedWorkspaces(this, "test-allowed-workspaces", new AgentPoolAllowedWorkspacesConfig {
            AgentPoolId = Token.AsString(tfeAgentPoolTestAgentPool.Id),
            AllowedWorkspaceIds = new [] { Token.AsString(tfeWorkspaceTestWorkspace.Id) }
        });
    }
}
```

## Argument Reference

The following arguments are supported:

* `AgentPoolId` - (Required) The ID of the agent pool.
* `AllowedWorkspaceIds` - (Required) IDs of workspaces to be added as allowed workspaces on the agent pool.


## Import

A resource can be imported; use `<AGENT POOL ID>` as the import ID. For example:

```shell
terraform import tfe_agent_pool_allowed_workspaces.foobar apool-rW0KoLSlnuNb5adB
```


<!-- cache-key: cdktf-0.17.0-pre.15 input-9d6c804f088514863a2d3b994f35e4cd7e510b8364e61ba1fa41165766b7d693 -->