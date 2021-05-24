import { HttpClient, HttpRequest } from '@angular/common/http';
import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { AppConfigService } from 'src/app/srv/app-config.service';
import { ConnectionStatus, ContactedEvent, EventBusService } from 'src/app/srv/events/event-bus.service'
import { catchError } from 'rxjs/operators';
import { throwError } from 'rxjs';

interface Collector {
  CollectorId: string
  State: string
  UpdatesInInterval: number
  BlockHeight: number
  CreatedAt: string // "2021-03-25T01:42:19Z"
}

interface CollectorReply {
  Data: Collector[]
}

@Component({
  selector: 'app-collector-status-detector',
  template: '<p>Collectors status {{  CollectorStatus }}</p>'
})
export class CollectorStatusDetectorComponent implements OnInit, OnDestroy {

  CollectorStatus = "Loading...."
  private outstandingTimer: any;

  constructor(
    private config: AppConfigService,
    private bus: EventBusService,
    private http: HttpClient) {
      this.outstandingTimer = null;
  }

  ngOnInit(): void {
    this.pingCollectors();
  }

  ngOnDestroy(): void {
    console.log("CollectorStatusDetectorComponent ngOnDestroy")
    if (this.outstandingTimer != null) {
      console.log("CollectorStatusDetectorComponent clear timers")
      clearInterval(this.outstandingTimer)
    }
  }

  private pingCollectors(): void {
    this.http
      .get<CollectorReply>(this.config.getRecentCollectorsUrl())
      .subscribe(
        data => {
          this.processReply(data);
        },
        err => {
          this.CollectorStatus = "error: " + err;
          console.log("ping error", err);
          this.setupTryAgain();
        });
    }

    private setupTryAgain(): void {
      setTimeout(() => {
        this.outstandingTimer = null;
        this.pingCollectors();
      }, this.config.getCollectorPingRetryMs());
    }

    private processReply(reply: CollectorReply): void {
      this.bus.SuccessfulCollectorContact.next(new ContactedEvent())

      const summary = new Map();
      reply.Data.forEach(e => {
        if (summary.has(e.CollectorId)) {
          if (e.State === 'U') {
          } else {
            summary.set(e.CollectorId, <Collector> e);
          }
        } else {
          summary.set(e.CollectorId, <Collector> e);
        }
      });

      let upCount = 0;
      let downCount = 0;
      let eventsProcessed = 0;
      for (let v of summary.values()) {
        const collector = <Collector>v
        if (collector.State === 'U') {
          upCount ++;
        } else {
          downCount ++;
        }
        eventsProcessed += collector.UpdatesInInterval;
      }

      console.log(`Collectors UP(${upCount}) Down(${downCount}) -- Events streamed ${eventsProcessed}`);
      this.CollectorStatus = `UP(${upCount}) Down(${downCount}) -- Events streaming`;

      if (upCount > 0) {
        this.bus.ConnectionStatus.next(ConnectionStatus.Connected);
      } else {
        this.bus.ConnectionStatus.next(ConnectionStatus.Disconnected);
      }
      this.setupTryAgain();
    }

}
