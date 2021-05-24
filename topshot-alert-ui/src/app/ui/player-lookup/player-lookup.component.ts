import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet } from 'src/app/srv/api/sets-database.service';
import { PlayerLookupPlaySetSelectedEvent } from '../player-lookup-result-list/player-lookup-result-list.component';

@Component({
  selector: 'app-player-lookup',
  templateUrl: './player-lookup.component.html',
  styleUrls: ['./player-lookup.component.scss']
})
export class PlayerLookupComponent implements OnInit {

  @Input() plays?: IPlay[]
  @Input() sets?: ISet[]
  @Input() maxSearchRows = 20;
  @Output() selectionEvent = new EventEmitter<PlayerLookupPlaySetSelectedEvent>();

  IsMinimized = true;
  SearchInputText = ""

  DisplayPlays: IPlay[] = []

  constructor() {
  }

  ngOnInit(): void {
  }

  isLoading(): boolean {
    return  (!this.plays || !this.sets);
  }

  onPlayerLookupPlaySetSelectedEvent(event: PlayerLookupPlaySetSelectedEvent): void {
    this.selectionEvent.emit(event);
  }

  onPlayerNameInput(input: string) {
    var searchFor = input.trim().toLocaleLowerCase();
    if (searchFor.length > 1) {
      this.executeSearch(searchFor);
    } else {
      this.DisplayPlays = []
    }
  }

  executeSearch(searchFor: string) {
    const results: IPlay[] = [];
    for (var i = 0; i < this.plays!.length; i++) {
      const includePlayer = this.plays![i]._fullNameLowerCase?.includes(searchFor);
      if (includePlayer) {
        results.push(this.plays![i]);
        if (results.length > this.maxSearchRows) {
          break;
        }
      }
    }
    this.DisplayPlays = results;
  }

}
