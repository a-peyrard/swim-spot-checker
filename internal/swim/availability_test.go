package swim

import (
	"reflect"
	"testing"
)

func Test_sanitizeAndParseResponse(t *testing.T) {
	type args struct {
		rawResponse string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse map[string]any
		wantErr      bool
	}{
		{
			name: "it should parse well formed JSON",
			args: args{
				rawResponse: "{\"available\": false, \"explanation\": \"The new version removed the specific times (3:00 & 6:00) previously listed, making it unclear if any spots are currently available.  No new session types or times were added.\"}",
			},
			wantResponse: map[string]any{
				"available":   false,
				"explanation": "The new version removed the specific times (3:00 & 6:00) previously listed, making it unclear if any spots are currently available.  No new session types or times were added.",
			},
			wantErr: false,
		},
		{
			name: "it should parse JSON markdown formatted",
			args: args{
				rawResponse: "```json\n{\n  \"available\": false,\n  \"explanation\": \"The new version removed the specific times (3:00 & 6:00) previously listed, making it unclear if any spots are currently available.  No new session types or times were added.\"\n}\n```",
			},
			wantResponse: map[string]any{
				"available":   false,
				"explanation": "The new version removed the specific times (3:00 & 6:00) previously listed, making it unclear if any spots are currently available.  No new session types or times were added.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := sanitizeAndParseResponse(tt.args.rawResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("sanitizeAndParseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("sanitizeAndParseResponse() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
