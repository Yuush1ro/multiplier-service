package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"net/http"
	"sync/atomic"
)

type Response struct {
	Result float64 `json:"result"`
}

type Generator struct {
	rtp     float64
	phi     float64
	counter uint64
}

func NewGenerator(rtp float64) *Generator {
	phi := (math.Sqrt(5.0) - 1.0) / 2.0 // 0.618...
	return &Generator{
		rtp: rtp,
		phi: phi,
	}
}

func (g *Generator) GetMultiplier() float64 {
	i := atomic.AddUint64(&g.counter, 1)

	f := float64(i) * g.phi
	_, frac := math.Modf(f)
	u := frac

	if u <= 0 {
		u = math.SmallestNonzeroFloat64
	}

	m := g.rtp / u

	if m < 1.0 {
		m = 1.0
	}
	if m > 10000.0 {
		m = 10000.0
	}
	return m
}

func main() {
	rtp := flag.Float64("rtp", 1.0, "target RTP in (0, 1]")
	flag.Parse()

	if *rtp <= 0 || *rtp > 1.0 {
		log.Fatalf("Invalid RTP value: must be in (0,1], got %f", *rtp)
	}

	gen := NewGenerator(*rtp)

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		m := gen.GetMultiplier()
		resp := Response{Result: m}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Server started on :64333 with RTP=%.6f", *rtp)
	log.Fatal(http.ListenAndServe(":64333", nil))
}
