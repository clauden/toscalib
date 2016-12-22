
package main

import( 
  "fmt"
  "reflect"
  "github.com/davecgh/go-spew/spew"
)

func decode(v interface{}) string {

  if (v == nil) {
    return "nil"
  }

  switch t := v.(type) { 

    case string:
      _ = t
      return fmt.Sprintf("string: %s", "a")

    case map[string]interface {}:
      // m := make(map[string]interface{})
      return fmt.Sprintf("map[string]interface{}: %s", v)

    case map[interface {}]interface {}:
      switch reflect.TypeOf(v).Kind() {

        case reflect.String:
          return fmt.Sprintf("reflected string: %s", "c")

        case reflect.Int:
          return fmt.Sprintf("reflected int: %s", "d")

        case reflect.Map:
          s := ""

          m, ok := v.(map[string]interface{})
          if !ok {
            s = s +  ("unable to promote to map[string]")
            m, ok := v.(map[interface{}]interface{})
            if !ok {
                s = s +  (", unable to promote to map[interface{}]")
                s = s + ", giving up"
            } else {
                s = s + fmt.Sprintf(", promoted to interface %s, %s = %d", m, "foo", m["foo"])
            }
          } else {
            s = s + fmt.Sprintf(", promoted to string %s", m)
          }
          return s
          // return fmt.Sprintf("reflected map: %s", m)
      } 
      default:
        return fmt.Sprintf("no clue")
  }
  return "wtf"
}

func main() {

  x := make(map[interface{}]interface{})
  x["foo"] = 123
  s := decode(x)
  fmt.Printf("%s\n", s)
  fmt.Printf("DUMP:\n")
  spew.Dump(s)
}
