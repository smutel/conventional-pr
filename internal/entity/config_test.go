package entity

import (
	"os"
	"reflect"
	"testing"

	"github.com/Namchee/conventional-pr/internal/constants"
)

func TestReadConfig(t *testing.T) {
	type expected struct {
		config *Config
		err    error
	}
	tests := []struct {
		name    string
		mocks   map[string]string
		want    expected
		wantErr bool
	}{
		{
			name: "should read config correctly",
			mocks: map[string]string{
				"INPUT_ACCESS_TOKEN":         "foo_bar",
				"INPUT_DRAFT":                "false",
				"INPUT_ISSUE":                "true",
				"INPUT_BOT":                  "false",
				"INPUT_MAXIMUM_FILE_CHANGES": "11",
				"INPUT_VERIFIED_COMMITS":     "true",
				"INPUT_IGNORED_USERS":        "Namchee, snyk-bot",
				"INPUT_REPORT":               "false",
			},
			want: expected{
				config: &Config{
					Token:        "foo_bar",
					Draft:        false,
					Issue:        true,
					Bot:          false,
					FileChanges:  11,
					Verified:     true,
					IgnoredUsers: []string{"Namchee", "snyk-bot"},
					Report:       false,
				},
				err: nil,
			},
			wantErr: false,
		},
		{
			name:  "should throw an error when token is empty",
			mocks: map[string]string{},
			want: expected{
				config: nil,
				err:    constants.ErrMissingToken,
			},
			wantErr: true,
		},
		{
			name: "should throw an error when fileChanges is negative",
			mocks: map[string]string{
				"INPUT_ACCESS_TOKEN":         "foo",
				"INPUT_MAXIMUM_FILE_CHANGES": "-1",
			},
			want: expected{
				config: nil,
				err:    constants.ErrNegativeFileChange,
			},
			wantErr: true,
		},
		{
			name: "should throw an error when title pattern is invalid",
			mocks: map[string]string{
				"INPUT_TITLE_PATTERN": "[",
			},
			want: expected{
				config: nil,
				err:    constants.ErrInvalidTitlePattern,
			},
			wantErr: true,
		},
		{
			name: "should throw an error when commit pattern is invalid",
			mocks: map[string]string{
				"INPUT_TITLE_PATTERN":  "a",
				"INPUT_COMMIT_PATTERN": "[",
			},
			want: expected{
				config: nil,
				err:    constants.ErrInvalidCommitPattern,
			},
			wantErr: true,
		},
		{
			name: "should throw an error when branch pattern is invalid",
			mocks: map[string]string{
				"INPUT_TITLE_PATTERN":  "a",
				"INPUT_COMMIT_PATTERN": "a",
				"INPUT_BRANCH_PATTERN": "[",
			},
			want: expected{
				config: nil,
				err:    constants.ErrInvalidBranchPattern,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.mocks {
				os.Setenv(key, val)
				defer os.Unsetenv(key)
			}

			got, err := ReadConfig()

			if tc.wantErr && err == nil {
				t.Fatalf("ReadConfig() err = %v, wantErr = %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(got, tc.want.config) {
				t.Fatalf("ReadConfig() = %v, want = %v", got, tc.want)
			}
		})
	}
}
