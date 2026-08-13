// Copyright (c) 2026 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quota

import (
	"math"
	"testing"
)

func TestIsRMB(t *testing.T) {
	cases := []struct {
		unit string
		want bool
	}{
		{"RMB", true},
		{"total_token", false},
		{"", false},
	}
	for _, c := range cases {
		if got := IsRMB(c.unit); got != c.want {
			t.Errorf("IsRMB(%q) = %v, want %v", c.unit, got, c.want)
		}
	}
}

func TestToRedisValue(t *testing.T) {
	cases := []struct {
		quota float64
		unit  string
		want  int64
	}{
		{1.0, UnitRMB, 1e8},
		{0.000001, UnitRMB, 100},
		{0.00000001, UnitRMB, 1},
		{100.0, UnitTotalToken, 100},
		{1.5, UnitTotalToken, 1},
	}
	for _, c := range cases {
		if got := ToRedisValue(c.quota, c.unit); got != c.want {
			t.Errorf("ToRedisValue(%v, %q) = %d, want %d", c.quota, c.unit, got, c.want)
		}
	}
}

func TestFromRedisValue(t *testing.T) {
	cases := []struct {
		value int64
		unit  string
		want  float64
	}{
		{1e8, UnitRMB, 1.0},
		{100, UnitRMB, 0.000001},
		{1, UnitRMB, 0.00000001},
		{100, UnitTotalToken, 100.0},
	}
	for _, c := range cases {
		got := FromRedisValue(c.value, c.unit)
		if math.Abs(got-c.want) > 1e-12 {
			t.Errorf("FromRedisValue(%d, %q) = %v, want %v", c.value, c.unit, got, c.want)
		}
	}
}

func TestPtrToRedisValue(t *testing.T) {
	q := 1.23
	u := UnitRMB
	got := PtrToRedisValue(&q, &u)
	want := int64(1.23 * RmbPrecision)
	if got != want {
		t.Errorf("PtrToRedisValue(%v, %q) = %d, want %d", q, u, got, want)
	}

	if got := PtrToRedisValue(nil, &u); got != 0 {
		t.Errorf("PtrToRedisValue(nil, ...) = %d, want 0", got)
	}
}

func TestRmbToFixedPoint(t *testing.T) {
	cases := []struct {
		yuan float64
		want int64
	}{
		{1.0, 1e8},
		{0.000001, 100},
		{0.00000001, 1},
	}
	for _, c := range cases {
		if got := RmbToFixedPoint(c.yuan); got != c.want {
			t.Errorf("RmbToFixedPoint(%v) = %d, want %d", c.yuan, got, c.want)
		}
	}
}

func TestFixedPointToRmb(t *testing.T) {
	cases := []struct {
		value int64
		want  float64
	}{
		{1e8, 1.0},
		{100, 0.000001},
		{1, 0.00000001},
	}
	for _, c := range cases {
		got := FixedPointToRmb(c.value)
		if math.Abs(got-c.want) > 1e-12 {
			t.Errorf("FixedPointToRmb(%d) = %v, want %v", c.value, got, c.want)
		}
	}
}
