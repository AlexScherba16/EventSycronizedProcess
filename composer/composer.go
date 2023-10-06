package composer

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"time"
)

type IComposer interface {
	RunApplication()
	StopApplication()
}

type composer struct {
}

func NewComposer() IComposer {
	return &composer{}
}

func (c *composer) RunApplication() {
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithOutput(color.Output),
		mpb.WithAutoRefresh())

	numBars := 2
	for i := 0; i < numBars; i++ {

		total := 10
		taskFirst := fmt.Sprintf("[CreateQapiStep]")
		taskSecond := fmt.Sprintf("[WaitQapiReadyStep]")
		queue := make([]*mpb.Bar, 2)

		queue[0] = p.AddBar(int64(total),
			mpb.PrependDecorators(
				decor.CountersNoUnit("%d / %d", decor.WC{W: 10}, decor.WCSyncWidth),
			),
			mpb.AppendDecorators(
				decor.Name(taskFirst, decor.WC{W: len(taskFirst) + 5, C: decor.DidentRight}),
				decor.Elapsed(decor.ET_STYLE_GO),
			),
			mpb.BarFillerClearOnComplete(),
		)
		queue[1] = p.AddBar(int64(total)+6,
			mpb.BarQueueAfter(queue[0]), // this bar is queued
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
				//decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_MMSS, 0, decor.WCSyncWidth), ""),
			),
			mpb.AppendDecorators(
				decor.Name(taskSecond, decor.WC{W: len(taskSecond) + 5, C: decor.DidentRight}),
				decor.Elapsed(decor.ET_STYLE_GO),
				//decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
			),
		)

		go func() {
			for _, b := range queue {
				complete(b)
			}
		}()
	}
	//
	//for i := 0; i < total; i++ {
	//	time.Sleep(1300 * time.Millisecond) //time.Duration(500+rand.Intn(1500)) * time.Millisecond) // случайная задержка для демонстрации
	//	bar.IncrBy(1)
	//}

	p.Wait()
}

func (c *composer) StopApplication() {

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
