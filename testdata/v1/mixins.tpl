<div class="one">
    <div class="dialog">
        <h1>Whoop</h1>
        <p>stuff</p>
    </div>
</div>
<div class="two">{{ $title := "Whoop" }} 
    <div class="dialog">
        <h1>{{ title }}</h1>
        <p>stuff</p>
    </div>
</div>
<div class="three">{{ $title := "Whoop" }}{{ $desc := "Just a mixin" }} 
    <div class="dialog">
        <h1>{{ title }}</h1>
        <p>{{ desc }}</p>
    </div>
</div>
<div id="profile">{{ $user := user }} 
    <div class="user">
        <h2>{{ user.name }}</h2>
        <ul class="pets">{{/* _, pet */}}{{ range pets }}
            <li>{{ pet }}</li>
            {{ end }}
        </ul>
    </div>
</div>
