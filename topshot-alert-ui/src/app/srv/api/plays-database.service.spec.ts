import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';

import { PlaysDatabaseService, PlaysPayload } from './plays-database.service';
import { AppConfigService } from '../app-config.service';
import { HttpClient } from '@angular/common/http';
import { Injectable, NgModule } from '@angular/core';

const dummyPlays: PlaysPayload = {
  Data: [
    {
      "PlayId": 1,
      "NbaSeason": "2019-20",
      "TeamAtMomentNBAID": "1610612737",
      "PlayCategory": "Handles",
      "JerseyNumber": "11",
      "PlayerPosition": "G",
      "DateOfMoment": "2019-11-06 00:30:00 +0000 UTC",
      "PlayType": "Handles",
      "FullName": "Trae Young",
      "PrimaryPosition": "PG",
      "TeamAtMoment": "Atlanta Hawks",
      "Sets": [1,3,6]
    },
    {
      "PlayId": 2,
      "NbaSeason": "2019-20",
      "TeamAtMomentNBAID": "1610612751",
      "PlayCategory": "Handles",
      "JerseyNumber": "11",
      "PlayerPosition": "G",
      "DateOfMoment": "2019-10-27 22:00:00 +0000 UTC",
      "PlayType": "Handles",
      "FullName": "Kyrie Irving",
      "PrimaryPosition": "PG",
      "TeamAtMoment": "Brooklyn Nets",
      "Sets": [1, 8, 4]
    },
  ]
};

@Injectable()
class UnitData extends PlaysDatabaseService {
  async numberPlays(): Promise<number> {
    return Promise.resolve(dummyPlays.Data.length);
  }
}

@NgModule({
  declarations: [],
  imports: [
    HttpClientTestingModule
  ],
  providers: [ {
    provide: PlaysDatabaseService,
    useClass: UnitData
  }, AppConfigService, HttpClient],
  bootstrap: []
})
export class PlaysDatabaseTestingModule {}

describe('PlaysDatabaseService', () => {
  let service: PlaysDatabaseService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [PlaysDatabaseTestingModule],
      providers: [],
    });
    service = TestBed.inject(PlaysDatabaseService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get the mocked data', (done) => {
    service.getPlays().then(val => {
      expect(val.length).toBeGreaterThan(2);
      done()
    });
  });

});
