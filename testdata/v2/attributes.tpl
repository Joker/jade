<a href="{{ print `google.com`+`google.com` }}">Google</a>
<a class="button" href="google.com">Google</a>
<a class="button" href="google.com">Google</a>
{{/* var authenticated = true */}}
<body class="{{ print authenticated ? "authed" : "anon" }}"></body>
<input type="checkbox" name="agreement" checked="checked"/>
<input data-json="
  {
    &#34;very-long&#34;: &#34;piece of &#34;,
    &#34;data&#34;: true
  }
"/>
<!-- pug error -->
<div class="div-class" (click)="play()"></div>
<div class="div-class" (click)="play()"></div>
<div class="div-class" '(click)'="play()"></div>
<a href="/#{url}">Link{{/* var url = "pug-test.html" */}}
    <a href="{{ print "/" + url }}">Link</a>
    {{/* url = "https://example.com/" */}}
    <a href="{{ print url }}">Another link</a>
    {{/* var btnType = "info" */}}{{/* var btnSize = "lg" */}}
    <button type="button" class="{{ print "btn btn-" + btnType + " btn-" + btnSize }}"></button>
    <button type="button" class="{{ print `btn btn-`+btnType+` btn-`+btnSize+`` }}"></button>
</a>
<div escaped="&lt;code&gt;"></div>
<div unescaped="<code>"></div>
<input type="checkbox" checked="checked"/>
<input type="checkbox" checked="checked"/>
<input type="checkbox"/>
<input type="checkbox" checked="true"/>
<!DOCTYPE html>
<input type="checkbox" checked="checked"/>
<input type="checkbox" checked="checked"/>
<input type="checkbox"/>
<input type="checkbox" checked="{{ print true && "checked" == "checked" }}"/>
<a style="{{ print map[string]string{"color": "red", "background": "green"} }}"></a>
{{/* var classes = []string{"foo", "bar", "baz"} */}}
<a class="{{ print classes }}"></a>
<a class="bang classes [&#39;bing&#39;]"></a>
{{/* var currentUrl = "/about" */}}
<a class="{{ print currentUrl == "/" ? "active" : "" }}" href="/">Home</a>
<a class="{{ print currentUrl == "/about" ? "active" : "" }}" href="/about">About</a>
<a class="button"></a>
<div class="content"></div>
<a id="main-link"></a>
<div id="content"></div>
<div id="foo" data-bar="foo"></div>
{{/* var attributes = struct{class string}{}; */}}{{/* attributes.class = "baz"; */}}
<div id="foo" data-bar="foo"></div>
<zxc class="asd qwe zxc" num="{{ print 1 }}"></zxc>
<zxc num="{{ print 1.1 }}"></zxc>
