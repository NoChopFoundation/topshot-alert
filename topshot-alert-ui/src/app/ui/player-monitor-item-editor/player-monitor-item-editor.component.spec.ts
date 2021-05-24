import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlayerMonitorItemEditorComponent } from './player-monitor-item-editor.component';

describe('PlayerMonitorItemEditorComponent', () => {
  let component: PlayerMonitorItemEditorComponent;
  let fixture: ComponentFixture<PlayerMonitorItemEditorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PlayerMonitorItemEditorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PlayerMonitorItemEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
