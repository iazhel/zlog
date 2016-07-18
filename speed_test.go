package zlog

import (
	//	"fmt"
	"testing"
)

func Benchmark_Step_100x1000Info(b *testing.B)    { bench_Step(100, b) }
func Benchmark_StepMap_100x1000Info(b *testing.B) { bench_StepM(100, b) }
func Benchmark_Info_1000(b *testing.B)            { bench_info(1000, b) }

//func Benchmark_Info_100000(b *testing.B)          { bench_info(100000, b) }
func Benchmark_InfoMap_1000(b *testing.B) { bench_infoM(1000, b) }

//func Benchmark_InfoMap_100000(b *testing.B)       { bench_infoM(100000, b) }

func bench_info(n int, b *testing.B) {
	log := NewZLA()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			log.Info("Step 1 -------------------------------------------------")
		}
	}
}

func bench_infoM(n int, b *testing.B) {
	log := NewZL()
	log1 := log.NewStep("step 1")
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			log1.Info("Step 1 -------------------------------------------------")
		}
	}
}

func bench_Step(n int, b *testing.B) {
	logs := make([]*ZLA, n)
	logs[0] = NewZLA()
	for i := 0; i < b.N; i++ {
		for j := 1; j < n; j++ {
			logs[j] = logs[0].NewStep("step")
			for k := 0; k < 1000; k++ {
				logs[j].Info("Step 1 -------------------------------------------------")
			}
		}
		_ = logs[0].GetAllLog()
	}
}
func bench_StepM(n int, b *testing.B) {
	logs := make([]*ZL, n)
	logs[0] = NewZL()
	// bench cycle
	for i := 0; i < b.N; i++ {
		// NewStep cycle
		for j := 1; j < n; j++ {
			logs[j] = logs[0].NewStep("step")
			// Info cycle
			for k := 0; k < 1000; k++ {
				logs[j].Info("Step 1 -------------------------------------------------")
			}
		}

		_ = logs[0].GetAllLog()
	}
}
