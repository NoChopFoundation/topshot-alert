import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet } from 'src/app/srv/api/sets-database.service';
import { ModelStorageService, MonitoredItem, MonitoredMomentsV09 } from 'src/app/srv/model-storage.service';
import { PlayerLookupPlaySetSelectedEvent } from '../player-lookup-result-list/player-lookup-result-list.component';
import { PlaysDatabaseService } from "src/app/srv/api/plays-database.service";
import { SetsDatabaseService } from "src/app/srv/api/sets-database.service";
import { EventBusService, PurchaseListedEvent } from 'src/app/srv/events/event-bus.service';
import { ConsiderService } from "src/app/srv/consider/consider.service";
import { Subscription } from 'rxjs';
import { AppConfigService } from 'src/app/srv/app-config.service';
import { Md5 } from 'ts-md5/dist/md5'

@Component({
  selector: 'app-player-monitor',
  templateUrl: './player-monitor.component.html',
  styleUrls: ['./player-monitor.component.scss']
})
export class PlayerMonitorComponent implements OnInit, OnDestroy {

  IsMinimized = true;
  IsLoading = true;
  model: MonitoredMomentsV09;
  liveFeedSubscription: Subscription | undefined;
  liveFeedSubscriptionMsg = "";
  IsImportMinimized = true;

  plays: IPlay[] = []
  sets: ISet[] = []

  constructor(
    private config: AppConfigService,
    private runtimeSrv: ConsiderService,
    private storage: ModelStorageService,
    private playsDb: PlaysDatabaseService,
    private setsDb: SetsDatabaseService,
    private bus: EventBusService) {
    this.model = <MonitoredMomentsV09>{}
  }

  ngOnDestroy(): void {
    this.liveFeedSubscription?.unsubscribe();
    this.liveFeedSubscription = undefined;
  }

  onReStartMonitor(): void {
    if (!this.liveFeedSubscription) {
      this.liveFeedSubscription = this.bus.LiveStream.subscribe(
        (event: PurchaseListedEvent) => {
          this.liveFeedSubscriptionMsg = "";
          this.onPurchaseListedEvent(event);
        },
        (error) => {
          this.liveFeedSubscriptionMsg = JSON.stringify(error);
          console.log("monitor error received", error);
          this.liveFeedSubscription = undefined;
        },
        () => {
          // Completed, shouldn't as we want to poll until we unsubscribe
          this.liveFeedSubscription = undefined;
        }
      );
    }
  }

  private onPurchaseListedEvent(evt: PurchaseListedEvent): void {
    this.runtimeSrv.onMomentEvent(this.model, evt);
  }

  ngOnInit(): void {
    this.model = this.storage.loadMonitoredMoments();
    this.playsDb.getPlays().then(aPlays => {
      this.plays = aPlays;
      this.setsDb.getSetList().then(s => {
        this.sets = s;
        this.IsLoading = false;
      });
    });
    this.onReStartMonitor();
  }

  onUserInitiatePause(): void {
    this.liveFeedSubscriptionMsg = "user paused";
    this.liveFeedSubscription?.unsubscribe();
    this.liveFeedSubscription = undefined;
  }

  onPlayerLookupPlaySetSelectedEvent(event: PlayerLookupPlaySetSelectedEvent): void {
    console.log("bro update model", event.PlayId, event.SetId);

    const idx = this.model.eventMonitoring.findIndex(e => e.playId === event.PlayId && e.setId === event.SetId);
    if (idx === -1) {
      const considerScript = this.config.getInitialConsiderScript(this.plays, this.sets, event.PlayId, event.SetId);
      this.model.eventMonitoring.push({
        playId: event.PlayId,
        setId: event.SetId,
        considerScript: considerScript,
        considerScriptMd5: <string>Md5.hashAsciiStr(considerScript)
      });
      this.storage.persistMonitoredMoments(this.model);
    }
  }

  onTopButtonClicked(): void {
    if (this.IsMinimized && this.model.eventMonitoring.length === 0) {
      // Ignore, nothing to show
    } else {
      this.IsMinimized = !this.IsMinimized;
    }
  }

  onDeleteAll(): void {
    if (confirm("Delete all monitored items?")) {
      this.model.eventMonitoring = [];
      this.storage.persistMonitoredMoments(this.model);
    }
  }

  findPlayFromModel(itemModel: MonitoredItem): IPlay {
    return this.plays.find(e => e.PlayId === itemModel.playId)!;
  }

  findSetFromModel(itemModel: MonitoredItem): ISet {
    return this.sets.find(e => e.SetId === itemModel.setId)!;
  }

  onDelete(event: MonitoredItem): void {
    const idx = this.model.eventMonitoring.findIndex(e => e.playId === event.playId && e.setId === event.setId);
    if (idx !== -1) {
      this.model.eventMonitoring.splice(idx, 1);
      this.storage.persistMonitoredMoments(this.model);
    }
  }

  onSave(event: MonitoredItem): void {
    const idx = this.model.eventMonitoring.findIndex(e => e.playId === event.playId && e.setId === event.setId);
    if (idx !== -1) {
      this.model.eventMonitoring[idx] = event;
      this.storage.persistMonitoredMoments(this.model);
    }
  }

  copyModelToClipboard(): void {
    this.storage.loadMonitoredMoments();

    const selBox = document.createElement('textarea');
    selBox.style.position = 'fixed';
    selBox.style.left = '0';
    selBox.style.top = '0';
    selBox.style.opacity = '0';
    selBox.value = this.storage.exportModel();
    document.body.appendChild(selBox);
    selBox.focus();
    selBox.select();
    document.execCommand('copy');
    document.body.removeChild(selBox);

    alert("The current monitoring model save to your clipboard. Save this in a text file for later import or backup.");
  }

  fileChanged(evt: any): void {
    let fileReader = new FileReader();
    fileReader.onload = (e) => {
      if (this.storage.importModel(<string> fileReader.result)) {
        this.model = this.storage.loadMonitoredMoments();
      }
    }
    fileReader.readAsText(evt.target.files[0]);
    evt.srcElement.value = "";
  }

}
