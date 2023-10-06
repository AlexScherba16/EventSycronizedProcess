package main

import (
	"eventsyncprocess/internal/pkg/event_bus"
	"eventsyncprocess/process/prepare_qapi_process"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb/v8"
	"sync"
	"time"
)

func main() {
	bus := event_bus.NewEventBus()

	progress := mpb.New(
		mpb.WithWidth(80),
		mpb.WithOutput(color.Output),
		mpb.WithAutoRefresh())

	wg := &sync.WaitGroup{}

	qapiProcess := prepare_qapi_process.NewProcess(progress, &bus)
	//remoteProcess := remote_process.NewProcess(progress, &bus)
	wg.Add(1)
	go func() {
		defer wg.Done()
		qapiProcess.Run()
	}()
	//go func() {
	//	defer wg.Done()
	//	remoteProcess.Run()
	//}()

	wg.Wait()
	progress.Wait()

	//numBars := 2
	//// to support color in Windows following both options are required
	//p := mpb.New(
	//	mpb.WithOutput(color.Output),
	//	mpb.WithAutoRefresh(),
	//)
	//
	//red, green := color.New(color.FgRed), color.New(color.FgGreen)
	//
	//for i := 0; i < numBars; i++ {
	//	task := fmt.Sprintf("Task#%02d:", i)
	//	queue := make([]*mpb.Bar, 2)
	//	queue[0] = p.AddBar(50, //rand.Int63n(201)+100,
	//		mpb.BarFillerClearOnComplete(),
	//		mpb.PrependDecorators(
	//			decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
	//			decor.Name("downloading", decor.WCSyncSpaceR),
	//			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
	//		),
	//		mpb.AppendDecorators(
	//			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
	//		),
	//	)
	//	queue[1] = p.AddBar(23, //rand.Int63n(101)+100,
	//		mpb.BarQueueAfter(queue[0]), // this bar is queued
	//		mpb.BarFillerClearOnComplete(),
	//		mpb.PrependDecorators(
	//			decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
	//			decor.OnCompleteMeta(
	//				decor.OnComplete(
	//					decor.Meta(decor.Name("installing", decor.WCSyncSpaceR), toMetaFunc(red)),
	//					"done!",
	//				),
	//				toMetaFunc(green),
	//			),
	//			decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_MMSS, 0, decor.WCSyncWidth), ""),
	//		),
	//		mpb.AppendDecorators(
	//			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
	//		),
	//	)
	//
	//	go func() {
	//		for _, b := range queue {
	//			complete(b)
	//		}
	//	}()
	//}
	//
	//p.Wait()
}

func complete(bar *mpb.Bar) {
	//max := 200 * time.Millisecond
	for !bar.Completed() {

		// start variable is solely for EWMA calculation
		// EWMA's unit of measure is an iteration's duration
		//start := time.Now()
		//time.Sleep(time.Duration(rand.Intn(10)+1) * max / 10)
		// we need to call EwmaIncrement to fulfill ewma decorator's contract
		//bar.EwmaIncrInt64(rand.Int63n(5)+1, time.Since(start))

		time.Sleep(500 * time.Millisecond) //time.Duration(500+rand.Intn(1500)) * time.Millisecond) // случайная задержка для демонстрации
		bar.IncrBy(1)
	}
}

func toMetaFunc(c *color.Color) func(string) string {
	return func(s string) string {
		return c.Sprint(s)
	}
}
