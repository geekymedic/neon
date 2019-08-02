package prometheus

import (
	"context"
	"fmt"
	"github.com/geekymedic/neon/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSys(t *testing.T) {
	version.PRONAME = "boss-system-bff-admin"
	idx := strings.Index(version.PRONAME, "system-")
	s := len("system-")
	if idx > 0 {
		idx = idx + s
		sysLabs.subSystem = strings.ReplaceAll(version.PRONAME[0:idx-1], "-", "_")
		sysLabs.nameSpace = strings.ReplaceAll(version.PRONAME[idx:], "-", "_")
	}
	t.Log(sysLabs)
}

func TestNewCounter(t *testing.T) {
	sysLabs.subSystem = "admin"
	sysLabs.nameSpace = "boss"
	t.Run("Newgague", func(t *testing.T) {
		name := fmt.Sprintf("%d", time.Now().UnixNano())
		gauge, err := NewGague(name)
		assert.Nil(t, err)
		assert.NotNil(t, gauge)
		gauge, err = NewGague(name)
		assert.NotNil(t, err)
		assert.Nil(t, gauge)
	})

	t.Run("", func(t *testing.T) {
		name := fmt.Sprintf("%d", time.Now().UnixNano())
		gauge, err := NewCounterWithLabelNames(name, "foo")
		assert.Nil(t, err)
		assert.NotNil(t, gauge)
		gauge, err = NewCounterWithLabelNames(name, "foo")
		assert.NotNil(t, err)
		assert.Nil(t, gauge)
	})

	t.Run("MustGague", func(t *testing.T) {
		name := fmt.Sprintf("%d", time.Now().UnixNano())
		assert.NotNil(t, MustGague(name))
		defer func() {
			err := recover()
			assert.NotNil(t, err)
		}()
		MustGague(name)
	})

	t.Run("MustCounterWithLabelNames", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			assert.NotNil(t, MustCounterWithLabelNames(fmt.Sprintf("%d", time.Now().UnixNano())))
			assert.NotNil(t, MustCounterWithLabelNames(fmt.Sprintf("%d", time.Now().UnixNano())))

			cc := MustCounterWithLabelNames(fmt.Sprintf("%d", time.Now().UnixNano()), "method", "code")
			cc.Add(1)
		})
		t.Run("register same name and labels", func(t *testing.T) {
			name := fmt.Sprintf("%d", time.Now().UnixNano())
			assert.NotNil(t, MustCounterWithLabelNames(name, "b"))
			defer func() {
				err := recover()
				assert.NotNil(t, err)
			}()
			MustCounterWithLabelNames(name, "b")
		})

		t.Run("register same name, two different labels", func(t *testing.T) {
			name := fmt.Sprintf("%d", time.Now().UnixNano())
			assert.NotNil(t, MustCounterWithLabelNames(name, "b"))
			defer func() {
				err := recover()
				assert.NotNil(t, err)
			}()
			MustCounterWithLabelNames(name, "c")
		})
	})

	t.Run("http-metrics", func(t *testing.T) {
		name := fmt.Sprintf("%d", time.Now().UnixNano())
		c := MustCounter(name)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			for {
				if ctx.Err() != nil {
					return
				}
				c.Add(1)
				time.Sleep(time.Millisecond * 100)
			}
		}()

		name1 := fmt.Sprintf("%d", time.Now().UnixNano())
		fmt.Println("name", name1)
		cc := MustCounterWithLabelNames(name1, "method", "code")
		for i := 0; i < 10; i++ {
			go func() {
				for {
					if ctx.Err() != nil {
						return
					}
					cc.Inc()
					cc.With("GET", "405").Add(2)
					time.Sleep(time.Millisecond * 100)
				}
			}()
		}
		go func() {
			time.Sleep(time.Second * 20)
			cancel()
			os.Exit(0)
		}()
		err := StartMetricsServer(":18090", "/metrics")
		assert.Nil(t, err)
	})
}

func TestNewSummary(t *testing.T) {
	t.Run("", func(t *testing.T) {
		summary := MustSummaryWithLabelNames("pond_temperature_celsius", map[float64]float64{0.5: 0.005, 0.9: 0.01, 0.99: 0.001}, "method")
		for i := 0; i < 10; i++ {
			value := 30 + math.Floor(120*math.Sin(float64(i)*0.1))/10
			summary.With("get").Observe(value)
			fmt.Println(value)
		}
		StartMetricsServer(":18090", "/metrics")
	})

	t.Run("", func(t *testing.T) {
		temps := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "pond_temperature_celsius",
			//Help:       "The temperature of the frog pond.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"metho"})
		prometheus.Register(temps)
		// Simulate some observations.
		for i := 0; i < 10; i++ {
			value := 30 + math.Floor(120*math.Sin(float64(i)*0.1))/10
			temps.WithLabelValues("get").Observe(value)
			fmt.Println(value)
		}

		// Just for demonstration, let's check the state of the summary by
		// (ab)using its Write method (which is usually only used by Prometheus
		// internally).

		StartMetricsServer(":18090", "/metrics")
	})

	t.Run("", func(t *testing.T) {
		metrics, err := NewSummary("request_summary", map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001})
		assert.Nil(t, err)
		assert.NotNil(t, metrics)
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		go func() {
			tm := time.Now()
			count := float64(0)
			for {
				if ctx.Err() != nil {
					return
				}
				time.Sleep(time.Millisecond * 100)
				metrics.Observe(count / float64(time.Since(tm)))
			}
		}()
		go func() {
			time.Sleep(time.Second * 20)
			cancel()
			os.Exit(0)
		}()
		err = StartMetricsServer(":18090", "/metrics")
		assert.Nil(t, err)
		cancel()
	})

}
