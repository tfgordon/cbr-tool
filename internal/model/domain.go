// Copyright © 2015-16 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package model

// See https://en.wikipedia.org/wiki/Multidimensional_scaling
// for ideas about how to visualize similarity measures

import (
	// "fmt"
	"math"
)

type Dimension struct {
	Description string
	Factors     []string           // ordered set of factor ids
	Default     string             // id of the default factor
	rank        map[string]float64 // The rank of each factor in the dimension
}

type Domain struct {
	Id          string `json:"_id"`
	Revision    string `json:"_rev"`
	Owner       string // userid of the owner of the domain model
	Issue       string // name of the issue modelled in this domain model
	Description string
	Options     []string          // statement ids of the options of the issue
	Statements  map[string]string // factor id to natural language text
	Dimensions  map[string]*Dimension
	index       map[string]string // from factor id to dimension id
	initialized bool
}

type CaseFactors map[string]string // dimension id to factor id

// Initialize the domain model by computing its factor index and ranking
// the factors of the dimensions
func (m *Domain) init() {
	// compute the factor index
	m.index = make(map[string]string)
	for k, d := range m.Dimensions {
		for _, f := range d.Factors {
			m.index[f] = k
		}
	}
	// rank the factors of each dimension
	for _, d := range m.Dimensions {
		d.rank = make(map[string]float64)
		for i, f := range d.Factors {
			d.rank[f] = float64(i)
		}
	}
	m.initialized = true
}

// Returns the statement of the factor with the given id.
// If no statement is defined for the factor, the id is used.
func (m *Domain) FactorStatement(id string) string {
	s, ok := m.Statements[id]
	if ok {
		return s
	} else {
		return id
	}
}

// Returns the dimension id of the factor with the given id.
// The boolean flag returned is false if there is no dimension for the
// given factor in the domain model
func (m *Domain) FactorDimension(factorId string) (string, bool) {
	if !m.initialized {
		m.init()
	}
	d, ok := m.index[factorId]
	if !ok {
		return "", false
	} else {
		return d, true
	}
}

// Checks whether the statement id is an option
// of the model.
func (m *Domain) IsOption(sid string) bool {
	result := false
	for _, id := range m.Options {
		if sid == id {
			result = true
		}
	}
	return result
}

// Case similarity is measured using the geometric mean of the ordinal
// distance between the factors of the dimensions
func (m *Domain) GeometricMeanSimilarity(c1 CaseFactors, c2 CaseFactors) float64 {

	if !m.initialized {
		m.init()
	}

	deltas := []float64{}
	for k, d := range m.Dimensions {
		f1, ok1 := c1[k]
		if !ok1 {
			f1 = d.Default
		}
		f2, ok2 := c2[k]
		if !ok2 {
			f2 = d.Default
		}

		deltas = append(deltas, math.Abs(d.rank[f1]-d.rank[f2]))
	}

	// compute the geometric mean of the deltas
	var product float64 = 1.0
	for _, delta := range deltas {
		// add 1 to the differences to avoid multiplying by zero
		product = product * (delta + 1)
	}
	n := len(deltas)
	mean := math.Pow(product, 1.0/float64(n))

	return 1.0 / mean
}

func (m *Domain) ArithmetricMeanSimilarity(c1 CaseFactors, c2 CaseFactors) float64 {
	if !m.initialized {
		m.init()
	}

	// deltas is a list of the differences of the *normalized* factors of
	// each dimension

	deltas := []float64{}

	normalize := func(f string, d *Dimension) float64 {
		n := float64(len(d.Factors))
		i := d.rank[f]
		return (i * float64(100)) / n
	}

	for k, d := range m.Dimensions {
		f1, ok1 := c1[k]
		if !ok1 {
			f1 = d.Default
		}
		f2, ok2 := c2[k]
		if !ok2 {
			f2 = d.Default
		}

		deltas = append(deltas, math.Abs(normalize(f1, d)-normalize(f2, d)))
	}

	// compute the arithmetric mean of the deltas
	var sum float64 = 0.0
	for _, delta := range deltas {
		sum = sum + delta
	}
	n := len(deltas)
	mean := sum / float64(n)

	return (100.0 - mean) / 100.0
}
