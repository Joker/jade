<?xml version="1.0" encoding="utf-8" ?>
<rss version="2.0"></rss>
<channel>
    <title>RSS Title</title>
    <description>Some description here</description>
    <link href="http://google.com"/>
    <lastbuilddate>Mon, 06 Sep 2010 00:01:00 +0000</lastbuilddate>
    <pubdate>Mon, 06 Sep 2009 16:45:00 +0000</pubdate>
    {{/* _, item */}}{{ range items }}
    <item>
        <title>{{ item.title }}</title>
        <description>{{ item.description }}</description>
        <link/>
        {{ item.link }}
    </item>
    {{ end }}
</channel>
