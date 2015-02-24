# Jade.go - template engine for Go (golang)
Package jade (github.com/Joker/jade) implements Jade-lang templates for generating Go html/template output.

## Jade syntax
example:
```jade
doctype html
html(lang="en")
  head
    title= pageTitle
    script(type='text/javascript').
      if (foo) {
         bar(1 + 5)
      }
  body
    h1 Jade - template engine
    #container.col
      if youAreUsingJade
        p You are amazing
      else
        p Get on it!
      p.
        Jade is a terse and simple
        templating language with a
        strong focus on performance
        and powerful features.
```
becomes
```html
<!DOCTYPE html>                             
<html lang="en">                            
    <head>                                  
        <title>{{ pageTitle }}</title>      
        <script type='text/javascript'>     
            if (foo) {                      
                bar(1 + 5)                  
            }                               
        </script>                           
    </head>                                 
    <body>                                  
        <h1>Jade - template engine</h1>
        <div id="container" class="col">    
            {{ if youAreUsingJade }}        
                <p>You are amazing</p>      
            {{ else }}                      
                <p>Get on it!</p>           
            {{ end }}                       
            <p>                             
                Jade is a terse and simple  
                templating language with a  
                strong focus on performance 
                and powerful features.      
            </p>                            
        </div>                              
    </body>                                 
</html>
```





