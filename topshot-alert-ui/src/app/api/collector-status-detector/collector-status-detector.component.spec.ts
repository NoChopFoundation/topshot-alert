import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CollectorStatusDetectorComponent } from './collector-status-detector.component';

describe('CollectorStatusDetectorComponent', () => {
  let component: CollectorStatusDetectorComponent;
  let fixture: ComponentFixture<CollectorStatusDetectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CollectorStatusDetectorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CollectorStatusDetectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
