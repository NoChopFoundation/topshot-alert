import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet } from 'src/app/srv/api/sets-database.service';
import { Util_GetSetNameWithEditionInfo_ViaSetId } from 'src/app/srv/api/util';
import { MonitoredItem } from 'src/app/srv/model-storage.service';

@Component({
  selector: 'app-player-monitor-item',
  templateUrl: './player-monitor-item.component.html',
  styleUrls: ['./player-monitor-item.component.scss']
})
export class PlayerMonitorItemComponent implements OnInit {

  @Input() model: MonitoredItem
  @Input() play: IPlay
  @Input() set: ISet
  @Output() deletionEvent = new EventEmitter<MonitoredItem>();
  @Output() saveEvent = new EventEmitter<MonitoredItem>();

  @Input() CompleteSets: ISet[] = []

  IsMinimized = true;

  constructor() {
    this.play = <IPlay>{};
    this.set = <ISet>{};
    this.model = <MonitoredItem>{};
  }

  ngOnInit(): void {
  }

  onDelete(): void {
    this.deletionEvent.emit(this.model);
  }

  onSave(): void {
    this.saveEvent.emit(this.model);
  }

  GetSetNameWithEditionInfo_ViaSetId(play: IPlay, setId: number): string {
    return Util_GetSetNameWithEditionInfo_ViaSetId(this.CompleteSets, play, setId);
  }

}
