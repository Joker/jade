<html>
    <head>
        <title>Dynamic Inline JavaScript</title>
        <script>var users = {{ JSON.stringify(users).replace(/<\//g, "<\\/") }}</script>
    </head>
</html>
