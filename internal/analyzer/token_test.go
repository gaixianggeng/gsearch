package analyzer

import (
	"reflect"
	"testing"
)

func Test_ignoredChar(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test1", args{"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"}, ""},
		{"test2", args{"test  test"}, "testtest"},
		{"test2", args{"test!test"}, "testtest"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ignoredChar(tt.args.str); got != tt.want {
				t.Errorf("ignoredChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNGram(t *testing.T) {
	type args struct {
		content string
		n       int32
	}
	tests := []struct {
		name    string
		args    args
		want    []Tokenization
		wantErr bool
	}{
		{
			"test1",
			args{"北京北京", 2},
			[]Tokenization{
				{("北京"), 0},
				{("京北"), 1},
				{("北京"), 2}},
			false},
		{
			"test2",
			args{"北京北京", 3},
			[]Tokenization{
				{("北京北"), 0},
				{("京北京"), 1},
				{("北京"), 2}},
			true},
		{
			"test3",
			args{"北京北京", 4},
			[]Tokenization{
				{("北京北京"), 0},
				{("北京"), 1}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NGram(tt.args.content, tt.args.n)
			t.Logf("got:%v", got)
			if err != nil {
				t.Errorf("NGram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.want) == tt.wantErr {
				t.Errorf("NGram() = %v, want %v", got, tt.want)
			}
		})
	}
}
