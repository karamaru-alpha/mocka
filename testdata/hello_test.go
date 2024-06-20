package testdata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/mocka/testdata/mocks/fmt"
	"github.com/karamaru-alpha/mocka/testdata/mocks/github.com/karamaru-alpha/mocka/testdata"
)

type mocks struct {
	stringer *fmt.MockStringer
	human    *testdata.MockHuman
}

func newMocks(t *testing.T) *mocks {
	return &mocks{
		stringer: fmt.NewMockStringer(t),
		human:    testdata.NewMockHuman(t),
	}
}

func TestHello(t *testing.T) {
	t.Run("pass: expect", func(t *testing.T) {
		m := newMocks(t)
		m.stringer.MockString.Expect().Times(1).Return("Alice")
		m.human.MockSay.Expect("Alice").Times(1).Return("Hello, Alice")

		result := Hello(m.stringer, m.human)
		assert.Equal(t, "Hello, Alice", result)
	})
	t.Run("pass: conditional expect", func(t *testing.T) {
		m := newMocks(t)
		m.stringer.MockString.ConditionalExpect(func() bool {
			return true
		}).Times(1).Return("Alice")
		m.human.MockSay.ConditionalExpect(func(name string) bool {
			return name == "Alice" || name == "Bob"
		}).Times(1).Return("Hello, Alice or Bob")

		result := Hello(m.stringer, m.human)
		assert.Equal(t, "Hello, Alice or Bob", result)
	})
	t.Run("pass: stabilize", func(t *testing.T) {
		m := newMocks(t)
		m.stringer.MockString.Stabilize(func() string {
			return "Alice"
		})
		m.human.MockSay.Stabilize(func(name string) string {
			return "Hello, " + name
		})

		result := Hello(m.stringer, m.human)
		assert.Equal(t, "Hello, Alice", result)
	})

	t.Run("fail: unexpected call", func(t *testing.T) {
		t.Skip()
		m := newMocks(t)
		m.stringer.MockString.Expect().Times(1).Return("Alice")
		// m.human.MockSay.Expect("Alice").Times(1).Return("Hello, Alice")

		Hello(m.stringer, m.human) // fail: Human.Say unexpected call. name: Alice
	})
	t.Run("fail: over expected", func(t *testing.T) {
		t.Skip()
		m := newMocks(t)
		m.stringer.MockString.Expect().Times(1).Return("Alice")
		m.human.MockSay.Expect("Alice").Times(2).Return("Hello, Alice")

		Hello(m.stringer, m.human) // fail: Human.Say expected more 1 times. args: {name:Alice}, expectNumber: 1
	})
	t.Run("fail: collision　expect and stab", func(t *testing.T) {
		t.Skip()
		m := newMocks(t)
		m.stringer.MockString.Expect().Times(1).Return("Alice")
		m.human.MockSay.Expect("Alice").Times(1).Return("Hello, Alice")
		m.human.MockSay.Stabilize(func(name string) string {
			return "Hello, " + name
		})

		Hello(m.stringer, m.human) // fail: Human.Say cannot stabilize after expect
	})
}
