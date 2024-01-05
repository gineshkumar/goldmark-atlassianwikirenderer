package atlassianwikirenderer

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	input    string
	expected string
}

func TestHeadingRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "# Heading 1",
		expected: "h1.Heading 1",
	}, testCase{
		input:    "## Heading 2",
		expected: "h2.Heading 2",
	}, testCase{
		input:    "### Heading 3",
		expected: "h3.Heading 3",
	}, testCase{
		input:    "#### Heading 4",
		expected: "h4.Heading 4",
	}, testCase{
		input:    "##### Heading 5",
		expected: "h5.Heading 5",
	}, testCase{
		input:    "###### Heading 6",
		expected: "h6.Heading 6",
	}, testCase{
		input:    "## Subheading",
		expected: "h2.Subheading",
	})
}

func TestBlockquoteRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "> This is a blockquote",
		expected: "{quote}\nThis is a blockquote\n{quote}",
	})
}

func TestFencedCodeBlockRendering(t *testing.T) {
	runTest(t, testCase{
		input:    fmt.Sprintln("```\nHello World\n```"),
		expected: "{code}\nHello World\n{code}",
	}, testCase{
		input:    fmt.Sprintf("```go\nfmt.Println(\"Hello World\")\n```"),
		expected: "{code:go}\nfmt.Println(\"Hello World\")\n{code}",
	})
}

func TestParagraphRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "This is a **paragraph**. _It can contain_ text,\nline breaks, and `various` formatting options.\nYou can separate paragraphs by leaving a blank line between them.",
		expected: "This is a *paragraph*. _It can contain_ text,\nline breaks, and {{various}} formatting options.\nYou can separate paragraphs by leaving a blank line between them.",
	})
}

func TestCodeSpanRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "This is an `inline code span`",
		expected: "This is an {{inline code span}}",
	}, testCase{
		input:    "Change the variable name from `CodeSpan` to `codeSpan`.",
		expected: "Change the variable name from {{CodeSpan}} to {{codeSpan}}.",
	})
}

func TestEmphasisRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "*Emphasised text1*",
		expected: "_Emphasised text1_",
	}, testCase{
		input:    "_Emphasised text2_",
		expected: "_Emphasised text2_",
	}, testCase{
		input:    "<i>Emphasised text3</i>",
		expected: "_Emphasised text3_",
	}, testCase{
		input:    "**Emphasised text4**",
		expected: "*Emphasised text4*",
	}, testCase{
		input:    "__*Emphasised text5*__",
		expected: "*_Emphasised text5_*",
	})
}

func TestLinkRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "[Link](https://www.jfrog.com)",
		expected: "[Link|https://www.jfrog.com]",
	})
}

func TestAutoLinkRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "<https://www.jfrog.com>",
		expected: "[https://www.jfrog.com]",
	})
}

func TestStrikethroughRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "~~Strikethrough1~~",
		expected: "-Strikethrough1-",
	}, testCase{
		input:    "<s>Strikethrough2</s>",
		expected: "-Strikethrough2-",
	}, testCase{
		input:    "<strike>Strikethrough3</strike>",
		expected: "<strike>Strikethrough3</strike>",
	})
}

func TestListItemRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "- List\n  - List item",
		expected: "*  List\n**  List item",
	}, testCase{
		input:    "- List\n  - List item\n    - Nested List Item",
		expected: "*  List\n**  List item\n***  Nested List Item",
	}, testCase{
		input:    "1. fruits\n     * apple\n     * banana\n  2. vegetables\n     - carrot\n     - broccoli",
		expected: "#   fruits\n**    apple\n**    banana\n#     vegetables\n**  carrot\n**  broccoli",
	})
}

func TestTableRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "| Month    | Assignee | Backup |\n| -------- | -------- | ------ |\n| January  | Dave     | Steve  |\n| February | Gregg    | Karen  |\n| March    | Diane    | Jorge  |",
		expected: "||Month||Assignee||Backup||\n|January|Dave|Steve|\n|February|Gregg|Karen|\n|March|Diane|Jorge|",
	}, testCase{
		input:    "<table><tr><th>Month</th><th>Assignee</th><th>Backup</th></tr><tr><td>January</td><td>Dave</td><td>Steve</td></tr><tr><td>February</td><td>Gregg</td><td>Karen</td></tr><tr><td>March</td><td>Diane</td><td>Jorge</td></tr></table>",
		expected: "{code:html}<table><tr><th>Month</th><th>Assignee</th><th>Backup</th></tr><tr><td>January</td><td>Dave</td><td>Steve</td></tr><tr><td>February</td><td>Gregg</td><td>Karen</td></tr><tr><td>March</td><td>Diane</td><td>Jorge</td></tr></table>{code}",
	})
}

func TestImageRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "![JFrog logo](https://speedmedia.jfrog.com/08612fe1-9391-4cf3-ac1a-6dd49c36b276/https://media.jfrog.com/wp-content/uploads/2021/12/29113553/jfrog-logo-2022.svg)",
		expected: "!https://speedmedia.jfrog.com/08612fe1-9391-4cf3-ac1a-6dd49c36b276/https://media.jfrog.com/wp-content/uploads/2021/12/29113553/jfrog-logo-2022.svg!",
	})
}

func TestFootnoteRendering(t *testing.T) {
	t.Skip("Test skipped as Atlassian currently does not support Footnotes")
	runTest(t, testCase{
		input:    "[^1]: This is the footnote.",
		expected: "[^1]: This is the footnote.",
	})
}

func TestFootNoteLinkRendering(t *testing.T) {
	t.Skip("Test skipped as Atlassian currently does not support Footnotes")
	runTest(t, testCase{
		input:    "Here is some text with a footnote[^1].\n[^1]: This is the footnote content.",
		expected: "Here is some text with a footnote[^1].\n[^1]: This is the footnote content.",
	})
}

func TestFootNoteListRendering(t *testing.T) {
	t.Skip("Test skipped as Atlassian currently does not support Footnotes")
	runTest(t, testCase{
		input:    "[^1]: This is the content of the first footnote.\n[^2]: This is the content of the second footnote.\n",
		expected: "[^1]: This is the content of the first footnote.\n[^2]: This is the content of the second footnote.\n",
	})
}

func TestDefinitionTermDescriptionRendering(t *testing.T) {
	t.Skip("Test skipped as Atlassian currently does not support Definitions")
	runTest(t, testCase{
		input:    "Definition Term\n:   This is the description for Term.",
		expected: "Definition Term\n:   This is the description for Term.",
	})
}

func TestHTMLBlockRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "```html\n<div>\n  <p>This is an HTML block.</p>\n</div>",
		expected: "{code:html}\n<div>\n  <p>This is an HTML block.</p>\n</div>{code}",
	})
}

func TestRawHTMLRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "<div>\n  <p>This is a raw HTML block.</p>\n</div>",
		expected: "{code:html}<div>\n  <p>This is a raw HTML block.</p>\n</div>{code}",
	})
}

func TestThematicBreakRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "Text above the thematic break.\n\n---\n\nText below the thematic break.",
		expected: "Text above the thematic break.\n\n----\n\nText below the thematic break.",
	})
}

func TestEntireDocument(t *testing.T) {
	inputFile := "./input.txt"
	expectedFile := "./expected.txt"

	inputContent, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Error reading input file: %v", err)
	}

	expectedContent, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Error reading expected file: %v", err)
	}

	runTest(t, testCase{
		input:    string(inputContent),
		expected: string(expectedContent),
	})
}

func runTest(t *testing.T, tests ...testCase) {
	renderer := New()
	md := goldmark.New(
		goldmark.WithRenderer(renderer),
		goldmark.WithExtensions(extension.TaskList),
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(extension.Linkify),
		goldmark.WithExtensions(extension.CJK),
		goldmark.WithExtensions(extension.Strikethrough),
		goldmark.WithExtensions(extension.Footnote),
		goldmark.WithExtensions(extension.DefinitionList),
		goldmark.WithExtensions(extension.Table),
	)
	for _, tc := range tests {
		var outputBuffer bytes.Buffer
		err := md.Convert([]byte(tc.input), &outputBuffer)
		if err != nil {
			t.Fatalf("Error converting: %v", err)
		}

		actual := strings.TrimSpace(outputBuffer.String())
		expected := strings.TrimSpace(tc.expected)

		if actual != expected {
			t.Errorf("Unexpected output.\nExpected:\n%s\nActual:\n%s", tc.expected, actual)
		}
	}
}
