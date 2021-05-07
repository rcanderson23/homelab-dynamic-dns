package dnsproviders

import "testing"

func TestIsCurrentEndpoint(t *testing.T) {
	type args struct {
		host    string
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "correct google dnsproviders",
			args: args{
				host:    "dnsproviders.google",
				address: "8.8.8.8",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "incorrect google dnsproviders",
			args: args{
				host:    "dnsproviders.google",
				address: "1.1.1.1",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsCurrentEndpoint(tt.args.host, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsCurrentEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsCurrentEndpoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}
