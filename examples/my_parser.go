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
	"fmt"
  "reflect"
	"github.com/awalterschulze/gographviz"
	"github.com/owulveryck/toscalib"
	"github.com/owulveryck/toscalib/toscaexec"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
)



// pretty-print a value
// use 'name' if no value available (actually, just if empty string)
func notreallycomposeHelper(val, name string) string {
  s := ""
  if val != "" {  
    s = s + val
  } else {
    s = fmt.Sprintf("__%s-unknown__", name)
  } 
  return s
}

// try pretty-printing with spew
func composeHelper(val interface{}, name string) string {
    

  /**
  s := "@ "
  r := reflect.ValueOf(val)
  t := r.Type()
  k := r.Kind()
  // v := r.Value()

  s = s + fmt.Sprintf("Type:%s / Kind:%s / Value:%v", t, k, val)
  **/

  s := spew.Sprintf("@ %v @ ", val)

  // s = s + " @"
  return s
}


// assemble the text of a node ... tediously
func compose(node toscalib.NodeTemplate) string {
  s := ""

  s = s + `<<table border="0" cellspacing="0">` + "\n"
  s = s + `  <tr><td colspan="2" port="port1" border="1" bgcolor="lightblue">` + composeHelper(node.Name, "name") + `</td></tr>` + "\n"
	s = s + `  <tr><td colspan="2" port="port2" border="1">` + composeHelper(node.Type, "type") + `</td></tr>` + "\n"

  // Requirements
  s = s + `  <tr>` + "\n"
  s = s + `    <td port="port2" border="1">` +"\n"
  s = s + `      <table border="1" cellspacing="0">` + "\n"
                   s = s + `<tr> <td colspan="5">REQUIREMENTS</td> </tr>` + "\n"
                   s = s + `<tr> <td>index</td> <td>name</td> <td>capability</td> <td> node </td> <td> relationship </td> </tr>` + "\n" 

                   for k,v := range node.Requirements {
                     s = s + `<tr>` + "\n"

                     // row index
                     s = s + `<td>` + fmt.Sprintf("%d", k) + `</td>` + "\n"
                     // s = s + `<td>` + composeHelper(fmt.Sprintf("%d", k), "") + `</td>` + "\n"

                     for n,r := range v {
                       // r := toscalib.RequirementAssignment(v)

                       // requirement name
                       s += `<td>` + composeHelper(n, "") + `</td>` + "\n"

                       // required capability
                       s += `<td>` + composeHelper(r.Capability, "") + `</td>` + "\n"

                       // required node type
                       s += `<td>` + composeHelper(r.Node, "") + `</td>` + "\n"

                       // required relationship
                       s += `<td>` + composeHelper(r.RelationshipName, "") + `</td>` + "\n"

                       s = s + `</tr>` + "\n"

                       // s = s + `<tr><td>` + composeHelper(fmt.Sprintf("%s = %s", k, v), "requirements") + `</td></tr>` + "\n"
                     }
                  }
  s = s + `      </table>` + "\n"
  s = s + `    </td>` + "\n"
  s = s + `  </tr>` + "\n"

  // Capabilities 
  s = s + `  <tr>` + "\n"
  s = s + `    <td port="port2" border="1">` +"\n"
  s = s + `      <table border="1" cellspacing="0">` + "\n"
                   s = s + `<tr> <td colspan="5">CAPABILITIES</td> </tr>` + "\n"
                   if len(node.Capabilities) == 0 {
                     s = s + `<tr> <td colspan="5"> none </td> </tr>` + "\n"
                   } else {
                   
                    for k,v := range node.Capabilities {
                         s = s + `<tr>` + "\n"
                         s = s + `<td>` + fmt.Sprintf("%s", k) + `</td>` + "\n"

                         fmt.Fprintf(os.Stderr, "\n\nk: %v\n", k)
                         if (v != nil) {
                            fmt.Fprintf(os.Stderr, "    Cap type: %T\n", v)
                            switch t := v.(type) {
                              case string:
                                fmt.Fprintf(os.Stderr, "    Capability_string: %s\n", v)

                              case map[string]interface {}:
                                m := make(map[string]interface{})
                                fmt.Fprintf(os.Stderr, "    Capability_map[string]: %s\n", m)

                              case map[interface {}]interface {}:

                                switch reflect.TypeOf(v).Kind() {
                                  case reflect.String:
                                    fmt.Fprintf(os.Stderr, "FOO: %s\n", v)
                                  case reflect.Int:
                                    fmt.Fprintf(os.Stderr, "BAR: %s\n", v)
                                  case reflect.Map:
                                   // Convert it to a PropertyAssignment
                                   pa := reflect.ValueOf(v).Interface().(map[interface{}]interface{})
                                   paa := make(toscalib.PropertyAssignment, 0) 
                                   for kk, vv := range pa {
                                      fmt.Fprintf(os.Stderr, "    Attribute %s:\n", kk)
                                      fmt.Fprintf(os.Stderr, "    Kind: %s\n", reflect.TypeOf(vv).Kind())
                                      fmt.Fprintf(os.Stderr, "    Value-raw: %s\n", vv)
                                      fmt.Fprintf(os.Stderr, "    Value-interface: %s\n", 
                                          reflect.ValueOf(vv).Interface())
                                          // .([]interface{}))

                                      froob, ok := v.(map[interface{}] interface{})
                                      if !ok {
                                        fmt.Fprintf(os.Stderr, "    blorting froob!!!!!!\n")
                                      } else {
                                        fmt.Fprintf(os.Stderr, "    blorting %s\n", froob)
                                        for k,v := range froob {
                                          fmt.Fprintf(os.Stderr, "    %s -> %s\n", k, v)
                                          ass, ok := v.(map[interface{}]interface{}) 
                                          if ok {
                                            for kk, vv := range ass {
                                              fmt.Fprintf(os.Stderr, "         %s -> %s\n", kk, vv)
                                            }
                                          }
                                        }
                                      }
                                        
 /*****
                                      if vv.Kind() == reflect.Map {
                                        for _, key := range vv.MapKeys() {
                                          strct := vv.MapIndex(key)
                                          fmt.Fprintln(os.Stderr, key.Interface(), strct.Interface())
                                        }
                                      } else {
                                        fmt.Fprintf(os.Stderr, "    ....... blort me\n")
                                      }
***/
                                   }
                                    
                                    //for kk, vv := range pa {
							                      //paa[kk.(string)] = reflect.ValueOf(vv).Interface().([]interface{})
                                    //}

                                    fmt.Fprintf(os.Stderr, "    Packed: %s\n", paa)
                                }




                                _, ok := v.(map[string]interface {})
                                if (!ok) {
                                  fmt.Fprintf(os.Stderr, "    Capability_map[not-a-string]: %v\n", v)
                                  fmt.Fprintf(os.Stderr, "    ... %s\n", reflect.TypeOf(v))
                                  fmt.Fprintf(os.Stderr, "    ... %s\n", reflect.TypeOf(v).Key())
                                  // z := reflect.Value(v) 
                                  // fmt.Fprintf(os.Stderr, "    ... %v\n", z)


/***
                                  _, ok2 := v.(map[ map[string]interface{} ]interface {})
                                  if (!ok2) {
                                    fmt.Fprintf(os.Stderr, "    Capability_map[no-clue]: %v\n", v)
                                    fmt.Fprintf(os.Stderr, "    ... %s\n", reflect.TypeOf(v))
                                    fmt.Fprintf(os.Stderr, "    ... %s\n", reflect.TypeOf(v).Key())
                                  } else { 
                                    fmt.Fprintf(os.Stderr, "    Capability_map[turns-out-to-be-a-map]: %v\n", v)
                                    fmt.Fprintf(os.Stderr, "    ... %s\n", reflect.TypeOf(v))
                                    fmt.Fprintf(os.Stderr, "    ... %s\n", reflect.TypeOf(v).Key())
                                  }
***/
                                } else {
                                  fmt.Fprintf(os.Stderr, "    Capability_map[turns-out-to-be-a-string]: %v\n", v)
                                }
                                
                              default:
                                fmt.Fprintf(os.Stderr, "Capability_other: %T\n", t)
                            }
                           // fmt.Fprintf(os.Stderr, "Capability: %s\n", v.(map[string]string))
                           /*******
                              m := make(map[string]interface{})
                                  
                                  
                              fmt.Fprintf(os.Stderr, "map: %s\n", m)
                            *******/
                         } else {
                           fmt.Fprintf(os.Stderr, "v is nil\n")
                        }

                         
                       // for j,z := range v {
                         // fmt.Fprintf(os.Stderr, "%T, %T", j, z)
                        //}

                       // s = s + `<td>` + composeHelper(fmt.Sprintf("%v", v["properties"]), "") + `</td>` + "\n"
                       s = s + `</tr>` + "\n"

                     }
                  }
  // s = s + `    <td port="port8" border="1">` + composeHelper(fmt.Sprintf("%s", node.Capabilities), "capabilities") + `</td>` + "\n"

  s = s + `      </table>` + "\n"
  s = s + `    </td>` + "\n"

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
	/***
  fmt.Fprintf(os.Stderr, "ServiceTemplate: %s\n", t.Description)
	for _, p := range t.NodeTypes {
		// fmt.Fprintf(os.Stderr, "Type: %+v\n", p)
		spew.Fdump(os.Stderr, p)
	}
  ***/

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
	// fmt.Fprintf(os.Stderr, "Playbook: %s\n", e)

  // step through the workflow
	for i, p := range e.Index {
		var label string

    // this would just run template agains the node...
		//// err = template.ExecuteTemplate(&out, "NODE", p.NodeTemplate)
		//// fmt.Fprintf(os.Stderr, "template error: %s\n", err)

    // i don't grok template, so instead just make a nice graphviz-formatted label
    label = compose(p.NodeTemplate)
    // fmt.Fprintf(os.Stderr, "Label: %s\n", label)
 
		g.AddNode("G", fmt.Sprintf("%v", i),
			map[string]string {
				"id":    fmt.Sprintf("\"%v\"", i),

				"label": label,

        // already computed above
				////"label": fmt.Sprintf("\"%v|%v\"", p.NodeTemplate.Name, p.OperationName),

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
