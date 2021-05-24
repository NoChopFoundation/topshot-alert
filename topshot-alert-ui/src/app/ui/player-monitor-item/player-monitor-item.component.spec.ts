import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlayerMonitorItemComponent } from './player-monitor-item.component';

describe('PlayerMonitorItemComponent', () => {
  let component: PlayerMonitorItemComponent;
  let fixture: ComponentFixture<PlayerMonitorItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PlayerMonitorItemComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PlayerMonitorItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
