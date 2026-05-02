package telegram

import "testing"

func TestReplyForCommand(t *testing.T) {
	tests := map[string]string{
		"/start":   StartText,
		"/help":    HelpText,
		"/profile": ProfileStartText,
		"hello":    HelpText,
	}

	for command, want := range tests {
		t.Run(command, func(t *testing.T) {
			got := ReplyForCommand(command)
			if got != want {
				t.Fatalf("ReplyForCommand(%q) = %q, want %q", command, got, want)
			}
		})
	}
}
