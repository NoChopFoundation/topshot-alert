import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet } from 'src/app/srv/api/sets-database.service';
import { MonitoredItem } from 'src/app/srv/model-storage.service';
import { Md5 } from 'ts-md5';

//export interface PlayerMonitorItemEditorDeleteEvent {
//}

@Component({
  selector: 'app-player-monitor-item-editor',
  templateUrl: './player-monitor-item-editor.component.html',
  styleUrls: ['./player-monitor-item-editor.component.scss']
})
export class PlayerMonitorItemEditorComponent implements OnInit {

  @Input() model: MonitoredItem;
  @Output() deletionEvent = new EventEmitter<MonitoredItem>();
  @Output() saveEvent = new EventEmitter<MonitoredItem>();

  @Input() CompleteSets: ISet[] = []

  content = `{}`;

  constructor() {
    this.model = <MonitoredItem>{}
  }

  ngOnInit(): void {
    this.content = this.model.considerScript;
  }

  onDelete(): void {
    this.deletionEvent.emit(this.model);
  }

  onSave(): void {
    this.model.considerScript = this.content;
    this.model.considerScriptMd5 = <string>Md5.hashAsciiStr(this.model.considerScript);
    this.saveEvent.emit(this.model);
  }

}
