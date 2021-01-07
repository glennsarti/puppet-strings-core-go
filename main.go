// +build go1.15

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/glennsarti/puppet-strings-core-go/yard"
	"github.com/glennsarti/puppet-strings-core-go/yard/puppet"

	"encoding/json"
)

func main() {
	registry := yard.NewRegistry()

	filename := "/workspaces/puppet-strings-core-go/tests/fixtures/plan.pp"
	//filename := "/workspaces/puppet-strings-core-go/tests/fixtures/type_alias.pp"

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// TODO: Is plan mode always true? probably not
	puppet.Parse(filename, string(content), true, registry)

	// Convert the Yard Registry into the expected JSON output from Puppet-Strings
	result := Output{}
	result.DataTypeAliases = make([]puppet.DataTypeAlias, 0)
	result.PuppetPlans = make([]puppet.PuppetPlan, 0)

	// Convert the registry into an output
	for _, item := range registry.All {
		switch item.(type) {
		case puppet.DataTypeAlias:
			{
				result.DataTypeAliases = append(result.DataTypeAliases, item.(puppet.DataTypeAlias))
			}
		case puppet.PuppetPlan:
			{
				result.PuppetPlans = append(result.PuppetPlans, item.(puppet.PuppetPlan))
			}
		}
	}

	emitJSON(result)
}

func emitJSON(value interface{}) {
	out, _ := json.MarshalIndent(value, "", "  ")
	fmt.Println(string(out))
}

type (
	// Output Puppet Strings YARD output
	Output struct {
		DataTypeAliases []puppet.DataTypeAlias `json:"data_type_aliases"`
		PuppetPlans     []puppet.PuppetPlan    `json:"puppet_plans"`
	}
)
