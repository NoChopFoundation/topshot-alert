import { HttpClient } from '@angular/common/http';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Injectable, NgModule } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { AppConfigService } from '../app-config.service';

import { ISetPayload, SetsDatabaseService } from './sets-database.service';

const dummySets: ISetPayload = {
  "Data": [
    {
      "SetId": 1,
      "Name": "Genesis"
    },
    {
      "SetId": 2,
      "Name": "Base Set"
    },
    {
      "SetId": 3,
      "Name": "Platinum Ice"
    },
    {
      "SetId": 4,
      "Name": "Holo MMXX"
    },
    {
      "SetId": 5,
      "Name": "Metallic Gold LE"
    },
    {
      "SetId": 6,
      "Name": "Early Adopters"
    },
    {
      "SetId": 7,
      "Name": "Rookie Debut"
    },
    {
      "SetId": 8,
      "Name": "Cosmic"
    },
    {
      "SetId": 9,
      "Name": "For the Win"
    },
    {
      "SetId": 10,
      "Name": "Denied!"
    },
    {
      "SetId": 11,
      "Name": "Throwdowns"
    },
    {
      "SetId": 12,
      "Name": "From the Top"
    },
    {
      "SetId": 13,
      "Name": "With the Strip"
    },
    {
      "SetId": 14,
      "Name": "Hometown Showdown: Cali vs. NY"
    },
    {
      "SetId": 15,
      "Name": "So Fresh"
    },
    {
      "SetId": 16,
      "Name": "First Round"
    },
    {
      "SetId": 17,
      "Name": "Conference Semifinals"
    },
    {
      "SetId": 18,
      "Name": "Western Conference Finals"
    },
    {
      "SetId": 19,
      "Name": "Eastern Conference Finals"
    },
    {
      "SetId": 20,
      "Name": "2020 NBA Finals"
    },
    {
      "SetId": 21,
      "Name": "The Finals"
    },
    {
      "SetId": 22,
      "Name": "Got Game"
    },
    {
      "SetId": 23,
      "Name": "Lace 'Em Up"
    },
    {
      "SetId": 24,
      "Name": "MVP Moves"
    },
    {
      "SetId": 25,
      "Name": "Run It Back"
    },
    {
      "SetId": 26,
      "Name": "Base Set"
    },
    {
      "SetId": 27,
      "Name": "Platinum Ice"
    },
    {
      "SetId": 28,
      "Name": "Holo Icon"
    },
    {
      "SetId": 29,
      "Name": "Metallic Gold LE"
    },
    {
      "SetId": 30,
      "Name": "Season Tip-off"
    },
    {
      "SetId": 31,
      "Name": "Deck the Hoops"
    },
    {
      "SetId": 32,
      "Name": "Cool Cats"
    },
    {
      "SetId": 33,
      "Name": "The Gift"
    },
    {
      "SetId": 34,
      "Name": "Seeing Stars"
    },
    {
      "SetId": 35,
      "Name": "Rising Stars"
    },
    {
      "SetId": 36,
      "Name": "2021 All-Star Game"
    }
  ]
};

@Injectable()
class SetUnitData extends SetsDatabaseService {

  protected async ensurePayload() : Promise<void> {
    this.setPayload(dummySets);
  }

}

@NgModule({
  declarations: [],
  imports: [
    HttpClientTestingModule
  ],
  providers: [ {
    provide: SetsDatabaseService,
    useClass: SetUnitData
  }, AppConfigService, HttpClient],
  bootstrap: []
})
export class SetsDatabaseTestingModule {}

describe('SetsDatabaseService', () => {
  let service: SetsDatabaseService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [SetsDatabaseTestingModule]
    });
    service = TestBed.inject(SetsDatabaseService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
