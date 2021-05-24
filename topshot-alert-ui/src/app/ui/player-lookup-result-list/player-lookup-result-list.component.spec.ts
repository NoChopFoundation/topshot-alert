import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlayerLookupResultListComponent } from './player-lookup-result-list.component';

describe('PlayerLookupResultListComponent', () => {
  let component: PlayerLookupResultListComponent;
  let fixture: ComponentFixture<PlayerLookupResultListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PlayerLookupResultListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PlayerLookupResultListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
