package attendess

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
	return &dsc{strings.NewReader(`indexer.token_size=3`)}, nil
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
			want: "mohemmed",
		}, {
			name: "test with some vowels 2",
			args: args{
				s: "ali ahmed",
			},
			want: "ele ehmed",
		}, {
			name: "test with some vowels 2",
			args: args{
				s: "mahmoud",
			},
			want: "mehmood",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substituteVowels(tt.args.s); got != tt.want {
				t.Errorf("substituteVowels() = %v, want %v", got, tt.want)
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
		want []token
	}{
		{
			name: "token on 3 chars",
			args: args{s: "hello", ts: 3},
			want: []token{"hel", "ell", "llo"},
		}, {
			name: "token on 2 chars",
			args: args{s: "hello", ts: 2},
			want: []token{"he", "el", "ll", "lo"},
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

	attendees := []Attendee{
		{DisplayName: "Mohammad Hewedy"},
		{DisplayName: "Asif Ahmed"},
		{DisplayName: "Ibrahim Ahmed"},
		{DisplayName: "Ibrahim Mostafa"},
		{DisplayName: "Saif Ibrahim"},
		{DisplayName: "Asif Ali"},
		{DisplayName: "Rashad Saif"},
		{DisplayName: "Ali Ibrahim"},
		{DisplayName: "Ahmad Altihami"},
	}

	index(attendees)

	result := search("Ibrahim Ahmed")

	fmt.Println(result)
}
