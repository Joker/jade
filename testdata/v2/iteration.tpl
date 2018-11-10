<ul>{{/* _, val */}}{{ range []int{1, 2, 3, 4, 5} }}
    <li>{{ val }}</li>
    {{ end }}
</ul>
<ul>{{/* index, val */}}{{ range []string{"zero", "one", "two"} }}
    <li>{{ strconv.Itoa(index) + ": " + val }}</li>
    {{ end }}
</ul>
<ul>{{/* index, val */}}{{ range map[int]string{1:"one",2:"two",3:"three"} }}
    <li>{{ strconv.Itoa(index) + ": " + val }}</li>
    {{ end }}
</ul>
{{/*  
 qfs := func (condition bool, iftrue, iffalse []string) []string {
       if condition {
           return iftrue
       } else {
           return iffalse
       }
   }
 var values = []string{}
*/}}
<ul>
    {{/* _, val */}}{{ range qfs(len(values)>0, values, []string{"There are no values"}) }}
    <li>{{ val }}</li>
    {{ end }}
</ul>
{{/* var values1 = []string{} */}}
<ul>{{ if gt len values1 0 }}{{/* _, val */}}{{ range values1 }}
    <li>{{ val }}</li>
    {{ end }}{{ else }}
    <li>There are no values1</li>
    {{ end }}
</ul>
{{/* var n = 0; */}}
<ul>{{ range n < 4 }}
    <li>{{ n ; n++ }}</li>
    {{ end }}
</ul>
