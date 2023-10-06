package remote_process

import (
	bus "eventsyncprocess/internal/pkg/event_bus"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"time"
)

type RemoteProcess struct {
	progress *mpb.Progress
	eventBus *bus.EventBus
}

func NewProcess(p *mpb.Progress, eventBus *bus.EventBus) RemoteProcess {
	return RemoteProcess{
		progress: p,
		eventBus: eventBus,
	}
}

func (p *RemoteProcess) Run() {
	queue := make([]*mpb.Bar, 3)

	taskFirst := fmt.Sprintf("[InstallAgentStep]")
	taskSecond := fmt.Sprintf("[ShotDownAgent]")
	taskThird := fmt.Sprintf("[DeployAgent]")

	waitQapiChannel := bus.NewEventChannel(2)
	p.eventBus.Subscribe("remoteProcess", "WaitQapiReadyStep", waitQapiChannel)

	//waitBuildAgent := bus.NewEventChannel(2)
	//p.eventBus.Subscribe("remoteProcess", "BuildAgentStep", waitBuildAgent)

	queue[0] = p.progress.AddBar(14,
		mpb.PrependDecorators(
			decor.CountersNoUnit("RemoteProcess %d / %d", decor.WC{W: 30, C: decor.DidentRight}, decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Name(taskFirst, decor.WC{W: len(taskFirst) + 5, C: decor.DidentRight}),
			decor.Elapsed(decor.ET_STYLE_GO),
		),
		mpb.BarFillerClearOnComplete(),
	)
	queue[1] = p.progress.AddBar(6,
		mpb.BarQueueAfter(queue[0]), // this bar is queued
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.CountersNoUnit("RemoteProcess %d / %d", decor.WC{W: 30, C: decor.DidentRight}, decor.WCSyncWidth),
			//decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_MMSS, 0, decor.WCSyncWidth), ""),
		),
		mpb.AppendDecorators(
			decor.Name(taskSecond, decor.WC{W: len(taskSecond) + 5, C: decor.DidentRight}),
			decor.Elapsed(decor.ET_STYLE_GO),
			//decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
		),
	)
	queue[2] = p.progress.AddBar(8,
		mpb.BarQueueAfter(queue[1]), // this bar is queued
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.CountersNoUnit("RemoteProcess %d / %d", decor.WC{W: 30, C: decor.DidentRight}, decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Name(taskThird, decor.WC{W: len(taskThird) + 5, C: decor.DidentRight}),
			decor.Elapsed(decor.ET_STYLE_GO),
		),
	)
	<-waitQapiChannel
	complete(queue[0])
	complete(queue[1])

	//<-waitBuildAgent
	complete(queue[2])

	//e := event.NewEvent("WaitQapiReadyStep", 42)

	//fmt.Println("Process publish")
	//p.eventBus.Publish("WaitQapiReadyStep", e)
}

func complete(bar *mpb.Bar) {
	for !bar.Completed() {
		time.Sleep(800 * time.Millisecond)
		bar.IncrBy(1)
	}
}
