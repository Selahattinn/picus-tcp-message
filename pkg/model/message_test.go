package model

import (
	"strconv"
	"testing"
)

func TestMessage_ToString(t *testing.T) {
	type fields struct {
		ID   int64
		From string
		To   string
		Text string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: " String format correct", fields: fields{ID: 1, From: "Test_From", To: "Test_To", Text: "Test Text"}, want: "ID: " + strconv.FormatInt(1, 10) + "\n\tFrom: " + "Test_From" + "\n\tTo: " + "Test_To" + "\n\tmessage: " + "Test Text" + "\n"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Message{
				ID:   tt.fields.ID,
				From: tt.fields.From,
				To:   tt.fields.To,
				Text: tt.fields.Text,
			}
			if got := m.ToString(); got != tt.want {
				t.Errorf("Message.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
