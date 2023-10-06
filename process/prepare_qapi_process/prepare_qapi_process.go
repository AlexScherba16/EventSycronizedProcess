package prepare_qapi_process

import (
	"eventsyncprocess/internal/pkg/event"
	bus "eventsyncprocess/internal/pkg/event_bus"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"time"
)

type PrepareQapiProcess struct {
	progress *mpb.Progress
	eventBus *bus.EventBus
}

func NewProcess(p *mpb.Progress, eventBus *bus.EventBus) PrepareQapiProcess {
	return PrepareQapiProcess{
		progress: p,
		eventBus: eventBus,
	}
}

func (p *PrepareQapiProcess) Run() {
	queue := make([]*mpb.Bar, 2)

	taskColor := color.New(color.FgCyan)
	taskFirst := taskColor.Sprint("[CreateQapiStep]")
	taskSecond := taskColor.Sprint("[WaitQapiReadyStep]")
	processName := "PrepareQapiProcess : "

	queue[0] = p.progress.AddBar(8,
		mpb.PrependDecorators(
			decor.Name(processName, decor.WC{W: 30, C: decor.DidentRight}),
			//decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
			decor.OnComplete(
				decor.Name(taskFirst, decor.WC{W: 30, C: decor.DidentRight}, decor.WCSyncSpaceR), ""),
		),
		mpb.AppendDecorators(
			//decor.Name(taskFirst, decor.WC{W: len(taskFirst) + 5, C: decor.DidentRight}),
			decor.Elapsed(decor.ET_STYLE_GO),
			//decor.OnComplete(decor.Elapsed(decor.ET_STYLE_GO), "done"),
		),
		mpb.BarFillerClearOnComplete(),
	)
	queue[1] = p.progress.AddBar(12,
		mpb.BarQueueAfter(queue[0]), // this bar is queued
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(processName, decor.WC{W: 30, C: decor.DidentRight}),
			decor.OnComplete(
				decor.Name(taskSecond, decor.WC{W: 30}, decor.WCSyncSpaceR), ""),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Elapsed(decor.ET_STYLE_GO), "done"),
		),
	)
	complete(queue[0])
	complete(queue[1])

	e := event.NewEvent("WaitQapiReadyStep", 42)

	//fmt.Println("Process publish")
	p.eventBus.Publish("WaitQapiReadyStep", e)
	//fmt.Println("Process publish 2")
	//p.eventBus.Publish("WaitQapiReadyStep", e)
	//fmt.Println("Process publish exit")
}

func complete(bar *mpb.Bar) {
	for !bar.Completed() {
		time.Sleep(500 * time.Millisecond)
		bar.IncrBy(1)
	}
}
