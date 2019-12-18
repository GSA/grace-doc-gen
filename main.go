package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

var variableSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "variable", LabelNames: []string{"Name"}},
		{Type: "output", LabelNames: []string{"Name"}},
	},
}

func main() {
	files := []string{"variables.tf", "outputs.tf"}
	for _, f := range files {
		fmt.Printf("\nFile: %s\n", f)
		fmt.Print("|     Name     | Description |\n| ------------ | ----------- |\n")
		p := hclparse.NewParser()
		v, diags := p.ParseHCLFile(f)
		if diags != nil {
			fmt.Printf("error parsing file %q: %v\n", f, diags)
			return
		}
		c, diags := v.Body.Content(&variableSchema)
		if diags != nil {
			fmt.Printf("error getting Content: %v", diags)
			return
		}
		for _, b := range c.Blocks {
			name := b.Labels[0]
			a, diags := b.Body.JustAttributes()
			if diags != nil {
				fmt.Printf("error getting attributes: %v", diags)
				return
			}
			var description string
			if attr, exists := a["description"]; exists {
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &description)
				diags = append(diags, valDiags...)
			}
			fmt.Printf("| %s | %s |\n", name, description)
		}
	}
}
