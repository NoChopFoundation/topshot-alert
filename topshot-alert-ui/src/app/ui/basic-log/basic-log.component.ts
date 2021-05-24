import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet, SetsDatabaseService } from 'src/app/srv/api/sets-database.service';
import { PlaysDatabaseService } from 'src/app/srv/api/plays-database.service';
import { DetectedConsiderEvent, EventBusService } from 'src/app/srv/events/event-bus.service';

@Component({
  selector: 'app-basic-log',
  templateUrl: './basic-log.component.html',
  styleUrls: ['./basic-log.component.scss']
})
export class BasicLogComponent implements OnInit, OnDestroy {

  logEntries: DetectedConsiderEvent[] = [];

  private plays: IPlay[] = []
  private sets: ISet[] = []
  private subscription: Subscription | undefined;

  constructor(
    private playsDb: PlaysDatabaseService,
    private setsDb: SetsDatabaseService,
    private bus: EventBusService) { }

  ngOnInit(): void {
    this.subscription = this.bus.DetectedMoments.subscribe(
      data => {
        this.logEntries.unshift(data);
      }
    );
    this.playsDb.getPlays().then(aPlays => {
      this.plays = aPlays;
      this.setsDb.getSetList().then(s => {
        this.sets = s;
      });
    });
  }
  ngOnDestroy(): void {
    this.subscription?.unsubscribe();
  }

  stringify(evt: DetectedConsiderEvent): string {
    let hdr = "";
    const play = this.findPlayFromModel(evt.Event.PlayId);
    if (play) {
      hdr += play!.FullName + " ";
    }
    const set = this.findSetFromModel(evt.Event.SetId);
    if (set) {
      hdr += set!.Name + " ";
    }
    return hdr + JSON.stringify(evt);
  }

  clearLog(): void {
    this.logEntries = [];
  }

  private findPlayFromModel(playId: number): IPlay | undefined {
    return this.plays.find(e => e.PlayId === playId);
  }

  private findSetFromModel(setId: number): ISet | undefined {
    return this.sets.find(e => e.SetId === setId);
  }


}
