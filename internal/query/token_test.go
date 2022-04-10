package query

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

func TestNgram(t *testing.T) {
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
		// TODO: Add test cases.
		{"test1", args{"北京北京", 2}, []Tokenization{
			{[]rune("北京"), 0},
			{[]rune("京北"), 1},
			{[]rune("北京"), 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ngram(tt.args.content, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ngram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ngram() = %v, want %v", got, tt.want)
			}
		})
	}
}
