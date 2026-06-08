package template

import (
	"strings"
	"testing"
)

func TestRenderInjectsTagsIntoOldTemplate(t *testing.T) {
	old := `---
title: "{title}"
tags: []
---`
	result := Render(old, "lsof", "linux", "linux")
	if strings.Contains(result, "tags: []") {
		t.Errorf("tags: [] was not replaced, got:\n%s", result)
	}
	if !strings.Contains(result, "tags: [\"linux\"") {
		t.Errorf("expected auto-injected tags, got:\n%s", result)
	}
	t.Logf("Result:\n%s", result)
}

func TestRenderNewTemplate(t *testing.T) {
	newTmpl := `---
title: "{title}"
tags: {tags}
---`
	result := Render(newTmpl, "lsof", "linux", "linux")
	if strings.Contains(result, "{tags}") {
		t.Errorf("{tags} was not replaced, got:\n%s", result)
	}
	t.Logf("Result:\n%s", result)
}
