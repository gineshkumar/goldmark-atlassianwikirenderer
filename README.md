# goldmark-atlassianwikirenderer

goldmark-atlassianwikirenderer is a renderer for [goldmark](https://github.com/yuin/goldmark) that allows rendering to [Atlassian Wiki Formatting Notation](https://jira.atlassian.com/secure/WikiRendererHelpAction.jspa?section=all).

# Usage

```go
	md := goldmark.New(
		goldmark.WithExtensions(extension.TaskList),
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(extension.Linkify),
		goldmark.WithExtensions(extension.CJK),
		goldmark.WithExtensions(extension.Strikethrough),
		goldmark.WithExtensions(extension.Footnote),
		goldmark.WithExtensions(extension.DefinitionList),
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(New()),
	)
```