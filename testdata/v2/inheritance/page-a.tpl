<html>
    <head>
        <title>My Site - {{ title }}</title>
        <script src="/jquery.js"></script>
        <script src="/pets.js"></script>
    </head>
    <body>
        <h1>{{ title }}</h1>
        {{/* var pets = []string{"cat", "dog"} */}}{{/* _, petName */}}{{ range pets }}
        <p>{{ petName }}</p>
        {{ end }}
        <div id="footer">
            <p>some footer content</p>
        </div>
    </body>
</html>
