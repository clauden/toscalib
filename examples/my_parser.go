/*
Copyright 2015 - Olivier Wulveryck

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
package main

import (
  ////	"bytes"
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/owulveryck/toscalib"
	"github.com/owulveryck/toscalib/toscaexec"
	// "gopkg.in/yaml.v2"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	//// "text/template"
)

func composeHelper(val, name string) string {
  s := ""
  if val != "" {  
    s = s + val
  } else {
    s = fmt.Sprintf("__%s-unknown__", name)
  } 
  return s
}


// assemble the text of a node 
func compose(node toscalib.NodeTemplate) string {
  s := ""

  s = s + `<<table border="0" cellspacing="0">` + "\n"
  s = s + `  <tr><td colspan="2" port="port1" border="1" bgcolor="lightblue">` + composeHelper(node.Name, "name") + `</td></tr>` + "\n"
	s = s + `  <tr><td colspan="2" port="port2" border="1">` + composeHelper(node.Type, "type") + `</td></tr>` + "\n"
  s = s + `  <tr>` + "\n"

  s = s + `    <td port="port2" border="1">` +"\n"
  s = s + `      <table border="0" cellspacing="0">` + "\n"
                   s = s + `<tr> <td>REQUIREMENTS</td> </tr>` + "\n"
                   for k,v := range node.Requirements {
                     s = s + `<tr>` + "\n"
                     s = s + `<td>` + composeHelper(fmt.Sprintf("%d", k), "req_number") + `</td>` + "\n"

                     for n,r := range v {
                       // r := toscalib.RequirementAssignment(v)

                       s += `<td>` + composeHelper(n, "__huh__") + `</td>` + "\n"
                       s += `<td>` + composeHelper(r.Capability, "__capab__") + `</td>` + "\n"
                       s += `<td>` + composeHelper(r.Node, "__node__") + `</td>` + "\n"
                       s += `<td>` + composeHelper(r.RelationshipName, "__rname__") + `</td>` + "\n"
                       s = s + `</tr>` + "\n"
                       // s = s + `<tr><td>` + composeHelper(fmt.Sprintf("%s = %s", k, v), "requirements") + `</td></tr>` + "\n"
                     }
                  }
  s = s + `      </table>` + "\n"
  s = s + `    </td>` + "\n"

  s = s + `    <td port="port8" border="1">` + composeHelper(fmt.Sprintf("%s", node.Capabilities), "capabilities") + `</td>` + "\n"
  s = s + `  </tr>` + "\n"
  s = s + `  <tr ><td colspan="2" port="port2" border="1">` + composeHelper(fmt.Sprintf("%s", node.Attributes), "attributes") + `</td></tr>` + "\n"
  s = s + `</table>>` + "\n"

  return s
}


/*
 *  given a TOSCA graph, display a workflow to resolve it
 */
func main() {
	var t toscalib.ServiceTemplateDefinition
	err := t.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

  // dump all the nodes for learning purposes
	fmt.Fprintf(os.Stderr, "ServiceTemplate: %s\n", t.Description)
	for _, p := range t.NodeTypes {
		// fmt.Fprintf(os.Stderr, "Type: %+v\n", p)
		spew.Dump(p)
	}

  // ???
	// out, err := yaml.Marshal(t)
	//fmt.Println(string(out))

	// Creates a new graph
	g := gographviz.NewGraph()
	g.AddAttr("", "rankdir", "LR")
	g.SetName("G")
	g.SetDir(true)

  // resolve the topology to a workflow
	e := toscaexec.GeneratePlaybook(t)
	fmt.Fprintf(os.Stderr, "Playbook: %s\n", e)

  // step through the workflow
	for i, p := range e.Index {
		var label string

		//// err = template.ExecuteTemplate(&out, "NODE", p.NodeTemplate)
		//// fmt.Fprintf(os.Stderr, "template error: %s\n", err)

    // make a nice graphviz-formatted label
    label = compose(p.NodeTemplate)
    fmt.Fprintf(os.Stderr, "Label: %s\n", label)
 
		g.AddNode("G", fmt.Sprintf("%v", i),
			map[string]string {
				"id":    fmt.Sprintf("\"%v\"", i),
				"label": label,
				//"label": fmt.Sprintf("\"%v|%v\"", p.NodeTemplate.Name, p.OperationName),
				"shape": "\"record\"",
			})
	}

	l := e.AdjacencyMatrix.Dim()
	for r := 0; r < l; r++ {
		for c := 0; c < l; c++ {
			if e.AdjacencyMatrix.At(r, c) == 1 {
				g.AddEdge(fmt.Sprintf("%v", r), fmt.Sprintf("%v", c), true, nil)
			}
		}

	}

	s := g.String()
	fmt.Println(s)
}
