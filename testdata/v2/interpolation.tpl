{{/* var title = "On Dogs: Man's Best Friend"; */}}{{/* var author = "enlore"; */}}{{/* var theGreat = "<span>escape!</span>"; */}}
<h1>{{ title }}</h1>
<p>Written with love by {{ author }}</p>
<p>This will be safe: {{ theGreat }}</p>
{{/* var msg = "not my inside voice"; */}}
<p>This is {{ strings.ToUpper(msg) }}</p>
<p>No escaping for {{ `}` }}!</p>
<p>Escaping works with \#{interpolation}</p>
<p>Interpolation works with {{ '#{interpolation}' }} too!</p>
{{/* var riskyBusiness = "<em>Some of the girls are wearing my mother's clothing.</em>"; */}}
<div class="quote">
    <p>Joel: {{ riskyBusiness }}</p>
</div>
<p>
    This is a very long and boring paragraph that spans multiple lines.
    Suddenly there is a <strong>strongly worded phrase</strong> that cannot be <em>ignored</em>.
</p>
<p>And here's an example of an interpolated tag with an attribute: <q lang="es">¡Hola Mundo!</q>
</p>
<p>If I don't write the paragraph with tag interpolation, tags like<strong>strong</strong>and<em>em</em>might produce unexpected results.
</p>
<p>If I do, whitespace is <strong>respected</strong> and <em>everybody</em> is happy. 
</p>
