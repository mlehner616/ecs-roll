package regex

import "testing"

func TestParseContainerInstanceId(t *testing.T) {
	cases := []struct {
		input      string
		shouldPass bool
	}{
		{"arn:aws:ecs:us-west-2:382937539128:container-instance/8406952b-095f-44ed-b7e6-b52cbe13eb6e", true},
		{"arn:aws:ecs:us-east-1:368129391883:container-instance/d67bd798-5f88-4142-8944-4f4c28eda446", true},
		{"d67bd798-5f88-4142-8944-4f4c28eda446", false},
		{"arn:aws:ecs:us-west-2:930485739301:container-instance/8406952b-095f-4ed-b7e6-b52cbe13eb6e", false},
		{"arn:aws:ecs:us-west-2:368129391883:container-instancefd101285-6a87-4589-81b3-8300ca3f51cc", false},
	}

	for _, c := range cases {
		out, err := parseContainerInstanceId(c.input)

		// The test case should have passed but an error occurred
		if c.shouldPass && err != nil {
			t.Errorf("error occurred on test case that '%s' should have passed: %v", c.input, err)
		}

		// Did not received any output when we should have
		if c.shouldPass && len(out) == 0 {
			t.Errorf("received zero output but expected non-zero output for '%s'", c.input)
		}

		// The test case should have failed but no error was returned
		if !c.shouldPass && err == nil {
			t.Errorf("purposely failed test case '%s' should have resulted an a returned error", c.input)
		}
	}
}

func BenchmarkParseContainerInstanceId(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parseContainerInstanceId("arn:aws:ecs:us-west-2:368129391883:container-instancefd101285-6a87-4589-81b3-8300ca3f51cc")
	}
}
