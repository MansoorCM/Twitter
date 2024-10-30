package auth

import "testing"

func TestCheckPassordHash(t *testing.T) {
	pass1 := "a new password 1"
	pass2 := "a second and longer password 2"
	hash1, _ := HashPassword(pass1)
	hash2, _ := HashPassword(pass2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{name: "correct pass and hash 1",
			password: pass1,
			hash:     hash1,
			wantErr:  false,
		},
		{name: "correct pass and hash 2",
			password: pass2,
			hash:     hash2,
			wantErr:  false,
		},
		{name: "wrong pass and hash 1",
			password: pass1,
			hash:     hash2,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
