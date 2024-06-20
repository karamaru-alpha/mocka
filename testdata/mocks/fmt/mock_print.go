// Code generated by mocka. DO NOT EDIT.
package fmt

import (
	"testing"
)

type MockStringer struct {
	t          *testing.T
	MockString *mStringer_String
}

func NewMockStringer(t *testing.T) *MockStringer {
	mockString := &mStringer_String{
		t: t,
	}
	t.Cleanup(func() {
		if mockString.stab != nil {
			return
		}
		for _, expect := range mockString.expects {
			if expect.times != 0 {
				if expect.condition != nil {
					t.Fatalf("Stringer.String conditional expected more %d times. expectNumber: %d", expect.times, expect.number)
				}
				t.Fatalf("Stringer.String expected more %d times. expectNumber %d", expect.times, expect.number)
			}
		}
	})
	return &MockStringer{
		t:          t,
		MockString: mockString,
	}
}

func (m *MockStringer) String() (result0 string) {
	if m.MockString.stab != nil {
		return m.MockString.stab()
	}
	for _, expect := range m.MockString.expects {
		if expect.times == 0 {
			continue
		}
		if expect.condition != nil {
			if !expect.condition() {
				continue
			}
		}
		expect.times--
		return expect.results.result0
	}
	m.t.Fatalf("Stringer.String unexpected call")
	panic("not reachable")
}

type mStringer_StringExpect struct {
	number    int
	condition func() bool
	times     int
	results   struct {
		result0 string
	}
}

type mStringer_String struct {
	t       *testing.T
	expects []*mStringer_StringExpect
	stab    func() string
}

func (m *mStringer_String) Expect() *mStringer_StringExpected {
	if m.stab != nil {
		m.t.Fatalf("Stringer.String cannot expect after stabilize")
	}
	return &mStringer_StringExpected{
		m:      m,
		expect: &mStringer_StringExpect{},
	}
}

func (m *mStringer_String) ConditionalExpect(condition func() bool) *mStringer_StringExpected {
	if m.stab != nil {
		m.t.Fatalf("Stringer.String cannot expect after stabilize")
	}
	return &mStringer_StringExpected{
		m: m,
		expect: &mStringer_StringExpect{
			condition: condition,
		},
	}
}

func (m *mStringer_String) Stabilize(stab func() string) {
	if len(m.expects) > 0 {
		m.t.Fatalf("Stringer.String cannot stabilize after expect")
	}
	m.stab = stab
}

type mStringer_StringExpected struct {
	m      *mStringer_String
	expect *mStringer_StringExpect
}

func (m *mStringer_StringExpected) Times(n int) *mStringer_StringTimesExpected {
	m.expect.times = n
	return &mStringer_StringTimesExpected{
		m:      m.m,
		expect: m.expect,
	}
}

type mStringer_StringTimesExpected struct {
	m      *mStringer_String
	expect *mStringer_StringExpect
}

func (m *mStringer_StringTimesExpected) Return(result0 string) {
	m.expect.results = struct {
		result0 string
	}{
		result0: result0,
	}
	m.expect.number = len(m.m.expects) + 1
	m.m.expects = append(m.m.expects, m.expect)
}
