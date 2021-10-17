package train

import "testing"

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "default",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := List(); (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
