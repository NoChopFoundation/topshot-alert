import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlayerMonitorComponent } from './player-monitor.component';

describe('PlayerMonitorComponent', () => {
  let component: PlayerMonitorComponent;
  let fixture: ComponentFixture<PlayerMonitorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PlayerMonitorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PlayerMonitorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
