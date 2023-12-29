package atlassianwikirenderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

type testCase struct {
	input    string
	expected string
}

func TestTaskCheckbox(t *testing.T) {
	runTest(t, testCase{
		input:    "- [ ] Task 1\n- [x] Task 2 (completed)\n- [ ] Task 3",
		expected: "*  [ ] Task 1\n*  [x] Task 2 (completed)\n*  [ ] Task 3",
	})
}

func TestHeadingRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "# Heading 1",
		expected: "h1.Heading 1",
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
		input:    "```This is a fenced code block```",
		expected: "{{This is a fenced code block}}",
	})
}

func TestParagraphRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "This is a paragraph. It can contain text,\nline breaks, and various formatting options.\n\nYou can separate paragraphs by leaving a blank line between them.",
		expected: "This is a paragraph. It can contain text,\n\nline breaks, and various formatting options.\n\nYou can separate paragraphs by leaving a blank line between them.",
	})
}

func TestCodeBlockRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "\t`Code Block`",
		expected: "{code}{code}",
	})
}

func TestCodeSpanRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "This is an `inline code span`",
		expected: "This is an {{inline code span}}",
	})
}

func TestEmphasisRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "*Emphasised text*",
		expected: "_Emphasised text_",
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
		input:    "~~Strikethrough~~",
		expected: "~~Strikethrough~~",
	})
}

func TestListRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "- List",
		expected: "*  List",
	})
}

func TestListItemRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "- List\n  - List item",
		expected: "*  List\n**  List item",
	})
}

func TestTableRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "| Month    | Assignee | Backup |\n| -------- | -------- | ------ |\n| January  | Dave     | Steve  |\n| February | Gregg    | Karen  |\n| March    | Diane    | Jorge  |",
		expected: "| Month    | Assignee | Backup |\n\n| -------- | -------- | ------ |\n\n| January  | Dave     | Steve  |\n\n| February | Gregg    | Karen  |\n\n| March    | Diane    | Jorge  |",
	})
}

func TestImageRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "![JFrog logo](https://speedmedia.jfrog.com/08612fe1-9391-4cf3-ac1a-6dd49c36b276/https://media.jfrog.com/wp-content/uploads/2021/12/29113553/jfrog-logo-2022.svg)",
		expected: "!https://speedmedia.jfrog.com/08612fe1-9391-4cf3-ac1a-6dd49c36b276/https://media.jfrog.com/wp-content/uploads/2021/12/29113553/jfrog-logo-2022.svg!",
	})
}

func TestFootnoteRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "[^1]: This is the footnote.",
		expected: "[^1]: This is the footnote.",
	})
}

func TestFootNoteLinkRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "Here is some text with a footnote[^1].\n[^1]: This is the footnote content.",
		expected: "Here is some text with a footnote[^1].\n\n[^1]: This is the footnote content.",
	})
}

func TestFootNoteListRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "[^1]: This is the content of the first footnote.\n[^2]: This is the content of the second footnote.\n",
		expected: "[^1]: This is the content of the first footnote.\n\n[^2]: This is the content of the second footnote.\n",
	})
}

func TestDefinitionTermDescriptionRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "Definition Term\n:   This is the description for Term.",
		expected: "Definition Term\n\n:   This is the description for Term.",
	})
}

func TestHTMLBlockRendering(t *testing.T) {
	runTest(t, testCase{
		input:    "```html\n<div>\n  <p>This is an HTML block.</p>\n</div>",
		expected: "{code:html}<div>\n  <p>This is an HTML block.</p>\n</div>{code}",
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

func runTest(t *testing.T, tc testCase) {
	var outputBuffer bytes.Buffer

	renderer := New()

	md := goldmark.New(
		goldmark.WithRenderer(renderer),
	)

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
