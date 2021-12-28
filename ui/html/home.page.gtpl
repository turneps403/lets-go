{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
    <h2>Latest Snippets</h2> 
    {{if .Snippets}}
        <table> 
            <tr></tr>
            {{range .Snippets}} 
                <tr>
                    <td><a href='/snippet?id={{.ID}}'>{{.Title}}</a></td> 
                    <td>{{humanDate .Created | printf "Created: %s"}}</td>
                    <td>#{{.ID}}</td>
                </tr>
            {{end}} 
        </table>
    {{else}}
    <p>There's nothing to see here... yet!</p>
    {{end}}
{{end}}