package index

import (
	"fmt"
	"github.com/mhewedy/go-conf"
	"io"
	"reflect"
	"strings"
	"testing"
)

type ds struct {
}
type dsc struct {
	io.Reader
}

func (d dsc) Close() error {
	return nil
}

func (d ds) Read() (io.ReadCloser, error) {
	return &dsc{strings.NewReader(`indexer.token_algo.token_size=4`)}, nil
}
func init() {

	conf.DefaultSource = ds{}
}

func Test_removeVowels(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test with no vowels",
			args: args{
				s: "mhmd",
			},
			want: "mhmd",
		},
		{
			name: "test with some vowels",
			args: args{
				s: "mohammed",
			},
			want: "mhmmd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeVowels(tt.args.s); got != tt.want {
				t.Errorf("removeVowels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tokenize(t *testing.T) {
	type args struct {
		s  string
		ts int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "token on 3 chars",
			args: args{s: "hello", ts: 3},
			want: []string{"hel", "ell", "llo"},
		}, {
			name: "token on 2 chars",
			args: args{s: "hello", ts: 2},
			want: []string{"he", "el", "ll", "lo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.args.s, tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Search(t *testing.T) {

	type myStruct struct {
		name string
	}
	Index([]Input{
		{Field: "Asif Ahmed", Ref: myStruct{name: "Asif Ahmed"}},
		{Field: "Saif Ibrahim", Ref: myStruct{name: "Saif Ibrahim"}},
		{Field: "Asif Ali", Ref: myStruct{name: "Asif Ali"}},
		{Field: "Rashad Saif", Ref: myStruct{name: "Rashad Saif"}},
	})

	fmt.Println(index)

	result := Search("Rashad saif")

	fmt.Println(result)
}
