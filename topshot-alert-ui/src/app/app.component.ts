import { Component, OnInit } from '@angular/core';
import { IPlay, PlaysDatabaseService } from './srv/api/plays-database.service';
import { ISet, SetsDatabaseService } from './srv/api/sets-database.service';
import { PlayerLookupPlaySetSelectedEvent } from './ui/player-lookup-result-list/player-lookup-result-list.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  sets?: ISet[];
  plays?: IPlay[];

  constructor(
    private playsDb: PlaysDatabaseService,
    private setsDb: SetsDatabaseService) {
  }

  ngOnInit() {
    this.playsDb.getPlays().then(p => {
      this.plays = p;
    });
    this.setsDb.getSetList().then(s => {
      this.sets = s;
    });
  }

}
