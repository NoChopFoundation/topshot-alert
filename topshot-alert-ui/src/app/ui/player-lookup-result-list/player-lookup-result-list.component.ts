import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet } from 'src/app/srv/api/sets-database.service';
import { Util_GetSetNameWithEditionInfo_ViaPlaySetIndex } from 'src/app/srv/api/util';
import { AppConfigService } from 'src/app/srv/app-config.service';
import { ModelStorageService } from 'src/app/srv/model-storage.service';

export interface PlayerLookupPlaySetSelectedEvent {
  PlayId: number;
  SetId: number;
}

interface ResultList {
  visable: boolean;
  enabledAddToMonitoring: boolean;
  set: ISet;
}

interface IPlayWitResultList extends IPlay {
  _resultList?: ResultList;
}

@Component({
  selector: 'app-player-lookup-result-list',
  templateUrl: './player-lookup-result-list.component.html',
  styleUrls: ['./player-lookup-result-list.component.scss']
})
export class PlayerLookupResultListComponent implements OnInit {

  @Input() DisplayPlays: IPlayWitResultList[] = []
  @Input() CompleteSets: ISet[] = []

  @Output() selectionEvent = new EventEmitter<PlayerLookupPlaySetSelectedEvent>();

  constructor(
    private config: AppConfigService,
    private modelSrv: ModelStorageService) { }

  ngOnInit(): void {
    // Clear out our metadata
    this.DisplayPlays.forEach(p => p._resultList = undefined);
  }

  GetSetNameWithEditionInfo(play: IPlay, setIndex: number): string {
    return Util_GetSetNameWithEditionInfo_ViaPlaySetIndex(this.CompleteSets, play, setIndex);
  }

  addToMonitoring(playId: number, setId: number, preview: ResultList): void {
    this.selectionEvent.emit({PlayId: playId, SetId: setId});
    preview.enabledAddToMonitoring = false;
  }

  previewSetId(playId: number, setId: number): void {
    this.resetAllPreviewAreaModels();
    this.initPreviewArea(playId, setId);
  }

  /**
  addPlayerMouseOver(playId: number, setId: number): void {
    console.log("addPlayerMouseOver", playId, setId);
    this.DisplayPlays.forEach(p => {
      if (p._resultList) {
        p._resultList.visable = false;
      }
    });
    const preview = this.DisplayPlays.find(p => p.PlayId === playId);
    preview!._resultList = {
      visable: true,
      set: this.CompleteSets.find(s => s.SetId === setId)!
    };
  }
   */

  private initPreviewArea(playId: number, setId: number) {
    const isMonitored = (this.modelSrv.loadMonitoredMoments().eventMonitoring.findIndex(m => m.playId === playId && m.setId === setId) != -1);
    const preview = this.DisplayPlays.find(p => p.PlayId === playId);
    preview!._resultList = {
      visable: true,
      enabledAddToMonitoring: !isMonitored,
      set: this.CompleteSets.find(s => s.SetId === setId)!
    };
  }

  private resetAllPreviewAreaModels(): void {
    this.DisplayPlays.forEach(p => {
      if (p._resultList) {
        p._resultList.visable = false;
        p._resultList.enabledAddToMonitoring = true; // Maybe wan to check model
      }
    });
  }

  addPlayerMouseOut(playId: number, setId: number): void {
    console.log("addPlayerMouseOut", playId, setId);
  }

  addPlayerId(playId: number, setIds: number[]): void {
    setIds.forEach(setId => {
      this.selectionEvent.emit({PlayId: playId, SetId: setId});
    });
  }

  getRawLinkPlaySet(playId: number, setId: number): string {
    return this.config.getRawLinkPlaySetUrl(playId, setId);
  }


}
