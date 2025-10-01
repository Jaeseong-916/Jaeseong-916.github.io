package main

import (
    "html/template"
    "os"
    "path/filepath"
)

type Post struct {
    Title   string
    Content string
    Date    string
}

func main() {
    posts := []Post{
        {Title: "첫 번째 포스트", Content: "내용...", Date: "2024-01-01"},
    }
    
    tmpl := `
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
    <h1>My Blog</h1>
    {{range .Posts}}
        <article>
            <h2>{{.Title}}</h2>
            <p>{{.Date}}</p>
            <p>{{.Content}}</p>
        </article>
    {{end}}
</body>
</html>`

    t := template.Must(template.New("blog").Parse(tmpl))
    
    file, _ := os.Create("index.html")
    defer file.Close()
    
    t.Execute(file, struct{ Title string; Posts []Post }{
        Title: "My Blog",
        Posts: posts,
    })
}