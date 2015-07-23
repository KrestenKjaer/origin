package app

import (
	imageapi "github.com/openshift/origin/pkg/image/api"
	templateapi "github.com/openshift/origin/pkg/template/api"
)

// ComponentMatch is a match to a provided component
type ComponentMatch struct {
	Value       string
	Argument    string
	Name        string
	Description string
	Score       float32
	Insecure    bool

	Builder     bool
	Image       *imageapi.DockerImage
	ImageStream *imageapi.ImageStream
	ImageTag    string
	Template    *templateapi.Template

	Meta map[string]string
}

func (m *ComponentMatch) String() string {
	return m.Argument
}

// IsImage returns whether or not the component match is an
// image or image stream
func (m *ComponentMatch) IsImage() bool {
	return m.Template == nil
}

// IsTemplate returns whether or not the component match is
// a template
func (m *ComponentMatch) IsTemplate() bool {
	return m.Template != nil
}

// Exact checks if the ComponentMatch is an exact match
func (m *ComponentMatch) Exact() bool {
	return m.Score == 0.0
}

// ComponentMatches holds multiple ComponentMatch
type ComponentMatches []*ComponentMatch

// Exact returns all ComponentMatch that are an exact match
func (m ComponentMatches) Exact() ComponentMatches {
	exact := ComponentMatches{}
	for _, match := range m {
		if match.Exact() {
			exact = append(exact, match)
		}
	}
	return exact
}

// Inexact returns all ComponentMatch that are not an exact match
func (m ComponentMatches) Inexact() ComponentMatches {
	inexact := ComponentMatches{}
	for _, match := range m {
		if !match.Exact() {
			inexact = append(inexact, match)
		}
	}
	return inexact
}

// ScoredComponentMatches is a set of component matches grouped by score
type ScoredComponentMatches []*ComponentMatch

func (m ScoredComponentMatches) Len() int           { return len(m) }
func (m ScoredComponentMatches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ScoredComponentMatches) Less(i, j int) bool { return m[i].Score < m[j].Score }

// Exact returns all the exact component matches
func (m ScoredComponentMatches) Exact() []*ComponentMatch {
	out := []*ComponentMatch{}
	for _, match := range m {
		if match.Score == 0.0 {
			out = append(out, match)
		}
	}
	return out
}