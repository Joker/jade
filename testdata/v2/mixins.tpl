<ul>
    <li>foo</li>
    <li>bar</li>
    <li>baz</li>
</ul>
<ul>
    <li>foo</li>
    <li>bar</li>
    <li>baz</li>
</ul>
<ul>{{ $name := "cat" }} 
    <li class="pet">{{ name }}</li>
    {{ $name := "dog" }} 
    <li class="pet">{{ name }}</li>
    {{ $name := "pig" }} 
    <li class="pet">{{ name }}</li>
</ul>
{{ $title := "Hello world" }} 
<div class="article">
    <div class="article-wrapper">
        <h1>{{ title }}</h1>
        {{ if len(block) > 0 }}{{/* block */}}{{ else }}
        <p>No content provided</p>
        {{ end }}
    </div>
</div>
{{ $title := "Hello world" }} 
<p>This is my</p>
<p>Amazing article</p>
<div class="article">
    <div class="article-wrapper">
        <h1>{{ title }}</h1>
        {{ if len(block) > 0 }}{{/* block */}}{{ else }}
        <p>No content provided</p>
        {{ end }}
    </div>
</div>
{{ $href := "/foo" }}{{ $name := "foo" }}
{{/* attributes := struct{class string}{class: "btn"} */}}
<a class="{{ print attributes.class }}" href="{{ print href }}">{{ name }}</a>
{{ $href := fn("/foo", "bar", "baz") }}{{ $name := "foo" }}
{{/* attributes := struct{class string}{class: "btn"} */}}
<a class="{{ print attributes.class }}" href="{{ print href }}">{{ name }}</a>
{{ $href := "/foo" }}{{ $name := "foo" }} 
<a href="{{ print href }}">{{ name }}</a>
{{ $title := "Default Title" }} 
<div class="article">
    <div class="article-wrapper">
        <h1>{{ title }}</h1>
    </div>
</div>
{{ $title := "Hello world" }} 
<div class="article">
    <div class="article-wrapper">
        <h1>{{ title }}</h1>
    </div>
</div>
<!-- TODO for string -->
{{ $items := []string{"\"string\"", "2", "3.5", "4"} }}{{ $id := fn("my-list") }} 
<ul id="{{ print id }}">{{/* _, item */}}{{ range items }}
    <li>{{ item }}</li>
    {{ end }}
</ul>
