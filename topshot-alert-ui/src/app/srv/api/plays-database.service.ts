import { Injectable } from '@angular/core';
import { AppConfigService } from '../app-config.service';
import { HttpClient } from '@angular/common/http';

export interface IPlay {
  PlayId: number
  NbaSeason: string         // 2019-20
  TeamAtMomentNBAID: string // "1610612737"
  PlayCategory: string      // Handles
  JerseyNumber: string      // "11"
  PlayerPosition: string    // G
  DateOfMoment: string      // 2019-11-06 00:30:00 +0000 UTC
  PlayType: string          // Handles
  FullName: string          // Trae Young
  PrimaryPosition: string   // PG
  TeamAtMoment: string      // "Atlanta Hawks"
  Sets: number[]            // 2,5,8
  EditionCounts: number[]   // 15000,1,10000

  _fullNameLowerCase?: string // set on the client
}

export interface PlaysPayload {
  Data: IPlay[]
}

@Injectable({
  providedIn: 'root'
})
export class PlaysDatabaseService {

  private cachedPayload: PlaysPayload | undefined;

  constructor(
    private config: AppConfigService,
    private http: HttpClient)
  {
  }

  async getPlays(): Promise<IPlay[]> {
    await this.ensurePayload();
    return this.cachedPayload!.Data
  }

  private async ensurePayload() : Promise<void> {
    if (!this.cachedPayload) {
      this.cachedPayload = await this.http.get<any>(this.config.getPlaysDataUrl()).toPromise();
      for (var i = 0; i < this.cachedPayload!.Data.length; i++) {
        this.cachedPayload!.Data[i]._fullNameLowerCase = this.cachedPayload!.Data[i].FullName.toLocaleLowerCase();
      }
    }
  }

}
