import { Injectable } from '@angular/core';
import { ISet, ISetPayload } from './sets-database.service';
import { IPlay, PlaysPayload } from './plays-database.service';
import { AppConfigService } from '../app-config.service';
import { HttpClient } from '@angular/common/http';
import { Observable, ReplaySubject } from 'rxjs';

export class StaticData {
  /**
   * List of known sets in the TopShot contract.
   */
  sets: ISet[] = [];
  /**
   *List of known plays in the TopShot contract.
   */
  plays: IPlay[] = [];

  constructor(plays: IPlay[], sets: ISet[]) {
    this.plays = plays;
    this.sets = sets;
  }

  findPlay(playId: number): IPlay | undefined {
    return this.plays.find(p => p.PlayId === playId);
  }

  findSet(setId: number): ISet | undefined {
    return this.sets.find(s => s.SetId === setId);
  }

}


@Injectable({
  providedIn: 'root'
})
export class StaticDataService {

  private replaySubject = new ReplaySubject<StaticData>(1);

  data$: Observable<StaticData> = this.replaySubject.asObservable();

  constructor(
    private config: AppConfigService,
    private http: HttpClient)
  {
    http.get<ISetPayload>(this.config.getSetsDataUrl()).subscribe(setsPayload => {
      http.get<PlaysPayload>(this.config.getPlaysDataUrl()).subscribe(playsPayload => {
        this.replaySubject.next(new StaticData(playsPayload.Data, setsPayload.Data));
      });
    });
  }

}
