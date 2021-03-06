package gow

import "testing"

type ControllerTest Controller

func (c *ControllerTest) HelloAction() {
	c.Response.RenderString("Hello World")
}

func TestGow(t *testing.T) {
	s := New()
	t.Logf("s:%v\n", s)
	s.SetBaseUrl("/test")
	s.RegisterController("/ct", &ControllerTest{})

	s.Run(":16000")
}

func TestNew(t *testing.T) {

	if s := New(); s == nil || s.config == nil {
		t.Fatalf("New error %s\n", "...")
	}
}
