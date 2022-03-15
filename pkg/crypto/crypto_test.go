package crypto

import (
	"testing"
)

func TestEncryptAES(t *testing.T) {
	type args struct {
		data string
		key  string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				data: `{
					"email": "admin@gmail.com",
					"password": "123456"
				}`,
				key: "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
			},
			want: "PfBYQ/35Gnu+9K/TjtT4H6wB9/XP1m+ssny6anIjJdD6q3rlJ9MZc0b8KB/rcF5HnxbIV2bKcRT5RJ1ynS4mig==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptAES(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncryptAES() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecryptAES(t *testing.T) {
	type args struct {
		data string
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				data: "PfBYQ/35Gnu+9K/TjtT4H6wB9/XP1m+ssny6anIjJdD6q3rlJ9MZc0b8KB/rcF5HnxbIV2bKcRT5RJ1ynS4mig==",
				key:  "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
			},
			want: `{
				"email": "admin@gmail.com",
				"password": "123456"
			}`,
		},
		{
			name: "test",
			args: args{
				data: "9ezlf/FMUpxWCrZZxtHKAtVHxdoJJF2Vvypg/XpgvdA+KKoMYzzIEv2YBkqXF/aIrK8sVfsciJkMosic0BXfCcasINOIcfjgQqFYNq+qnns=",
				key:  "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
			},
			want: `{"code":3, "message":"failed to find user: mongo: no documents in result"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptAES(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecryptAES() got = %v, want %v", got, tt.want)
			}
		})
	}
}
