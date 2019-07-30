package usage_test

import (
	"net"
	"testing"
	"time"

	"github.com/DeedleFake/cast/usage"
)

func TestInterface(t *testing.T) {
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Errorf("Failed to get interface list: %v", err)
	}

	name := interfaces[0].Name
	for _, inter := range interfaces {
		if inter.Name[0] == 'e' {
			name = inter.Name
		}
	}
	t.Logf("Using interface %q", name)

	inter, err := usage.OpenInterface(name)
	if err != nil {
		t.Errorf("Failed to open %q: %v", name, err)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)

		in, err := inter.In()
		if err != nil {
			t.Errorf("Failed to get incoming bytes: %v", err)
		}

		out, err := inter.Out()
		if err != nil {
			t.Errorf("Failed to get outgoing bytes: %v", err)
		}

		t.Logf("In: %v\n", in)
		t.Logf("Out: %v", out)
	}
}
